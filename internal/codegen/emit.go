package codegen

import (
	"fmt"
	"go/format"
	"path/filepath"
	"sort"
	"strings"
)

// Generated is one command's worth of code, ready to be written or spliced into
// an existing file.
type Generated struct {
	Dir     string // package directory, relative to the repo root
	PkgName string // Go package name, e.g. "mothership"
	File    string // file name within Dir
	Path    string // Dir/File

	Source   []byte   // the complete file, for a file that does not exist yet
	Imports  []string // imports the command needs
	Flags    []string // flag registrations for the init() block
	VarBlock string   // the `var xCommand = &cobra.Command{...}` declaration

	// NewColumns are context columns the command persists into that the schema
	// does not have yet.
	NewColumns []ContextField
	// Warnings records anything the recipe asked for that could not be honoured.
	Warnings []string
}

// Emit generates the code for one recipe.
func Emit(r *Recipe, op *Operation, schema *ContextSchema, ctors map[string]string) (*Generated, error) {
	if r.Command.Service() == "" {
		return nil, fmt.Errorf("recipe has no command path")
	}

	g := &Generated{
		PkgName: r.Command.Service(),
		File:    r.Command.File,
	}
	if g.File == "" {
		g.File = r.Command.Use() + ".go"
	}
	g.Dir = filepath.Join("internal", "cli", g.PkgName)
	g.Path = filepath.Join(g.Dir, g.File)

	e := &emitter{
		recipe: r,
		op:     op,
		schema: schema,
		client: op.Client,
		gen:    g,
		ctor:   constructorName(op.Client, ctors),
	}
	if !op.Callable() {
		return nil, fmt.Errorf("%s %s is not callable: %s; cmd/patcher needs to flatten that union", op.Method, op.Path, op.Blocked)
	}
	e.need("github.com/spf13/cobra")

	for _, seg := range strings.Split(op.Path, "/") {
		if strings.HasPrefix(seg, ":") {
			e.warn("the endpoint path contains a literal %q segment: the spec never declared it as a path parameter, so the client hardcodes it and the command cannot vary it", seg)
		}
	}

	body, err := e.runE()
	if err != nil {
		return nil, err
	}
	if !e.sendsAuthorization() {
		e.warn("nothing sends an authorization header: the operation declares no such parameter and the recipe adds none. If the endpoint needs credentials, add an arg with source: header, name: authorization, via: editor")
	}

	g.Flags = e.flags
	g.VarBlock = e.varBlock(body)
	g.Imports = e.importList()

	source, err := e.file()
	if err != nil {
		return nil, err
	}
	g.Source = source
	return g, nil
}

// emitter accumulates the lines, flags and imports of one generated command.
type emitter struct {
	recipe *Recipe
	op     *Operation
	schema *ContextSchema
	client string
	gen    *Generated
	ctor   string

	lines   []string
	flags   []string
	imports map[string]bool
}

func (e *emitter) line(format string, args ...any) {
	e.lines = append(e.lines, fmt.Sprintf(format, args...))
}

func (e *emitter) blank() {
	if len(e.lines) > 0 && e.lines[len(e.lines)-1] != "" {
		e.lines = append(e.lines, "")
	}
}

// errCheck emits the repo's ubiquitous error check.
func (e *emitter) errCheck() {
	e.line("if err != nil {")
	e.line("return err")
	e.line("}")
}

func (e *emitter) need(paths ...string) {
	if e.imports == nil {
		e.imports = map[string]bool{}
	}
	for _, p := range paths {
		e.imports[p] = true
	}
}

func (e *emitter) warn(format string, args ...any) {
	e.gen.Warnings = append(e.gen.Warnings, fmt.Sprintf(format, args...))
}

func (e *emitter) addFlag(cli, shorthand, flagType, def, description string) {
	target := e.recipe.Command.Var
	if shorthand != "" {
		e.flags = append(e.flags, fmt.Sprintf("%s.Flags().%sP(%q, %q, %s, %q)",
			target, flagType, cli, shorthand, def, description))
		return
	}
	e.flags = append(e.flags, fmt.Sprintf("%s.Flags().%s(%q, %s, %q)",
		target, flagType, cli, def, description))
}

