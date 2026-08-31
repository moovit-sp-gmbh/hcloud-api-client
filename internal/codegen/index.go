// Package codegen provides static analysis of the generated OpenAPI clients and
// of the cobra command tree, plus the emitters that turn an operation into a
// command.
//
// Nothing in this package may import internal/cli or internal/db. internal/db
// opens db.sqlite and runs migrations from init(), so importing the command tree
// would make the generator create and migrate a database in whatever directory
// it happens to run from.
package codegen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Location is where a parameter travels in the request.
type Location string

const (
	InPath   Location = "path"
	InQuery  Location = "query"
	InHeader Location = "header"
)

// Param is one path, query or header parameter of an operation.
type Param struct {
	Name     string   `json:"name"`   // wire name, e.g. "orgName"
	GoName   string   `json:"goName"` // argument or Params field name
	GoType   string   `json:"goType"` // as written in the client, e.g. "*float32"
	In       Location `json:"in"`
	Optional bool     `json:"optional"`         // pointer field in the Params struct
	SpecType string   `json:"specType"`         // declared OpenAPI type, e.g. "number"
	Scalar   string   `json:"scalar,omitempty"` // builtin a named type resolves to, e.g. AgentStatus -> string
	Doc      string   `json:"doc,omitempty"`
}

// Field is one field of a request or response DTO.
type Field struct {
	GoName   string `json:"goName"`
	JSONName string `json:"jsonName"`
	GoType   string `json:"goType"`
	Optional bool   `json:"optional"`
	Scalar   string `json:"scalar,omitempty"` // builtin a named type resolves to
	Form     bool   `json:"form,omitempty"`   // carries a form tag, i.e. a query parameter
	Doc      string `json:"doc,omitempty"`
}

// Body describes an operation's request body.
type Body struct {
	TypeName string  `json:"typeName"` // e.g. "PostV1RegisterJSONRequestBody"
	Encoding string  `json:"encoding"` // "json" or "formdata"
	DTO      string  `json:"dto"`      // what the alias resolves to, e.g. "AgentRegisterDto"
	Freeform bool    `json:"freeform"` // map[string]interface{} — the patcher genericized it
	Opaque   bool    `json:"opaque"`   // reaches a type with unexported fields; must be sent as raw JSON
	Fields   []Field `json:"fields,omitempty"`
}

// Constructable reports whether the body can be built field by field. When it
// cannot, the command has to take raw JSON through a --body flag, the way
// internal/cli/mothership/connect.go does.
func (b *Body) Constructable() bool {
	return b != nil && !b.Freeform && !b.Opaque && len(b.Fields) > 0
}

// Response is one documented status code of an operation.
type Response struct {
	Status int     `json:"status"`
	Field  string  `json:"field"`          // e.g. "JSON200"
	GoType string  `json:"goType"`         // e.g. "*StatsDto"
	DTO    string  `json:"dto"`            // e.g. "StatsDto"
	List   bool    `json:"list,omitempty"` // the body is an array, so it has no addressable fields
	Fields []Field `json:"fields,omitempty"`
}

// Operation is a single callable client operation.
type Operation struct {
	Client     string     `json:"client"` // package name, e.g. "mothershipClient"
	Name       string     `json:"name"`   // e.g. "GetV1OrgOrgNameStats"
	Method     string     `json:"method"`
	Path       string     `json:"path"` // templated, e.g. "/v1/org/{orgName}/stats"
	PathParams []Param    `json:"pathParams,omitempty"`
	Params     []Param    `json:"params,omitempty"` // query and header
	ParamsType string     `json:"paramsType,omitempty"`
	Body       *Body      `json:"body,omitempty"`
	Responses  []Response `json:"responses,omitempty"`

	// Blocked says why the operation cannot be called from outside the
	// generated package, and is empty when it can. oapi-codegen renders an
	// un-flattened union path parameter as a struct with an unexported union
	// field and no constructor, which makes the operation unreachable no matter
	// what a command does — cmd/patcher exists to prevent this, and anything
	// still marked here is a spec its transforms do not cover.
	Blocked string `json:"blocked,omitempty"`
}

// Callable reports whether a command can be generated for the operation.
func (o *Operation) Callable() bool { return o.Blocked == "" }

// Key identifies an operation across clients. Operation names collide between
// clients (both idp and high5 declare GetV1User), so nothing may key on the
// name alone.
func (o *Operation) Key() string { return o.Client + "." + o.Name }

