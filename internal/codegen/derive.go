package codegen

import (
	"fmt"
	"regexp"
	"strings"
)

// Derive builds a default recipe for an operation. The wizard presents this as
// the starting point, and `codegen regen` replays a saved one, so every default
// here is something a person can override rather than a rule.
func Derive(op *Operation, client *Client, specs *Specs, schema *ContextSchema) *Recipe {
	info := specs.For(op)

	r := &Recipe{
		Client:    op.Client,
		Operation: op.Name,
		Command: RecipeCommand{
			Path:  SuggestPath(op, client),
			Short: info.Summary,
		},
		Response: RecipeResponse{Render: RenderBody},
	}
	if info.Description != "" && info.Description != info.Summary {
		r.Command.Long = info.Description
	}
	r.Command.Var = CommandVar(r.Command.Path)
	r.Command.File = r.Command.Use() + ".go"

	for _, p := range op.PathParams {
		r.Args = append(r.Args, deriveArg(p, info, schema))
	}
	for _, p := range op.Params {
		r.Args = append(r.Args, deriveArg(p, info, schema))
	}

	// An endpoint can require credentials without declaring an authorization
	// parameter, in which case the client has no field for it and the header has
	// to go through a RequestEditorFn. Only the spec knows, so this is added
	// when a vendored spec says the operation is secured.
	if _, declared := op.HeaderParam("authorization"); !declared && info.Secured {
		r.Args = append(r.Args, AuthorizationArg())
	}

	switch {
	case op.Body == nil:
		r.Body.Mode = BodyNone
	case op.Body.Constructable() && !requiresRawBody(op.Body, op.Client):
		r.Body.Mode = BodyTyped
		for _, f := range op.Body.Fields {
			field := deriveBodyField(f, schema)
			// An optional field no flag can carry is dropped rather than
			// forcing the whole body to raw JSON for its sake.
			if !KindFor(f.GoType, f.Scalar, false, op.Client).Supported {
				field.Skip = true
			}
			r.Body.Fields = append(r.Body.Fields, field)
		}
	default:
		// A free-form body, or one reaching a type with unexported fields, has
		// to be handed through as raw JSON.
		r.Body.Mode = BodyRaw
		r.Body.Flag = "body"
	}

	if resp, ok := op.Success(); ok {
		r.Response.Field = resp.Field
	}
	resolveFlagCollisions(r)
	return r
}

// resolveFlagCollisions renames flags that would otherwise clash.
//
// Several endpoints name a body field the same as one of their path parameters
// — PUT /databases/{dbName}/documents/{key} takes a body with its own "key" —
// and registering one flag name twice makes cobra panic as the command tree is
// built. Parameters keep their name because they identify the resource; the body
// field is the one that moves.
func resolveFlagCollisions(r *Recipe) {
	seen := map[string]bool{}
	claim := func(name string) bool {
		if name == "" || seen[name] {
			return false
		}
		seen[name] = true
		return true
	}

	for i := range r.Args {
		arg := &r.Args[i]
		if arg.Skip {
			continue
		}
		if claim(arg.CLI) {
			continue
		}
		arg.CLI = disambiguate(seen, string(arg.Source), arg.Name)
		arg.Var = LocalVar(arg.CLI)
	}
	for i := range r.Body.Fields {
		f := &r.Body.Fields[i]
		if f.Skip {
			continue
		}
		if claim(f.CLI) {
			continue
		}
		f.CLI = disambiguate(seen, "body", f.Name)
		f.Var = LocalVar(f.CLI)
	}
}

// disambiguate builds an unused flag name by prefixing the value's origin, e.g.
// a body field named "key" becomes --bodyKey.
func disambiguate(seen map[string]bool, prefix, name string) string {
	candidate := prefix + strings.ToUpper(name[:1]) + name[1:]
	for n := 2; seen[candidate]; n++ {
		candidate = fmt.Sprintf("%s%s%d", prefix, strings.ToUpper(name[:1])+name[1:], n)
	}
	seen[candidate] = true
	return candidate
}

// AuthorizationArg is the synthetic authorization header: --token, resolved
// from the flag, then the JWT in context, then a prompt, and sent through a
// RequestEditorFn.
func AuthorizationArg() RecipeArg {
	return RecipeArg{
		Source:      InHeader,
		Name:        "authorization",
		CLI:         "token",
		Var:         "jwt",
		Type:        "string",
		Description: "authorization token",
		Resolve:     []string{FromFlag, FromContext, FromPrompt},
		Context:     "Jwt",
		Persist:     true,
		Prompt:      "authorization token",
		Via:         ViaEditor,
	}
}