// runE builds the body of the command's RunE function.
func (e *emitter) runE() ([]string, error) {
	e.need("context", "github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command")

	e.line("client, err := clients.%s(cmd)", e.ctor)
	e.errCheck()
	e.need("github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients")

	if e.usesContext() {
		e.blank()
		e.line("service := internal.New()")
		e.line("latest, err := service.GetLatest()")
		e.errCheck()
		e.need("github.com/moovit-sp-gmbh/hcloud-cli-repl/internal")
	}

	callArgs, err := e.resolveArgs()
	if err != nil {
		return nil, err
	}
	if err := e.resolveEditorHeaders(); err != nil {
		return nil, err
	}

	paramsExpr, err := e.buildParams()
	if err != nil {
		return nil, err
	}
	bodyExpr, err := e.buildBody()
	if err != nil {
		return nil, err
	}

	e.blank()
	e.emitCall(callArgs, paramsExpr, bodyExpr)
	e.blank()
	e.emitRender()
	e.emitPersist()
	e.blank()
	e.line("command.SetResult(cmd, text)")
	e.blank()
	e.line("return nil")

	return e.lines, nil
}

// sendsAuthorization reports whether the command sends credentials, by either
// route.
func (e *emitter) sendsAuthorization() bool {
	for _, a := range e.recipe.Args {
		if !a.Skip && a.Source == InHeader && strings.EqualFold(a.Name, "authorization") {
			return true
		}
	}
	return false
}

// usesContext reports whether the command touches the context table at all.
func (e *emitter) usesContext() bool {
	for _, a := range e.recipe.Args {
		if a.Skip {
			continue
		}
		if a.Context != "" || a.Persist {
			return true
		}
	}
	for _, f := range e.recipe.Body.Fields {
		if f.Skip {
			continue
		}
		if f.Context != "" || f.Persist {
			return true
		}
	}
	return len(e.recipe.Response.Persist) > 0
}

// resolveArgs emits the resolution of every path parameter and returns the
// expressions to pass positionally at the call site.
func (e *emitter) resolveArgs() ([]string, error) {
	var pathArgs []string

	for _, p := range e.op.PathParams {
		arg, ok := e.argFor(InPath, p.Name)
		if !ok || arg.Skip {
			return nil, fmt.Errorf("path parameter %q must be mapped: the endpoint cannot be called without it", p.Name)
		}
		expr, err := e.resolveOne(argSpec{
			cli:         arg.CLI,
			shorthand:   arg.Shorthand,
			variable:    arg.Var,
			goType:      p.GoType,
			scalar:      p.Scalar,
			description: firstNonEmpty(arg.Description, p.Doc),
			resolve:     arg.Resolve,
			contextName: arg.Context,
			persist:     arg.Persist,
			prompt:      arg.Prompt,
			required:    true,
		})
		if err != nil {
			return nil, fmt.Errorf("path parameter %q: %w", p.Name, err)
		}
		pathArgs = append(pathArgs, expr)
	}
	return pathArgs, nil
}

// resolveEditorHeaders resolves headers sent through a RequestEditorFn rather
// than a typed Params field.
//
// Some endpoints need credentials without declaring an authorization parameter
// — DeleteV1Agent is one — so the client has no field to put the header in and
// the only way to send it is an editor closure, which is what
// internal/cli/mothership/agent.go does by hand.
func (e *emitter) resolveEditorHeaders() error {
	for _, arg := range e.recipe.Args {
		if arg.Skip || arg.Source != InHeader || arg.Via != ViaEditor {
			continue
		}
		// A header the operation declares is resolved with the params struct.
		if _, declared := e.op.HeaderParam(arg.Name); declared {
			continue
		}
		goType := arg.Type
		if goType == "" {
			goType = "string"
		}
		if _, err := e.resolveOne(argSpec{
			cli:         arg.CLI,
			shorthand:   arg.Shorthand,
			variable:    arg.Var,
			goType:      goType,
			description: arg.Description,
			resolve:     arg.Resolve,
			contextName: arg.Context,
			persist:     arg.Persist,
			prompt:      arg.Prompt,
			required:    true,
		}); err != nil {
			return fmt.Errorf("header %q: %w", arg.Name, err)
		}
	}
	return nil
}

// argSpec is one argument to resolve, flattened from either a parameter or a
// body field so both go through the same emitter.
type argSpec struct {
	cli         string
	shorthand   string
	variable    string
	goType      string
	scalar      string
	description string
	resolve     []string
	contextName string
	persist     bool
	prompt      string
	required    bool
	def         string
}