// Success returns the response for the operation's 2xx status with a typed
// body, preferring 200 over 201. ok is false when the operation documents no
// typed success body, in which case a command must render the raw res.Body.
func (o *Operation) Success() (Response, bool) {
	var best Response
	for _, r := range o.Responses {
		if r.Status < 200 || r.Status > 299 || r.DTO == "" {
			continue
		}
		if best.Status == 0 || r.Status < best.Status {
			best = r
		}
	}
	return best, best.Status != 0
}

// HeaderParam returns the operation's header parameter with the given wire
// name. Commands prefer the typed Params field over a RequestEditorFn closure
// when the operation declares the header.
func (o *Operation) HeaderParam(name string) (Param, bool) {
	for _, p := range o.Params {
		if p.In == InHeader && strings.EqualFold(p.Name, name) {
			return p, true
		}
	}
	return Param{}, false
}

// Client is one generated client package.
type Client struct {
	Package string                `json:"package"`
	File    string                `json:"file"`
	Ops     map[string]*Operation `json:"ops"`
}

// Sorted returns the client's operations ordered by path then method, which is
// how coverage reports and the endpoint picker list them.
func (c *Client) Sorted() []*Operation {
	out := make([]*Operation, 0, len(c.Ops))
	for _, op := range c.Ops {
		out = append(out, op)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Path != out[j].Path {
			return out[i].Path < out[j].Path
		}
		return out[i].Method < out[j].Method
	})
	return out
}

// Index is every operation the CLI can call, per client package.
type Index struct {
	Clients map[string]*Client `json:"clients"`
}

// Op looks an operation up by client package and operation name.
func (ix *Index) Op(client, name string) (*Operation, bool) {
	c, ok := ix.Clients[client]
	if !ok {
		return nil, false
	}
	op, ok := c.Ops[name]
	return op, ok
}

// ClientNames returns the client package names in a stable order.
func (ix *Index) ClientNames() []string {
	out := make([]string, 0, len(ix.Clients))
	for name := range ix.Clients {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// LoadIndex parses every generated client under root and returns the operations
// they expose. root is the repository root.
func LoadIndex(root string) (*Index, error) {
	files, err := filepath.Glob(filepath.Join(root, "internal", "clients", "*", "*.gen.go"))
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no generated clients under %s/internal/clients", root)
	}
	sort.Strings(files)

	ix := &Index{Clients: map[string]*Client{}}
	for _, file := range files {
		c, err := loadClient(file)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", file, err)
		}
		ix.Clients[c.Package] = c
	}
	return ix, nil
}

func loadClient(path string) (*Client, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	reg := newRegistry(f)
	c := &Client{Package: f.Name.Name, File: path, Ops: map[string]*Operation{}}

	names, err := operationNames(f)
	if err != nil {
		return nil, err
	}

	builders := map[string]*ast.FuncDecl{}
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil || fn.Body == nil {
			continue
		}
		builders[fn.Name.Name] = fn
	}

	for _, name := range names {
		op := &Operation{Client: c.Package, Name: name}

		// The path and method live in New<Op>Request for parameter-only
		// operations and in New<Op>RequestWithBody for bodied ones, because the
		// former just marshals the body and delegates.
		builder := builders["New"+name+"RequestWithBody"]
		if builder == nil || !hasOperationPath(builder) {
			builder = builders["New"+name+"Request"]
		}
		if builder == nil {
			return nil, fmt.Errorf("operation %s: no request builder found", name)
		}
		method, tmpl, wire, err := parseBuilder(builder)
		if err != nil {
			return nil, fmt.Errorf("operation %s: %w", name, err)
		}
		op.Method = method

		sig, ok := interfaceMethod(f, name+"WithResponse")
		if !ok {
			return nil, fmt.Errorf("operation %s: not in ClientWithResponsesInterface", name)
		}
		pathArgs := signaturePathParams(sig, name)

		op.Path, op.PathParams, err = resolvePath(tmpl, wire, pathArgs)
		if err != nil {
			return nil, fmt.Errorf("operation %s: %w", name, err)
		}
		for i := range op.PathParams {
			p := &op.PathParams[i]
			p.Scalar = reg.scalar(p.GoType)
			if base := strings.TrimPrefix(p.GoType, "*"); !strings.ContainsAny(base, ".[]{}*") && reg.opaque(base) {
				op.Blocked = fmt.Sprintf("path parameter %s has type %s, which the client declares with an unexported union field and no constructor", p.Name, base)
			}
		}

		if reg.has(name + "Params") {
			op.ParamsType = name + "Params"
			op.Params = paramsFields(reg, op.ParamsType, wire)
		}
		op.Body = reg.body(name)
		op.Responses = reg.responses(name)

		c.Ops[name] = op
	}
	return c, nil
}

