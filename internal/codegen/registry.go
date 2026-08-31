package codegen

import (
	"go/ast"
	"go/token"
	"go/types"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// registry resolves the type declarations of one generated client file, so a
// DTO name can be turned into the fields a command needs to fill in or read
// back out.
type registry struct {
	decls map[string]*ast.TypeSpec
}

func newRegistry(f *ast.File) *registry {
	r := &registry{decls: map[string]*ast.TypeSpec{}}
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			// A lone TypeSpec carries no doc of its own; the GenDecl holds it.
			if ts.Doc == nil && len(gd.Specs) == 1 {
				ts.Doc = gd.Doc
			}
			r.decls[ts.Name.Name] = ts
		}
	}
	return r
}

func (r *registry) has(name string) bool {
	_, ok := r.decls[name]
	return ok
}

// underlying follows alias and defined-type chains until it reaches something
// that is not a bare name in this package. oapi-codegen writes request bodies
// as `type PostV1RegisterJSONRequestBody = AgentRegisterDto`, so a body type
// usually takes two hops to reach its struct.
func (r *registry) underlying(name string) (ast.Expr, string) {
	seen := map[string]bool{}
	resolved := name
	for {
		ts, ok := r.decls[resolved]
		if !ok || seen[resolved] {
			return nil, resolved
		}
		seen[resolved] = true
		id, ok := ts.Type.(*ast.Ident)
		if !ok {
			return ts.Type, resolved
		}
		if _, known := r.decls[id.Name]; !known {
			return ts.Type, resolved
		}
		resolved = id.Name
	}
}

// structOf resolves a type name to its struct declaration, if it has one.
func (r *registry) structOf(name string) (*ast.StructType, bool) {
	expr, _ := r.underlying(name)
	st, ok := expr.(*ast.StructType)
	return st, ok
}

// fields returns the fields of a named struct type. Anonymous and embedded
// fields are skipped; none of these specs produce them.
func (r *registry) fields(name string) []Field {
	st, ok := r.structOf(name)
	if !ok || st.Fields == nil {
		return nil
	}

	var out []Field
	for _, f := range st.Fields.List {
		if len(f.Names) == 0 {
			continue
		}
		goType := types.ExprString(f.Type)
		tag := fieldTag(f)
		jsonName, jsonOpts := tagValue(tag, "json")
		formName, _ := tagValue(tag, "form")

		for _, n := range f.Names {
			if !n.IsExported() {
				continue
			}
			name := jsonName
			if name == "" {
				name = formName
			}
			if name == "" {
				name = strings.ToLower(n.Name[:1]) + n.Name[1:]
			}
			out = append(out, Field{
				GoName:   n.Name,
				JSONName: name,
				GoType:   goType,
				Optional: strings.HasPrefix(goType, "*") || jsonOpts["omitempty"],
				Scalar:   r.scalar(goType),
				Form:     formName != "",
				Doc:      fieldDoc(f, n.Name),
			})
		}
	}
	return out
}

// builtinScalars are the Go types a flag can carry directly.
var builtinScalars = map[string]bool{
	"string": true, "bool": true, "int": true, "int8": true, "int16": true,
	"int32": true, "int64": true, "uint": true, "uint8": true, "uint16": true,
	"uint32": true, "uint64": true, "float32": true, "float64": true,
}

// scalar returns the builtin a named type resolves to, or "" when it resolves
// to something a flag cannot carry. AgentStatus is a string enum and so can be
// taken as a string; Sorting is a struct and cannot.
func (r *registry) scalar(goType string) string {
	base := strings.TrimPrefix(goType, "*")
	if builtinScalars[base] {
		return base
	}
	if strings.ContainsAny(base, ".[]{}*") {
		return ""
	}

	seen := map[string]bool{}
	for !seen[base] {
		seen[base] = true
		ts, ok := r.decls[base]
		if !ok {
			return ""
		}
		id, ok := ts.Type.(*ast.Ident)
		if !ok {
			return ""
		}
		if builtinScalars[id.Name] {
			return id.Name
		}
		base = id.Name
	}
	return ""
}

// opaque reports whether a type can be constructed from outside the generated
// package. oapi-codegen renders an unrenderable union as a struct with an
// unexported `union json.RawMessage` field and no constructor, so any body that
// reaches one has to be sent as raw JSON — this is why
// internal/cli/mothership/connect.go takes a --body flag.
func (r *registry) opaque(name string) bool {
	return r.opaqueWalk(name, map[string]bool{})
}