// resolveOne emits the statements that produce one value and returns the
// expression to use at the call site.
func (e *emitter) resolveOne(s argSpec) (string, error) {
	wantPrompt := hasSource(s.resolve, FromPrompt)
	wantContext := hasSource(s.resolve, FromContext) && s.contextName != ""

	kind := KindFor(s.goType, s.scalar, wantPrompt, e.client)
	if !kind.Supported {
		return "", fmt.Errorf("type %s cannot be carried on a flag", s.goType)
	}

	variable := s.variable
	if variable == "" {
		variable = LocalVar(s.cli)
	}
	literal, err := kind.DefaultLiteral(s.def)
	if err != nil {
		return "", fmt.Errorf("--%s default: %w", s.cli, err)
	}
	e.addFlag(s.cli, s.shorthand, kind.Flag, literal, s.description)

	e.blank()
	switch {
	case wantContext && kind.Flag == "String":
		// The flag, then whatever the context remembers, then a prompt — the
		// chain internal/command/args.go exists to express.
		chain := []string{fmt.Sprintf("command.StringFlagIfSet(%q)", s.cli)}
		chain = append(chain, fmt.Sprintf("command.ContextIfValidString(latest.%s)", s.contextName))
		if wantPrompt {
			chain = append(chain, fmt.Sprintf("command.PromptString(%q)", firstNonEmpty(s.prompt, s.cli)))
		}
		e.line("%s, err := command.GetArg(cmd, %s)", variable, strings.Join(chain, ", "))
		e.errCheck()

	case wantContext:
		// Context columns are text and GetArg's providers only exist for
		// strings, so there is nothing to read a non-string value out of. Fall
		// back to the flag rather than emitting a conversion that could not be
		// round-tripped.
		e.warn("%s is a %s, which the context table cannot hold: it is resolved from its flag only", s.cli, s.goType)
		e.line("%s, err := cmd.Flags().Get%s(%q)", variable, kind.Flag, s.cli)
		e.errCheck()

	case wantPrompt && kind.Helper != "":
		prompt := firstNonEmpty(s.prompt, s.cli)
		if prompt == s.cli {
			e.line("%s, err := command.FlagOrPrompt%s(cmd, %q)", variable, kind.Helper, s.cli)
		} else {
			e.line("%s, err := command.FlagOrDifferentPrompt%s(cmd, %q, %q)", variable, kind.Helper, s.cli, prompt)
		}
		e.errCheck()

	default:
		if wantPrompt {
			e.warn("%s is a %s, which internal/command/args.go cannot prompt for; the flag is required", s.cli, s.goType)
		}
		e.line("%s, err := cmd.Flags().Get%s(%q)", variable, kind.Flag, s.cli)
		e.errCheck()
	}

	e.need(kind.Imports...)

	expr := variable
	if kind.Parse {
		// A UUID has to be parsed, and the parse can fail, so it cannot be an
		// inline conversion.
		parsed := variable + "Parsed"
		e.line("%s, err := uuid.Parse(%s)", parsed, variable)
		e.errCheck()
		expr = parsed
	} else {
		expr = kind.Convert(variable)
	}

	// Only a value the context can hold is written back; asString would
	// otherwise stringify something that could not be read out again.
	if s.persist && s.contextName != "" && (kind.Flag == "String" || s.scalar == "string") {
		e.line("latest.%s = sql.NullString{Valid: true, String: %s}", s.contextName, e.asString(variable, s.goType, s.scalar))
		e.need("database/sql")
	}
	return expr, nil
}

// asString renders a resolved variable as the string a context column stores.
func (e *emitter) asString(variable, goType, scalar string) string {
	base := strings.TrimPrefix(goType, "*")
	switch {
	case base == "string":
		return variable
	case scalar == "string", base == "[]byte":
		return fmt.Sprintf("string(%s)", variable)
	case base == "openapi_types.Email", base == "types.Email":
		return variable
	default:
		e.need("fmt")
		return fmt.Sprintf("fmt.Sprint(%s)", variable)
	}
}