// operationNames returns the operation names declared by
// ClientWithResponsesInterface. That interface is the authoritative list of
// what the CLI can call: it holds exactly the operations oapi-codegen managed
// to render.
func operationNames(f *ast.File) ([]string, error) {
	var iface *ast.InterfaceType
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok || ts.Name.Name != "ClientWithResponsesInterface" {
				continue
			}
			if it, ok := ts.Type.(*ast.InterfaceType); ok {
				iface = it
			}
		}
	}
	if iface == nil {
		return nil, fmt.Errorf("no ClientWithResponsesInterface declaration")
	}

	var names []string
	for _, m := range iface.Methods.List {
		if len(m.Names) != 1 {
			continue
		}
		n := m.Names[0].Name
		if !strings.HasSuffix(n, "WithResponse") {
			continue
		}
		op := strings.TrimSuffix(n, "WithResponse")
		// Bodied operations also declare a variant that takes the body as an
		// io.Reader (WithBody) or as form values (WithFormdataBody). Those are
		// the same endpoint, so skip them and keep the typed call.
		if strings.HasSuffix(op, "Body") {
			continue
		}
		names = append(names, op)
	}
	sort.Strings(names)
	return names, nil
}

func interfaceMethod(f *ast.File, method string) (*ast.FuncType, bool) {
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok || ts.Name.Name != "ClientWithResponsesInterface" {
				continue
			}
			it, ok := ts.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}
			for _, m := range it.Methods.List {
				if len(m.Names) == 1 && m.Names[0].Name == method {
					if ft, ok := m.Type.(*ast.FuncType); ok {
						return ft, true
					}
				}
			}
		}
	}
	return nil, false
}

// signaturePathParams pulls the positional path arguments out of an operation's
// WithResponse signature: everything that is not the context, the params
// struct, the body, a content type or the request editors.
func signaturePathParams(sig *ast.FuncType, op string) []Param {
	var out []Param
	if sig.Params == nil {
		return out
	}
	for _, field := range sig.Params.List {
		goType := types.ExprString(field.Type)
		switch {
		case goType == "context.Context",
			goType == "*"+op+"Params",
			goType == op+"JSONRequestBody",
			goType == op+"FormdataRequestBody",
			goType == "io.Reader",
			strings.HasSuffix(goType, "RequestEditorFn"):
			continue
		}
		if _, isEllipsis := field.Type.(*ast.Ellipsis); isEllipsis {
			continue
		}
		for _, name := range field.Names {
			if name.Name == "contentType" {
				continue
			}
			out = append(out, Param{GoName: name.Name, GoType: goType, In: InPath})
		}
	}
	return out
}

// wireParam is one runtime.StyleParamWithOptions call found in a request
// builder: the parameter's wire name, where it travels and its declared
// OpenAPI type.
type wireParam struct {
	name string
	in   Location
	typ  string
}

func hasOperationPath(fn *ast.FuncDecl) bool {
	found := false
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		if as, ok := n.(*ast.AssignStmt); ok && len(as.Lhs) == 1 {
			if id, ok := as.Lhs[0].(*ast.Ident); ok && id.Name == "operationPath" {
				found = true
			}
		}
		return !found
	})
	return found
}

// parseBuilder extracts the HTTP method, the fmt.Sprintf path template and the
// styled parameters from a New<Op>Request function body.
func parseBuilder(fn *ast.FuncDecl) (method, tmpl string, wire []wireParam, err error) {
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.AssignStmt:
			if tmpl != "" || len(node.Lhs) != 1 {
				return true
			}
			id, ok := node.Lhs[0].(*ast.Ident)
			if !ok || id.Name != "operationPath" || len(node.Rhs) != 1 {
				return true
			}
			call, ok := node.Rhs[0].(*ast.CallExpr)
			if !ok || !isSelector(call.Fun, "fmt", "Sprintf") || len(call.Args) == 0 {
				return true
			}
			if lit, ok := stringLit(call.Args[0]); ok {
				tmpl = lit
			}
		case *ast.CallExpr:
			switch {
			case isSelector(node.Fun, "http", "NewRequest") && len(node.Args) > 0:
				if method == "" {
					method = httpMethod(node.Args[0])
				}
			case isSelector(node.Fun, "runtime", "StyleParamWithOptions") && len(node.Args) >= 3:
				if p, ok := styleParam(node); ok {
					wire = append(wire, p)
				}
			}
		}
		return true
	})

	if method == "" {
		return "", "", nil, fmt.Errorf("no http.NewRequest call")
	}
	if tmpl == "" {
		return "", "", nil, fmt.Errorf("no operationPath assignment")
	}
	return method, tmpl, wire, nil
}