// requiresRawBody reports whether a body has a required field no flag can
// carry — a list of UUIDs, a nested DTO, a map. Such a body can only be sent as
// raw JSON, so the command takes a --body flag instead of a flag per field.
func requiresRawBody(body *Body, client string) bool {
	for _, f := range body.Fields {
		if f.Optional {
			continue
		}
		if !KindFor(f.GoType, f.Scalar, false, client).Supported {
			return true
		}
	}
	return false
}

func deriveArg(p Param, info OpInfo, schema *ContextSchema) RecipeArg {
	// The spec is the better source: oapi-codegen emits no doc comment for a
	// path parameter, so p.Doc is always empty for those.
	arg := RecipeArg{
		Source:      p.In,
		Name:        p.Name,
		GoName:      p.GoName,
		CLI:         p.Name,
		Type:        p.GoType,
		Description: firstNonEmpty(info.Param(p.In, p.Name), p.Doc),
		Prompt:      p.Name,
		Resolve:     []string{FromFlag, FromPrompt},
	}

	// The authorization header is the JWT the whole CLI already keeps in
	// context, and every existing command spells its flag --token.
	if p.In == InHeader && strings.EqualFold(p.Name, "authorization") {
		arg.CLI = "token"
		arg.Var = "jwt"
		arg.Context = "Jwt"
		arg.Description = firstNonEmpty(arg.Description, "authorization token")
		arg.Prompt = "authorization token"
		arg.Persist = true
		arg.Via = ViaParams
		arg.Resolve = []string{FromFlag, FromContext, FromPrompt}
		return arg
	}
	if p.In == InHeader {
		arg.Via = ViaParams
	}

	arg.Var = LocalVar(arg.CLI)

	// A parameter the context already remembers should be read from it, which
	// is what makes `--orgName` optional on a second command.
	if f, ok := schema.Field(p.Name); ok {
		arg.Context = f.GoName
		arg.Persist = true
		arg.Resolve = []string{FromFlag, FromContext, FromPrompt}
	}

	// Optional query parameters are sent only when asked for, so they get a
	// flag and nothing else; prompting for a page number nobody set would be
	// noise. Whether an unset flag omits the parameter or sends a default is
	// the caller's call, so default to leaving it out and let the wizard ask.
	if p.In == InQuery && p.Optional {
		arg.Resolve = []string{FromFlag}
		arg.Prompt = ""
		arg.Unset = UnsetOmit
	}
	return arg
}

func deriveBodyField(f Field, schema *ContextSchema) RecipeBodyField {
	out := RecipeBodyField{
		Field:       f.GoName,
		Name:        f.JSONName,
		CLI:         f.JSONName,
		Var:         LocalVar(f.JSONName),
		Type:        f.GoType,
		Description: f.Doc,
		Optional:    f.Optional,
		Prompt:      f.JSONName,
		Resolve:     []string{FromFlag, FromPrompt},
	}
	if f.Optional {
		out.Resolve = []string{FromFlag}
		out.Prompt = ""
		out.Unset = UnsetOmit
	}
	if ctxField, ok := schema.Field(f.JSONName); ok {
		out.Context = ctxField.GoName
		out.Resolve = []string{FromFlag, FromContext, FromPrompt}
	}
	return out
}

// SuggestPath proposes where a command belongs, from the endpoint's path.
//
// The specs declare no operationIds, so the path is all there is to go on: drop
// the version, the path parameters and the ubiquitous org/{orgName} prefix, and
// what remains reads as the command path. A method verb is appended only when
// the same path carries more than one method, because that is exactly when the
// noun alone would be ambiguous.
func SuggestPath(op *Operation, client *Client) []string {
	path := CommandNouns(op)

	if client != nil && len(siblingMethods(op, client)) > 1 {
		if verb := MethodVerb(op.Method); verb != "" {
			path = append(path, verb)
		}
	}
	if len(path) == 1 {
		if verb := MethodVerb(op.Method); verb != "" {
			path = append(path, verb)
		}
	}
	return path
}