// buildParams emits the operation's Params struct and returns the expression to
// pass, or "nil" when nothing maps onto it.
func (e *emitter) buildParams() (string, error) {
	if e.op.ParamsType == "" {
		return "", nil
	}

	type mapped struct {
		param Param
		arg   RecipeArg
	}
	var required, optional []mapped
	for _, p := range e.op.Params {
		arg, ok := e.argFor(p.In, p.Name)
		if !ok || arg.Skip {
			continue
		}
		if arg.Via == ViaEditor {
			continue
		}
		if p.Optional {
			optional = append(optional, mapped{p, arg})
		} else {
			required = append(required, mapped{p, arg})
		}
	}
	if len(required) == 0 && len(optional) == 0 {
		return "nil", nil
	}
	e.needClient()

	// Required fields go in the literal; optional ones are set only when their
	// flag was given, so an unset --limit does not override a server default.
	var inline []string
	for _, m := range required {
		expr, err := e.resolveOne(e.specForParam(m.param, m.arg))
		if err != nil {
			return "", fmt.Errorf("%s parameter %s: %w", m.param.In, m.param.Name, err)
		}
		inline = append(inline, fmt.Sprintf("%s: %s,", m.param.GoName, expr))
	}

	// With nothing optional to guard, the struct reads better inline at the
	// call site, which is how internal/cli/mothership/stats.go writes it.
	if len(optional) == 0 {
		if len(inline) == 0 {
			return fmt.Sprintf("&%s.%s{}", e.client, e.op.ParamsType), nil
		}
		return fmt.Sprintf("&%s.%s{\n%s\n}", e.client, e.op.ParamsType, strings.Join(inline, "\n")), nil
	}

	e.blank()
	if len(inline) == 0 {
		e.line("params := &%s.%s{}", e.client, e.op.ParamsType)
	} else {
		e.line("params := &%s.%s{", e.client, e.op.ParamsType)
		for _, l := range inline {
			e.line("%s", l)
		}
		e.line("}")
	}

	for _, m := range optional {
		if err := e.emitOptional(optionalSpec{
			target:      "params",
			field:       m.param.GoName,
			goType:      m.param.GoType,
			scalar:      m.param.Scalar,
			cli:         m.arg.CLI,
			shorthand:   m.arg.Shorthand,
			variable:    m.arg.Var,
			description: firstNonEmpty(m.arg.Description, m.param.Doc),
			unset:       m.arg.Unset,
			def:         m.arg.Default,
		}); err != nil {
			return "", err
		}
	}
	return "params", nil
}

// optionalSpec is one optional value to wire onto a struct field: a query
// parameter onto the Params struct, or a body field onto the body.
type optionalSpec struct {
	target      string // "params" or "body"
	field       string // Go field name on that struct
	goType      string
	scalar      string
	cli         string
	shorthand   string
	variable    string
	description string
	unset       string // omit or default
	def         string
}

// emitOptional wires one optional value onto its struct.
//
// With unset "omit" the field is only set when the flag was given, so the
// server keeps applying its own default. With unset "default" the field is
// always sent, carrying the recipe's default when the flag was not given —
// which is the only way to make a command state a value explicitly rather than
// deferring to the API.
func (e *emitter) emitOptional(s optionalSpec) error {
	kind := KindFor(s.goType, s.scalar, false, e.client)
	if !kind.Supported {
		e.warn("optional %s field %s skipped: type %s cannot be carried on a flag", s.target, s.field, s.goType)
		return nil
	}

	literal, err := kind.DefaultLiteral(s.def)
	if err != nil {
		return fmt.Errorf("--%s default: %w", s.cli, err)
	}

	variable := firstNonEmpty(s.variable, LocalVar(s.cli))
	e.addFlag(s.cli, s.shorthand, kind.Flag, literal, s.description)
	e.need(kind.Imports...)

	// A pointer field takes the address; oapi-codegen renders optional fields
	// as pointers, but a slice or map optional-by-omitempty is not one.
	assign := "%s.%s = &%s"
	if !strings.HasPrefix(s.goType, "*") {
		assign = "%s.%s = %s"
	}

	// convert emits whatever turns the flag's value into the field's type. A
	// UUID has to be parsed rather than cast, and the parse can fail, so it
	// cannot be folded into the assignment.
	convert := func(variable string) string {
		if kind.Parse {
			parsed := variable + "Parsed"
			e.line("%s, err := uuid.Parse(%s)", parsed, variable)
			e.errCheck()
			return parsed
		}
		if kind.Wrap != "" {
			converted := variable + "Value"
			e.line("%s := %s", converted, kind.Convert(variable))
			return converted
		}
		return variable
	}

	e.blank()
	if s.unset == UnsetDefault {
		e.line("%s, err := cmd.Flags().Get%s(%q)", variable, kind.Flag, s.cli)
		e.errCheck()
		e.line(assign, s.target, s.field, convert(variable))
		return nil
	}

	e.line("if cmd.Flags().Changed(%q) {", s.cli)
	e.line("%s, err := cmd.Flags().Get%s(%q)", variable, kind.Flag, s.cli)
	e.errCheck()
	variable = convert(variable)
	e.line(assign, s.target, s.field, variable)
	e.line("}")
	return nil
}