func styleParam(call *ast.CallExpr) (wireParam, bool) {
	name, ok := stringLit(call.Args[2])
	if !ok {
		return wireParam{}, false
	}
	p := wireParam{name: name}
	opts, ok := call.Args[len(call.Args)-1].(*ast.CompositeLit)
	if !ok {
		return p, true
	}
	for _, elt := range opts.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}
		switch key.Name {
		case "ParamLocation":
			sel, ok := kv.Value.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			switch sel.Sel.Name {
			case "ParamLocationPath":
				p.in = InPath
			case "ParamLocationQuery":
				p.in = InQuery
			case "ParamLocationHeader":
				p.in = InHeader
			}
		case "Type":
			if lit, ok := stringLit(kv.Value); ok {
				p.typ = lit
			}
		}
	}
	return p, true
}

// resolvePath turns the generated fmt.Sprintf template back into an OpenAPI
// path. Every verb in these clients is %s, and the path parameters appear in
// the styled-parameter calls in the same order they appear in the template.
func resolvePath(tmpl string, wire []wireParam, args []Param) (string, []Param, error) {
	var pathWire []wireParam
	for _, w := range wire {
		if w.in == InPath {
			pathWire = append(pathWire, w)
		}
	}

	parts := strings.Split(tmpl, "%s")
	if len(parts)-1 != len(pathWire) {
		return "", nil, fmt.Errorf("path template %q has %d verbs but %d path parameters", tmpl, len(parts)-1, len(pathWire))
	}

	var sb strings.Builder
	for i, part := range parts {
		sb.WriteString(part)
		if i < len(pathWire) {
			sb.WriteString("{" + pathWire[i].name + "}")
		}
	}

	// Marry the wire names to the Go argument names from the signature. They
	// agree in order; when the counts disagree the signature wins for the Go
	// side and the wire name is left to fall back to the argument name.
	out := make([]Param, 0, len(pathWire))
	for i, w := range pathWire {
		p := Param{Name: w.name, GoName: w.name, GoType: "string", In: InPath, SpecType: w.typ}
		if i < len(args) {
			p.GoName = args[i].GoName
			p.GoType = args[i].GoType
		}
		out = append(out, p)
	}
	return sb.String(), out, nil
}

// paramsFields reads the query and header parameters off the generated Params
// struct. The struct carries the Go names, types and optionality; the builder's
// styled-parameter calls say where each one travels.
func paramsFields(reg *registry, typeName string, wire []wireParam) []Param {
	where := map[string]wireParam{}
	for _, w := range wire {
		if w.in == InQuery || w.in == InHeader {
			where[strings.ToLower(w.name)] = w
		}
	}

	var out []Param
	for _, f := range reg.fields(typeName) {
		p := Param{
			Name:     f.JSONName,
			GoName:   f.GoName,
			GoType:   f.GoType,
			Optional: f.Optional,
			Scalar:   f.Scalar,
			Doc:      f.Doc,
		}
		if w, ok := where[strings.ToLower(f.JSONName)]; ok {
			p.In = w.in
			p.SpecType = w.typ
		} else {
			// A field with a form tag is a query parameter; oapi-codegen only
			// emits that tag for query parameters.
			if f.Form {
				p.In = InQuery
			} else {
				p.In = InHeader
			}
		}
		out = append(out, p)
	}
	return out
}

func isSelector(e ast.Expr, pkg, name string) bool {
	sel, ok := e.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != name {
		return false
	}
	id, ok := sel.X.(*ast.Ident)
	return ok && id.Name == pkg
}

func stringLit(e ast.Expr) (string, bool) {
	lit, ok := e.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}
	v, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}
	return v, true
}

// httpMethod turns the http.MethodGet selector in a generated NewRequest call
// into "GET".
func httpMethod(e ast.Expr) string {
	if sel, ok := e.(*ast.SelectorExpr); ok {
		if strings.HasPrefix(sel.Sel.Name, "Method") {
			return strings.ToUpper(strings.TrimPrefix(sel.Sel.Name, "Method"))
		}
	}
	if lit, ok := stringLit(e); ok {
		return strings.ToUpper(lit)
	}
	return ""
}