// CommandNouns is the command path an endpoint reduces to before any verb is
// added: the service, then the path's literal segments.
func CommandNouns(op *Operation) []string {
	path := []string{Service(op.Client)}

	segments := strings.Split(strings.Trim(op.Path, "/"), "/")
	for i, seg := range segments {
		switch {
		case seg == "" || strings.HasPrefix(seg, "{"):
			continue
		case strings.HasPrefix(seg, ":"):
			// Three high5 paths carry an Express-style ":secretKey" segment the
			// spec never converted to an OpenAPI path parameter, so the client
			// treats it as a literal. Keep it out of the command path; Emit
			// warns that the endpoint is not actually parameterizable.
			continue
		case seg == "v1" || seg == "v2":
			continue
		case seg == "org" && i+1 < len(segments) && strings.HasPrefix(segments[i+1], "{"):
			// Practically every endpoint is org-scoped; carrying "org" into
			// every command path would say nothing.
			continue
		}
		path = append(path, KebabCase(seg))
	}
	return path
}

// siblingMethods returns the HTTP methods the client exposes on one path.
func siblingMethods(op *Operation, client *Client) []string {
	var out []string
	for _, other := range client.Ops {
		if other.Path == op.Path {
			out = append(out, other.Method)
		}
	}
	return out
}

// KebabCase turns a camelCase path segment into a command name, so
// /settings/general/lastVisit reads as "settings general last-visit" — the same
// shape as the hand-written "patch-password".
func KebabCase(s string) string {
	s = camelBoundary.ReplaceAllString(s, "${1}-${2}")
	s = nonKebab.ReplaceAllString(s, "-")
	s = strings.ToLower(s)
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return strings.Trim(s, "-")
}

// nonKebab matches anything that cannot appear in a command name.
var nonKebab = regexp.MustCompile(`[^A-Za-z0-9-]+`)

// MethodVerb is the command name an HTTP method reads as.
func MethodVerb(method string) string {
	switch strings.ToUpper(method) {
	case "GET":
		return "get"
	case "POST":
		return "create"
	case "PUT":
		return "update"
	case "PATCH":
		return "patch"
	case "DELETE":
		return "delete"
	}
	return ""
}

// CommandVar is the variable name a generated command is declared as, following
// the repo's <verb><Noun>Command convention.
func CommandVar(path []string) string {
	if len(path) == 0 {
		return "generatedCommand"
	}
	// Skip the service segment: the package already names the service.
	parts := path
	if len(parts) > 1 {
		parts = parts[1:]
	}

	var sb strings.Builder
	for i, p := range parts {
		p = nonIdent.ReplaceAllString(p, " ")
		for j, word := range strings.Fields(p) {
			if i == 0 && j == 0 {
				sb.WriteString(strings.ToLower(word[:1]) + word[1:])
				continue
			}
			sb.WriteString(strings.ToUpper(word[:1]) + strings.ToLower(word[1:]))
		}
	}
	sb.WriteString("Command")
	return sb.String()
}

// goKeywords are the identifiers a local variable must not take.
var goKeywords = map[string]bool{
	"break": true, "case": true, "chan": true, "const": true, "continue": true,
	"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
	"func": true, "go": true, "goto": true, "if": true, "import": true,
	"interface": true, "map": true, "package": true, "range": true, "return": true,
	"select": true, "struct": true, "switch": true, "type": true, "var": true,
	// Not keywords, but shadowing these inside a generated RunE would break it.
	"client": true, "cmd": true, "args": true, "err": true, "res": true,
	"service": true, "latest": true, "text": true, "data": true, "body": true,
	"context": true, "json": true, "command": true, "clients": true, "internal": true,
	"sql": true, "params": true,
}

// nonIdent matches anything that cannot appear in a Go identifier.
var nonIdent = regexp.MustCompile(`[^A-Za-z0-9]+`)

// LocalVar turns a flag name into a local variable name safe to declare inside
// a generated RunE.
func LocalVar(name string) string {
	name = nonIdent.ReplaceAllString(name, " ")
	var sb strings.Builder
	for i, word := range strings.Fields(name) {
		if i == 0 {
			sb.WriteString(strings.ToLower(word[:1]) + word[1:])
			continue
		}
		sb.WriteString(strings.ToUpper(word[:1]) + strings.ToLower(word[1:]))
	}
	v := sb.String()
	if v == "" {
		return "value"
	}
	if goKeywords[v] {
		return v + "Value"
	}
	return v
}