func (e *emitter) specForParam(p Param, arg RecipeArg) argSpec {
	return argSpec{
		cli:         arg.CLI,
		shorthand:   arg.Shorthand,
		variable:    arg.Var,
		goType:      p.GoType,
		scalar:      p.Scalar,
		description: firstNonEmpty(arg.Description, p.Doc),
		resolve:     arg.Resolve,
		contextName: arg.Context,
		persist:     arg.Persist,
		prompt:      arg.Prompt,
		def:         arg.Default,
		required:    !p.Optional,
	}
}

// buildBody emits the request body and returns the expression to pass.
func (e *emitter) buildBody() (string, error) {
	if e.op.Body == nil {
		return "", nil
	}
	e.needClient()

	switch e.recipe.Body.Mode {
	case BodyNone, "":
		return "", fmt.Errorf("operation %s takes a body but the recipe says mode: none", e.op.Name)

	case BodyRaw:
		// A free-form body, or one reaching a type with unexported fields, is
		// handed through as raw JSON — the same escape hatch
		// internal/cli/mothership/connect.go uses.
		flag := firstNonEmpty(e.recipe.Body.Flag, "body")
		e.addFlag(flag, "", "String", `"{}"`, "request body as JSON")
		e.blank()
		e.line("bodyJSON, err := cmd.Flags().GetString(%q)", flag)
		e.errCheck()
		e.line("var body %s.%s", e.client, e.op.Body.TypeName)
		e.line("err = json.Unmarshal([]byte(bodyJSON), &body)")
		e.errCheck()
		e.need("encoding/json")
		return "body", nil

	case BodyTyped:
		byName := map[string]Field{}
		for _, f := range e.op.Body.Fields {
			byName[f.GoName] = f
		}

		var inline []string
		var optional []RecipeBodyField
		for _, rf := range e.recipe.Body.Fields {
			if rf.Skip {
				continue
			}
			field, ok := byName[rf.Field]
			if !ok {
				e.warn("body field %s is not on %s any more; skipped", rf.Field, e.op.Body.DTO)
				continue
			}
			if field.Optional {
				optional = append(optional, rf)
				continue
			}
			expr, err := e.resolveOne(argSpec{
				cli:         rf.CLI,
				shorthand:   rf.Shorthand,
				variable:    rf.Var,
				goType:      field.GoType,
				scalar:      field.Scalar,
				description: firstNonEmpty(rf.Description, field.Doc),
				resolve:     rf.Resolve,
				contextName: rf.Context,
				persist:     rf.Persist,
				prompt:      rf.Prompt,
				def:         rf.Default,
				required:    true,
			})
			if err != nil {
				return "", fmt.Errorf("body field %s: %w", rf.Field, err)
			}
			inline = append(inline, fmt.Sprintf("%s: %s,", field.GoName, expr))
		}

		e.blank()
		if len(inline) == 0 {
			e.line("body := %s.%s{}", e.client, e.op.Body.TypeName)
		} else {
			e.line("body := %s.%s{", e.client, e.op.Body.TypeName)
			for _, l := range inline {
				e.line("%s", l)
			}
			e.line("}")
		}

		for _, rf := range optional {
			field := byName[rf.Field]
			if err := e.emitOptional(optionalSpec{
				target:      "body",
				field:       field.GoName,
				goType:      field.GoType,
				scalar:      field.Scalar,
				cli:         rf.CLI,
				shorthand:   rf.Shorthand,
				variable:    rf.Var,
				description: firstNonEmpty(rf.Description, field.Doc),
				unset:       rf.Unset,
				def:         rf.Default,
			}); err != nil {
				return "", fmt.Errorf("body field %s: %w", rf.Field, err)
			}
		}
		return "body", nil
	}
	return "", fmt.Errorf("unknown body mode %q", e.recipe.Body.Mode)
}