func (r *registry) opaqueWalk(name string, seen map[string]bool) bool {
	if name == "" || seen[name] {
		return false
	}
	seen[name] = true

	expr, resolved := r.underlying(name)
	if resolved != name {
		seen[resolved] = true
	}
	st, ok := expr.(*ast.StructType)
	if !ok || st.Fields == nil {
		return false
	}

	for _, f := range st.Fields.List {
		for _, n := range f.Names {
			if !n.IsExported() {
				return true
			}
		}
		if r.opaqueWalk(namedType(f.Type), seen) {
			return true
		}
	}
	return false
}

// body describes an operation's request body, or nil when it has none. Two idp
// operations take a form-encoded body rather than JSON.
func (r *registry) body(op string) *Body {
	typeName, encoding := op+"JSONRequestBody", "json"
	if !r.has(typeName) {
		typeName, encoding = op+"FormdataRequestBody", "formdata"
	}
	if !r.has(typeName) {
		return nil
	}

	expr, dto := r.underlying(typeName)
	b := &Body{TypeName: typeName, Encoding: encoding, DTO: dto}
	if mt, ok := expr.(*ast.MapType); ok {
		// The patcher replaces a body that reaches an ambiguous discriminator
		// mapping with a free-form object, which lands here.
		if types.ExprString(mt.Key) == "string" {
			b.Freeform = true
			return b
		}
	}
	b.Fields = r.fields(dto)
	b.Opaque = r.opaque(dto)
	return b
}

// responses describes every documented status code of an operation that carries
// a typed body.
func (r *registry) responses(op string) []Response {
	st, ok := r.structOf(op + "Response")
	if !ok || st.Fields == nil {
		return nil
	}

	var out []Response
	for _, f := range st.Fields.List {
		for _, n := range f.Names {
			if !strings.HasPrefix(n.Name, "JSON") {
				continue
			}
			status, err := strconv.Atoi(strings.TrimPrefix(n.Name, "JSON"))
			if err != nil {
				continue
			}
			goType := types.ExprString(f.Type)
			resp := Response{
				Status: status,
				Field:  n.Name,
				GoType: goType,
				List:   strings.Contains(goType, "[]"),
				DTO:    namedType(f.Type),
			}
			if !resp.List {
				resp.Fields = r.fields(resp.DTO)
			}
			out = append(out, resp)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Status < out[j].Status })
	return out
}

// namedType strips pointers, slices and package qualifiers to reach the bare
// type name, e.g. "*[]TargetAgentDto" -> "TargetAgentDto".
func namedType(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.StarExpr:
		return namedType(t.X)
	case *ast.ArrayType:
		return namedType(t.Elt)
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return t.Sel.Name
	}
	return ""
}

func fieldTag(f *ast.Field) reflect.StructTag {
	if f.Tag == nil {
		return ""
	}
	v, err := strconv.Unquote(f.Tag.Value)
	if err != nil {
		return ""
	}
	return reflect.StructTag(v)
}

// tagValue splits one struct tag into its name and its options, e.g.
// `json:"limit,omitempty"` -> "limit", {"omitempty"}.
func tagValue(tag reflect.StructTag, key string) (string, map[string]bool) {
	raw, ok := tag.Lookup(key)
	if !ok {
		return "", nil
	}
	parts := strings.Split(raw, ",")
	opts := map[string]bool{}
	for _, o := range parts[1:] {
		opts[o] = true
	}
	name := parts[0]
	if name == "-" {
		name = ""
	}
	return name, opts
}

// fieldDoc returns a field's comment as flag help. oapi-codegen writes the
// description prefixed with the field name ("// Filters A list of filters"),
// which reads wrong on a --help line, so drop that prefix.
func fieldDoc(f *ast.Field, goName string) string {
	var raw string
	switch {
	case f.Doc != nil:
		raw = f.Doc.Text()
	case f.Comment != nil:
		raw = f.Comment.Text()
	default:
		return ""
	}
	doc := strings.Join(strings.Fields(strings.TrimSpace(raw)), " ")
	doc = strings.TrimPrefix(doc, goName+" ")
	return doc
}
