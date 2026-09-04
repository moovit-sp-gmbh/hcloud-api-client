package codegen

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	goIdent = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	// Flag names mirror the API's own spelling, the way --orgName and
	// --cpuUtilization already do, so camelCase and the snake_case the OAuth
	// endpoints use both have to pass.
	flagName    = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_-]*$`)
	pathSegment = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)
	fileName    = regexp.MustCompile(`^[A-Za-z0-9_-]+\.go$`)
)

// ValidGoIdent reports whether s can be a Go identifier.
func ValidGoIdent(s string) bool { return goIdent.MatchString(s) }

// ValidFlagName reports whether s can be a cobra flag name.
func ValidFlagName(s string) bool { return flagName.MatchString(s) }

// ValidPathSegment reports whether s can be a command name.
func ValidPathSegment(s string) bool { return pathSegment.MatchString(s) }

// ValidFileName reports whether s can be a Go file name in a command package.
func ValidFileName(s string) bool { return fileName.MatchString(s) }

// Validate checks a recipe before anything is written.
//
// Recipes are hand-editable and replayed non-interactively, so a typo in one
// would otherwise surface as a migration adding a column called "name of the
// space" or as code that does not compile.
func (r *Recipe) Validate() error {
	var problems []string
	bad := func(format string, args ...any) {
		problems = append(problems, fmt.Sprintf(format, args...))
	}

	if r.Client == "" {
		bad("client is empty")
	}
	if !ValidGoIdent(r.Operation) {
		bad("operation %q is not an identifier", r.Operation)
	}

	if len(r.Command.Path) < 2 {
		bad("command path %v needs at least a service and a command name", r.Command.Path)
	}
	for _, segment := range r.Command.Path {
		if !ValidPathSegment(segment) {
			bad("command path segment %q must be lowercase letters, digits and dashes", segment)
		}
	}
	if !ValidGoIdent(r.Command.Var) {
		bad("command variable %q is not an identifier", r.Command.Var)
	}
	if r.Command.File != "" && !ValidFileName(r.Command.File) {
		bad("command file %q must be a bare .go file name", r.Command.File)
	}

	seenFlags := map[string]string{}
	checkFlag := func(what, cli, shorthand string) {
		if !ValidFlagName(cli) {
			bad("%s flag name %q must start with a letter and hold only letters, digits, dashes and underscores", what, cli)
		}
		if other, dup := seenFlags[cli]; dup {
			bad("%s and %s both use the flag --%s", other, what, cli)
		}
		seenFlags[cli] = what
		if shorthand != "" && len(shorthand) != 1 {
			bad("%s shorthand %q must be a single character", what, shorthand)
		}
	}
	checkContext := func(what, field string, resolve []string, persist bool) {
		if field == "" {
			if hasSource(resolve, FromContext) {
				bad("%s resolves from context but names no context field", what)
			}
			if persist {
				bad("%s persists to context but names no context field", what)
			}
			return
		}
		if !ValidGoIdent(field) {
			bad("%s context field %q is not an identifier", what, field)
		}
	}
	checkUnset := func(what, unset string) {
		switch unset {
		case "", UnsetOmit, UnsetDefault:
		default:
			bad("%s unset %q must be omit or default", what, unset)
		}
	}
	checkResolve := func(what string, resolve []string) {
		if len(resolve) == 0 {
			bad("%s has no resolution sources", what)
			return
		}
		for _, source := range resolve {
			switch source {
			case FromFlag, FromContext, FromPrompt:
			default:
				bad("%s has unknown resolution source %q; expected flag, context or prompt", what, source)
			}
		}
	}

	for _, a := range r.Args {
		if a.Skip {
			continue
		}
		what := fmt.Sprintf("%s parameter %q", a.Source, a.Name)
		switch a.Source {
		case InPath, InQuery, InHeader:
		default:
			bad("%s has unknown source %q; expected path, query or header", what, a.Source)
		}
		checkFlag(what, a.CLI, a.Shorthand)
		checkResolve(what, a.Resolve)
		checkContext(what, a.Context, a.Resolve, a.Persist)
		if a.Var != "" && !ValidGoIdent(a.Var) {
			bad("%s variable %q is not an identifier", what, a.Var)
		}
		if a.Via != "" && a.Via != ViaParams && a.Via != ViaEditor {
			bad("%s via %q must be params or editor", what, a.Via)
		}
		checkUnset(what, a.Unset)
	}

	switch r.Body.Mode {
	case BodyNone, BodyTyped, BodyRaw:
	case "":
		bad("body mode is empty; expected none, typed or raw")
	default:
		bad("unknown body mode %q", r.Body.Mode)
	}
	if r.Body.Mode == BodyRaw && r.Body.Flag != "" && !ValidFlagName(r.Body.Flag) {
		bad("body flag name %q is not a valid flag name", r.Body.Flag)
	}
	for _, f := range r.Body.Fields {
		if f.Skip {
			continue
		}
		what := fmt.Sprintf("body field %q", f.Name)
		if !ValidGoIdent(f.Field) {
			bad("%s targets %q, which is not an identifier", what, f.Field)
		}
		checkFlag(what, f.CLI, f.Shorthand)
		checkResolve(what, f.Resolve)
		checkContext(what, f.Context, f.Resolve, f.Persist)
		checkUnset(what, f.Unset)
		if f.Default != "" && !f.Optional && hasSource(f.Resolve, FromPrompt) {
			bad("%s has a default but is prompted for, so the default can never apply", what)
		}
	}

	switch r.Response.Render {
	case RenderBody:
	case RenderTyped:
		if r.Response.Field == "" {
			bad("response render is typed but names no response field")
		}
	case "":
		bad("response render is empty; expected typed or body")
	default:
		bad("unknown response render %q", r.Response.Render)
	}

	for _, p := range r.Response.Persist {
		if !ValidGoIdent(p.Context) {
			bad("response persist target %q is not an identifier", p.Context)
		}
		if p.Column != "" && !ValidPathSegment(strings.ReplaceAll(p.Column, "_", "-")) {
			bad("response persist column %q is not a valid column name", p.Column)
		}
		if _, ok := p.Header(); ok {
			continue
		}
		if _, _, ok := p.Body(); !ok {
			bad("response persist source %q must be header:<name> or <responseField>:<dtoField>", p.From)
		}
	}

	if len(problems) > 0 {
		return fmt.Errorf("invalid recipe:\n  - %s", strings.Join(problems, "\n  - "))
	}
	return nil
}