// emitCall emits the client call, including any header carried by a
// RequestEditorFn rather than by the typed params struct.
func (e *emitter) emitCall(pathArgs []string, paramsExpr, bodyExpr string) {
	call := []string{"context.TODO()"}
	call = append(call, pathArgs...)
	if paramsExpr != "" {
		call = append(call, paramsExpr)
	}
	if bodyExpr != "" {
		call = append(call, bodyExpr)
	}

	var editors []RecipeArg
	for _, a := range e.recipe.Args {
		if a.Source == InHeader && a.Via == ViaEditor && !a.Skip {
			editors = append(editors, a)
		}
	}

	if len(editors) == 0 {
		e.line("res, err := client.%sWithResponse(%s)", e.op.Name, strings.Join(call, ", "))
		e.errCheck()
		return
	}

	e.need("net/http")
	e.line("res, err := client.%sWithResponse(%s, func(_ context.Context, req *http.Request) error {",
		e.op.Name, strings.Join(call, ", "))
	for _, a := range editors {
		e.line("req.Header.Set(%q, %s)", a.Name, firstNonEmpty(a.Var, LocalVar(a.CLI)))
	}
	e.line("return nil")
	e.line("})")
	e.errCheck()
}

// emitRender turns the response into the command's result.
func (e *emitter) emitRender() {
	e.need("encoding/json")

	if e.recipe.Response.Render == RenderTyped && e.recipe.Response.Field != "" {
		field := e.recipe.Response.Field
		e.line("var text command.Text")
		e.line("if res.%s != nil {", field)
		e.line("text, err = command.JSON(res.%s)", field)
		e.errCheck()
		e.line("} else {")
		e.line("var data any")
		e.line("if err := json.Unmarshal(res.Body, &data); err != nil {")
		e.line("return err")
		e.line("}")
		e.line("text, err = command.JSON(data)")
		e.errCheck()
		e.line("}")
		return
	}

	// Unmarshalling res.Body and re-encoding prints both successes and errors
	// as readable JSON, which is what internal/cli/mothership/stats.go does.
	e.line("var data any")
	e.line("err = json.Unmarshal(res.Body, &data)")
	e.errCheck()
	e.line("text, err := command.JSON(data)")
	e.errCheck()
}

// emitPersist writes the chosen parts of a successful response into the context
// table and pushes the row.
func (e *emitter) emitPersist() {
	if !e.usesContext() {
		return
	}
	e.blank()
	e.line("if res.StatusCode() >= 200 && res.StatusCode() < 300 {")

	for _, p := range e.recipe.Response.Persist {
		if name, ok := p.Header(); ok {
			e.line("latest.%s = sql.NullString{Valid: true, String: res.HTTPResponse.Header.Get(%q)}", p.Context, name)
			e.need("database/sql")
			continue
		}
		respField, dtoField, ok := p.Body()
		if !ok {
			e.warn("cannot persist %q: expected header:<name> or <responseField>:<dtoField>", p.From)
			continue
		}
		expr, guard, err := e.persistExpr(respField, dtoField)
		if err != nil {
			e.warn("cannot persist %q: %v", p.From, err)
			continue
		}
		// The typed response field is nil on any other status, and a bare
		// dereference is how internal/cli/idp/user.go panics.
		e.line("if %s {", guard)
		e.line("latest.%s = sql.NullString{Valid: true, String: %s}", p.Context, expr)
		e.line("}")
		e.need("database/sql")
	}

	e.line("service.Push(latest)")
	e.line("}")
}

// persistExpr renders the expression that reads one field off a typed response,
// plus the nil guard it needs.
func (e *emitter) persistExpr(respField, dtoField string) (expr, guard string, err error) {
	var resp *Response
	for i := range e.op.Responses {
		if e.op.Responses[i].Field == respField {
			resp = &e.op.Responses[i]
		}
	}
	if resp == nil {
		return "", "", fmt.Errorf("the operation has no %s response", respField)
	}
	if resp.List {
		return "", "", fmt.Errorf("%s is an array and has no addressable fields", respField)
	}

	var field *Field
	for i := range resp.Fields {
		if resp.Fields[i].GoName == dtoField {
			field = &resp.Fields[i]
		}
	}
	if field == nil {
		return "", "", fmt.Errorf("%s has no field %s", resp.DTO, dtoField)
	}

	access := fmt.Sprintf("res.%s.%s", respField, dtoField)
	guard = fmt.Sprintf("res.%s != nil", respField)
	if field.Optional && strings.HasPrefix(field.GoType, "*") {
		guard += fmt.Sprintf(" && %s != nil", access)
		access = "*" + access
	}
	return e.asString(access, strings.TrimPrefix(field.GoType, "*"), field.Scalar), guard, nil
}

func (e *emitter) needClient() {
	e.need("github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/" + e.client)
}

// varBlock renders the cobra.Command declaration around a RunE body.
func (e *emitter) varBlock(body []string) string {
	var sb strings.Builder
	c := e.recipe.Command

	fmt.Fprintf(&sb, "var %s = &cobra.Command{\n", c.Var)
	fmt.Fprintf(&sb, "Use: %q,\n", c.Use())
	if c.Short != "" {
		fmt.Fprintf(&sb, "Short: %q,\n", c.Short)
	}
	if c.Long != "" {
		fmt.Fprintf(&sb, "Long: %q,\n", c.Long)
	}
	sb.WriteString("RunE: func(cmd *cobra.Command, args []string) error {\n")
	for _, l := range body {
		sb.WriteString(l)
		sb.WriteString("\n")
	}
	sb.WriteString("},\n}\n")
	return sb.String()
}

// file renders the complete source file for a command that gets its own file.
func (e *emitter) file() ([]byte, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "package %s\n\n", e.gen.PkgName)

	sb.WriteString("import (\n")
	for _, group := range groupImports(e.importList()) {
		for _, imp := range group {
			fmt.Fprintf(&sb, "%q\n", imp)
		}
		sb.WriteString("\n")
	}
	sb.WriteString(")\n\n")

	if len(e.flags) > 0 {
		sb.WriteString("func init() {\n")
		for _, f := range e.flags {
			sb.WriteString(f)
			sb.WriteString("\n")
		}
		sb.WriteString("}\n\n")
	}
	sb.WriteString(e.gen.VarBlock)

	formatted, err := format.Source([]byte(sb.String()))
	if err != nil {
		return nil, fmt.Errorf("generated code does not parse: %w\n%s", err, numbered(sb.String()))
	}
	return formatted, nil
}

func (e *emitter) importList() []string {
	out := make([]string, 0, len(e.imports))
	for imp := range e.imports {
		out = append(out, imp)
	}
	sort.Strings(out)
	return out
}

// groupImports splits imports into the standard-library group and the rest, the
// way gofmt-ed files in this repo are arranged.
func groupImports(imports []string) [][]string {
	var std, other []string
	for _, imp := range imports {
		if strings.Contains(strings.SplitN(imp, "/", 2)[0], ".") {
			other = append(other, imp)
		} else {
			std = append(std, imp)
		}
	}
	var out [][]string
	if len(std) > 0 {
		out = append(out, std)
	}
	if len(other) > 0 {
		out = append(out, other)
	}
	return out
}

// constructorName is the clients.New*Client function for a client package.
// Existing names win, because idp's is spelled NewIDPClient rather than
// NewIdpClient.
func constructorName(client string, ctors map[string]string) string {
	for name, pkg := range ctors {
		if pkg == client {
			return name
		}
	}
	service := Service(client)
	return "New" + strings.ToUpper(service[:1]) + service[1:] + "Client"
}

// argFor finds the recipe argument for one API parameter.
func (e *emitter) argFor(in Location, name string) (RecipeArg, bool) {
	for _, a := range e.recipe.Args {
		if a.Source == in && strings.EqualFold(a.Name, name) {
			return a, true
		}
	}
	return RecipeArg{}, false
}

func hasSource(resolve []string, want string) bool {
	for _, r := range resolve {
		if r == want {
			return true
		}
	}
	return false
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func numbered(src string) string {
	var sb strings.Builder
	for i, line := range strings.Split(src, "\n") {
		fmt.Fprintf(&sb, "%4d  %s\n", i+1, line)
	}
	return sb.String()
}
