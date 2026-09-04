package codegen

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

// Wizard walks a person through turning an endpoint into a command. It asks
// only what cannot be derived, and every answer it collects lands in a recipe
// so the same command can be regenerated later without asking again.
type Wizard struct {
	in  *bufio.Reader
	out io.Writer

	Index  *Index
	Tree   *Tree
	Specs  *Specs
	Schema *ContextSchema
}

// NewWizard builds a wizard reading from in and writing to out.
func NewWizard(in io.Reader, out io.Writer, ix *Index, tree *Tree, specs *Specs, schema *ContextSchema) *Wizard {
	return &Wizard{in: bufio.NewReader(in), out: out, Index: ix, Tree: tree, Specs: specs, Schema: schema}
}

func (w *Wizard) printf(format string, args ...any) {
	fmt.Fprintf(w.out, format, args...)
}

// ask prompts for a line of input, returning def when the answer is empty.
func (w *Wizard) ask(prompt, def string) (string, error) {
	if def != "" {
		w.printf("%s [%s]: ", prompt, def)
	} else {
		w.printf("%s: ", prompt)
	}
	line, err := w.in.ReadString('\n')
	if err != nil && (err != io.EOF || line == "") {
		return "", err
	}
	answer := strings.TrimSpace(line)
	if answer == "" {
		return def, nil
	}
	return answer, nil
}

// askValid re-asks until the answer passes, so a mistyped flag or context field
// is caught while the person can still fix it rather than surfacing later as a
// column called "name of the space".
func (w *Wizard) askValid(prompt, def string, valid func(string) bool, complaint string) (string, error) {
	for {
		answer, err := w.ask(prompt, def)
		if err != nil {
			return "", err
		}
		if valid(answer) {
			return answer, nil
		}
		w.printf("  %s\n", complaint)
	}
}

// askOptional is askValid for an answer that may be left empty.
func (w *Wizard) askOptional(prompt, def string, valid func(string) bool, complaint string) (string, error) {
	return w.askValid(prompt, def, func(s string) bool {
		return s == "" || valid(s)
	}, complaint)
}

// confirm asks a yes/no question.
func (w *Wizard) confirm(prompt string, def bool) (bool, error) {
	hint := "y/N"
	if def {
		hint = "Y/n"
	}
	answer, err := w.ask(fmt.Sprintf("%s (%s)", prompt, hint), "")
	if err != nil {
		return false, err
	}
	switch strings.ToLower(answer) {
	case "":
		return def, nil
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	}
	return def, nil
}

// choose presents a numbered list and returns the chosen index.
func (w *Wizard) choose(prompt string, options []string, def int) (int, error) {
	for i, o := range options {
		w.printf("  %2d) %s\n", i+1, o)
	}
	for {
		answer, err := w.ask(prompt, strconv.Itoa(def+1))
		if err != nil {
			return 0, err
		}
		n, err := strconv.Atoi(answer)
		if err != nil || n < 1 || n > len(options) {
			w.printf("pick a number between 1 and %d\n", len(options))
			continue
		}
		return n - 1, nil
	}
}

// PickOperation asks which endpoint to generate a command for. Only operations
// with no command are offered, since that is the point.
func (w *Wizard) PickOperation(manifest *Manifest, client, filter string) (*Operation, error) {
	if client == "" {
		names := w.Index.ClientNames()
		labels := make([]string, len(names))
		for i, name := range names {
			cov := Compare(w.Index, manifest, []string{name}).Clients
			labels[i] = fmt.Sprintf("%-18s %d missing of %d", name, len(cov[0].Missing), cov[0].Total)
		}
		w.printf("\nWhich service?\n")
		i, err := w.choose("service", labels, 0)
		if err != nil {
			return nil, err
		}
		client = names[i]
	}
	client = ClientPackage(client)
	if _, ok := w.Index.Clients[client]; !ok {
		return nil, fmt.Errorf("no client package %q; have %v", client, w.Index.ClientNames())
	}

	cov := Compare(w.Index, manifest, []string{client}).Clients[0]
	var candidates []*Operation
	for _, op := range cov.Missing {
		if !op.Callable() {
			continue
		}
		if filter != "" && !matchesFilter(op, filter) {
			continue
		}
		candidates = append(candidates, op)
	}
	if len(candidates) == 0 {
		if filter != "" {
			return nil, fmt.Errorf("no unimplemented %s operation matches %q", client, filter)
		}
		return nil, fmt.Errorf("%s has no unimplemented operations left", client)
	}

	// A hundred-odd endpoints is too many to scroll, so narrow before listing.
	for len(candidates) > 30 {
		w.printf("\n%d operations in %s. Narrow by path or method.\n", len(candidates), client)
		answer, err := w.ask("filter", "")
		if err != nil {
			return nil, err
		}
		if answer == "" {
			break
		}
		var kept []*Operation
		for _, op := range candidates {
			if matchesFilter(op, answer) {
				kept = append(kept, op)
			}
		}
		if len(kept) == 0 {
			w.printf("nothing matches %q\n", answer)
			continue
		}
		candidates = kept
	}

	labels := make([]string, len(candidates))
	for i, op := range candidates {
		info := w.Specs.For(op)
		labels[i] = fmt.Sprintf("%-6s %-52s %s", op.Method, op.Path, firstNonEmpty(info.Summary, info.Tag))
	}
	w.printf("\nWhich endpoint?\n")
	i, err := w.choose("endpoint", labels, 0)
	if err != nil {
		return nil, err
	}
	return candidates[i], nil
}

func matchesFilter(op *Operation, filter string) bool {
	filter = strings.ToLower(filter)
	return strings.Contains(strings.ToLower(op.Path), filter) ||
		strings.Contains(strings.ToLower(op.Method), filter) ||
		strings.Contains(strings.ToLower(op.Name), filter)
}

// Run asks everything needed to fill in a recipe, starting from the derived
// defaults.
func (w *Wizard) Run(op *Operation, r *Recipe) error {
	if err := w.commandPath(op, r); err != nil {
		return err
	}
	if err := w.arguments(op, r); err != nil {
		return err
	}
	if err := w.authorization(op, r); err != nil {
		return err
	}
	if err := w.body(op, r); err != nil {
		return err
	}
	return w.persistence(op, r)
}

// commandPath asks where the command belongs, and reports which parts of the
// path do not exist yet.
func (w *Wizard) commandPath(op *Operation, r *Recipe) error {
	w.printf("\n--- Command path ---\n")
	w.printf("%s %s\n", op.Method, op.Path)
	if info := w.Specs.For(op); info.Tag != "" {
		w.printf("spec tag: %s\n", info.Tag)
	}
	w.printf("existing commands under %s: %s\n", r.Command.Service(), strings.Join(siblingPaths(w.Tree, r.Command.Service()), ", "))

	for {
		answer, err := w.ask("command path", strings.Join(r.Command.Path, " "))
		if err != nil {
			return err
		}
		path := strings.Fields(answer)
		if len(path) < 2 {
			w.printf("a command needs at least a service and a name, e.g. %q\n", r.Command.Service()+" stats")
			continue
		}
		invalid := ""
		for _, segment := range path {
			if !ValidPathSegment(segment) {
				invalid = segment
				break
			}
		}
		if invalid != "" {
			w.printf("%q is not usable as a command name: use lowercase letters, digits and dashes\n", invalid)
			continue
		}
		if _, exact := w.Tree.Find(path); exact {
			w.printf("%q already exists; pick another path\n", strings.Join(path, " "))
			continue
		}
		r.Command.Path = path

		parent, missing := w.Tree.Resolve(path)
		if len(missing) > 1 {
			w.printf("will create %d new command group(s) under %q: %s\n",
				len(missing)-1, parent.String(), strings.Join(missing[:len(missing)-1], " > "))
		}
		break
	}

	varName, err := w.askValid("command variable", CommandVar(r.Command.Path), ValidGoIdent,
		"a command variable has to be a Go identifier")
	if err != nil {
		return err
	}
	r.Command.Var = varName

	def := r.Command.Use() + ".go"
	if r.Command.File != "" {
		def = r.Command.File
	}
	file, err := w.askValid("file under internal/cli/"+r.Command.Service(), def, ValidFileName,
		"give a bare file name ending in .go")
	if err != nil {
		return err
	}
	r.Command.File = file

	if r.Command.Short, err = w.ask("short description", r.Command.Short); err != nil {
		return err
	}
	if r.Command.Long, err = w.ask("long description (optional)", r.Command.Long); err != nil {
		return err
	}
	return nil
}

func siblingPaths(tree *Tree, service string) []string {
	node, ok := tree.Find([]string{service})
	if !ok {
		return []string{"(new service package)"}
	}
	var out []string
	for _, c := range node.Children {
		out = append(out, c.Use)
	}
	sort.Strings(out)
	if len(out) == 0 {
		return []string{"(none yet)"}
	}
	return out
}

// arguments walks each path, query and header parameter.
func (w *Wizard) arguments(op *Operation, r *Recipe) error {
	if len(r.Args) == 0 {
		return nil
	}
	w.printf("\n--- Arguments ---\n")

	for i := range r.Args {
		arg := &r.Args[i]
		required := arg.Source == InPath || !w.optional(op, *arg)

		w.printf("\n%s parameter %q (%s)", arg.Source, arg.Name, arg.Type)
		if arg.Description != "" {
			w.printf(" — %s", arg.Description)
		}
		w.printf("\n")

		if !required {
			keep, err := w.confirm("expose it", true)
			if err != nil {
				return err
			}
			if !keep {
				arg.Skip = true
				continue
			}
		}

		cli, err := w.askValid("  flag name", arg.CLI, ValidFlagName,
			"a flag name starts with a letter and holds only letters, digits, dashes and underscores")
		if err != nil {
			return err
		}
		arg.CLI = cli
		if arg.Var == "" || arg.Var == LocalVar(arg.CLI) {
			arg.Var = LocalVar(cli)
		}

		if arg.Shorthand, err = w.askOptional("  shorthand (optional)", arg.Shorthand, singleChar,
			"a shorthand is a single character"); err != nil {
			return err
		}
		if arg.Description, err = w.ask("  flag help", arg.Description); err != nil {
			return err
		}
		if err := w.resolution(&arg.Resolve, &arg.Context, &arg.Persist, &arg.Prompt, arg.CLI); err != nil {
			return err
		}
		if !required {
			if err := w.unsetBehaviour(&arg.Unset, &arg.Default, arg.CLI, arg.Type); err != nil {
				return err
			}
		} else if len(arg.Resolve) == 1 && arg.Resolve[0] == FromFlag {
			if err := w.defaultValue(&arg.Default, arg.CLI, arg.Type); err != nil {
				return err
			}
		}
	}
	return nil
}

// authorization asks whether the command has to send credentials, for the
// endpoints whose spec never declared an authorization parameter.
//
// The generated client has no field for such a header, so the only way to send
// it is a RequestEditorFn. Derive already adds one when a vendored spec says the
// operation is secured; this covers the case where there is no spec to consult.
func (w *Wizard) authorization(op *Operation, r *Recipe) error {
	if _, declared := op.HeaderParam("authorization"); declared {
		return nil
	}
	for _, a := range r.Args {
		if a.Source == InHeader && strings.EqualFold(a.Name, "authorization") {
			return nil // Derive added it from the spec
		}
	}

	w.printf("\n--- Authorization ---\n")
	w.printf("%s declares no authorization parameter, so nothing would be sent.\n", op.Name)
	if info := w.Specs.For(op); info.Known {
		w.printf("its spec declares no security requirement either.\n")
	} else {
		w.printf("no vendored spec to check; run `make generate` to fetch one.\n")
	}

	send, err := w.confirm("send an authorization header anyway (through a RequestEditorFn)", true)
	if err != nil {
		return err
	}
	if !send {
		return nil
	}

	arg := AuthorizationArg()
	if arg.CLI, err = w.ask("  flag name", arg.CLI); err != nil {
		return err
	}
	arg.Var = LocalVar(arg.CLI)
	if arg.Var == "token" {
		arg.Var = "jwt"
	}
	if err := w.resolution(&arg.Resolve, &arg.Context, &arg.Persist, &arg.Prompt, arg.CLI); err != nil {
		return err
	}
	r.Args = append(r.Args, arg)
	return nil
}

func (w *Wizard) optional(op *Operation, arg RecipeArg) bool {
	for _, p := range op.Params {
		if p.In == arg.Source && strings.EqualFold(p.Name, arg.Name) {
			return p.Optional
		}
	}
	return false
}

// resolution asks where a value comes from: the flag, the context table, or a
// prompt — the chain internal/command/args.go implements.
func (w *Wizard) resolution(resolve *[]string, contextField *string, persist *bool, prompt *string, cli string) error {
	options := []string{
		"flag, then context, then prompt",
		"flag, then prompt",
		"flag only",
	}
	def := 1
	switch {
	case hasSource(*resolve, FromContext):
		def = 0
	case hasSource(*resolve, FromPrompt):
		def = 1
	default:
		def = 2
	}

	choice, err := w.choose("  where does the value come from", options, def)
	if err != nil {
		return err
	}
	switch choice {
	case 0:
		*resolve = []string{FromFlag, FromContext, FromPrompt}
	case 1:
		*resolve = []string{FromFlag, FromPrompt}
		*contextField = ""
		*persist = false
	case 2:
		*resolve = []string{FromFlag}
		*contextField = ""
		*persist = false
		*prompt = ""
		return nil
	}

	if choice == 0 {
		field, err := w.askValid("  context field ("+strings.Join(w.Schema.Names(), ", ")+")",
			firstNonEmpty(*contextField, GoFieldName(ColumnName(cli))), ValidGoIdent,
			"a context field is a Go field name on db.Context, e.g. OrgName")
		if err != nil {
			return err
		}
		*contextField = field
		if _, known := w.Schema.Field(field); !known {
			w.printf("  %s is new; a migration will add the %s column\n", field, ColumnName(field))
		}
		save, err := w.confirm("  save the resolved value back to context", true)
		if err != nil {
			return err
		}
		*persist = save
	}

	p, err := w.ask("  prompt text", firstNonEmpty(*prompt, cli))
	if err != nil {
		return err
	}
	*prompt = p
	return nil
}

// unsetBehaviour asks what an optional value does when its flag is not given:
// leave the field out and let the API apply its own default, or always send it
// carrying a default the command states itself.
func (w *Wizard) unsetBehaviour(unset, def *string, cli, goType string) error {
	kind := KindFor(goType, "", false, "")
	options := []string{
		"omit it — leave the field out of the request, so the API applies its own default",
		"always send it, carrying a default value",
	}
	choice := 0
	if *unset == UnsetDefault {
		choice = 1
	}

	choice, err := w.choose("  when --"+cli+" is not given", options, choice)
	if err != nil {
		return err
	}
	if choice == 0 {
		*unset = UnsetOmit
		*def = ""
		return nil
	}
	*unset = UnsetDefault

	zero := strings.Trim(kind.Default, `"`)
	if zero == "nil" {
		zero = ""
	}
	value, err := w.askValid("  default value", firstNonEmpty(*def, zero), func(v string) bool {
		if !kind.Supported {
			return true
		}
		_, err := kind.DefaultLiteral(v)
		return err == nil
	}, "that is not a valid "+goType)
	if err != nil {
		return err
	}
	*def = value
	return nil
}

// defaultValue asks for a flag's default, for a value resolved from its flag
// alone. Without one the command silently sends the type's zero.
func (w *Wizard) defaultValue(def *string, cli, goType string) error {
	kind := KindFor(goType, "", false, "")
	if !kind.Supported {
		return nil
	}

	value, err := w.askValid("  default value when --"+cli+" is not given", *def, func(v string) bool {
		_, err := kind.DefaultLiteral(v)
		return err == nil
	}, "that is not a valid "+goType)
	if err != nil {
		return err
	}
	*def = value
	return nil
}

// body asks how the request body is built.
func (w *Wizard) body(op *Operation, r *Recipe) error {
	if op.Body == nil {
		return nil
	}
	w.printf("\n--- Request body (%s) ---\n", op.Body.DTO)

	switch {
	case op.Body.Freeform:
		w.printf("the spec patcher genericized this body, so it is passed as raw JSON through --%s\n", firstNonEmpty(r.Body.Flag, "body"))
		return nil
	case op.Body.Opaque:
		w.printf("this body reaches a type the client declares with an unexported union field,\n")
		w.printf("so it cannot be built field by field and is passed as raw JSON through --%s\n", firstNonEmpty(r.Body.Flag, "body"))
		return nil
	case r.Body.Mode == BodyRaw:
		w.printf("a required field of this body cannot be carried on a flag, so the whole body\n")
		w.printf("is passed as raw JSON through --%s\n", firstNonEmpty(r.Body.Flag, "body"))
		return nil
	}

	raw, err := w.confirm("take the whole body as raw JSON instead of a flag per field", false)
	if err != nil {
		return err
	}
	if raw {
		r.Body.Mode = BodyRaw
		r.Body.Flag = "body"
		r.Body.Fields = nil
		return nil
	}

	for i := range r.Body.Fields {
		f := &r.Body.Fields[i]
		if f.Skip {
			w.printf("\nfield %q (%s) cannot be carried on a flag and is left out\n", f.Name, f.Type)
			continue
		}

		w.printf("\nfield %q (%s)", f.Name, f.Type)
		if f.Optional {
			w.printf(" optional")
		}
		if f.Description != "" {
			w.printf(" — %s", f.Description)
		}
		w.printf("\n")

		if f.Optional {
			keep, err := w.confirm("expose it", true)
			if err != nil {
				return err
			}
			if !keep {
				f.Skip = true
				continue
			}
		}

		cli, err := w.askValid("  flag name", f.CLI, ValidFlagName,
			"a flag name starts with a letter and holds only letters, digits, dashes and underscores")
		if err != nil {
			return err
		}
		f.CLI = cli
		f.Var = LocalVar(cli)
		if f.Shorthand, err = w.askOptional("  shorthand (optional)", f.Shorthand, singleChar,
			"a shorthand is a single character"); err != nil {
			return err
		}
		if f.Description, err = w.ask("  flag help", f.Description); err != nil {
			return err
		}
		if err := w.resolution(&f.Resolve, &f.Context, &f.Persist, &f.Prompt, f.CLI); err != nil {
			return err
		}
		if f.Optional {
			if err := w.unsetBehaviour(&f.Unset, &f.Default, f.CLI, f.Type); err != nil {
				return err
			}
		} else if len(f.Resolve) == 1 && f.Resolve[0] == FromFlag {
			if err := w.defaultValue(&f.Default, f.CLI, f.Type); err != nil {
				return err
			}
		}
	}
	return nil
}

// persistence asks which parts of a successful response to keep in the context
// table, which is how one command's output feeds the next command's arguments.
func (w *Wizard) persistence(op *Operation, r *Recipe) error {
	w.printf("\n--- Response ---\n")

	resp, hasTyped := op.Success()
	if hasTyped {
		w.printf("%s returns %s on %d\n", op.Name, resp.GoType, resp.Status)
	} else {
		w.printf("%s documents no typed success body\n", op.Name)
	}

	options := []string{
		"pretty-print the response body whatever the status",
		"pretty-print the typed success field, falling back to the body",
	}
	def := 0
	if r.Response.Render == RenderTyped {
		def = 1
	}
	if !hasTyped {
		w.printf("rendering the raw body, since there is no typed field\n")
		r.Response.Render = RenderBody
	} else {
		choice, err := w.choose("how should the result be printed", options, def)
		if err != nil {
			return err
		}
		if choice == 1 {
			r.Response.Render = RenderTyped
			r.Response.Field = resp.Field
		} else {
			r.Response.Render = RenderBody
		}
	}

	// Response headers first: the JWT and the agent hello token both arrive
	// that way rather than in the body.
	for {
		add, err := w.confirm("keep a response header in the context table", false)
		if err != nil {
			return err
		}
		if !add {
			break
		}
		name, err := w.askValid("  header name", "authorization", ValidFlagName,
			"a header name starts with a letter and holds letters, digits, dashes and underscores")
		if err != nil {
			return err
		}
		field, err := w.askValid("  context field", "Jwt", ValidGoIdent,
			"a context field is a Go field name on db.Context, e.g. Jwt")
		if err != nil {
			return err
		}
		r.Response.Persist = append(r.Response.Persist, RecipePersist{
			From: "header:" + name, Context: field, Column: ColumnName(field),
		})
		w.noteColumn(field)
	}

	if !hasTyped {
		return nil
	}
	if resp.List {
		w.printf("%s is an array, so it has no addressable fields to keep\n", resp.Field)
		return nil
	}
	if len(resp.Fields) == 0 {
		return nil
	}

	keep, err := w.confirm(fmt.Sprintf("keep any of %s's %d fields in the context table", resp.DTO, len(resp.Fields)), false)
	if err != nil {
		return err
	}
	if !keep {
		return nil
	}

	for {
		labels := []string{"(done)"}
		for _, f := range resp.Fields {
			labels = append(labels, fmt.Sprintf("%-24s %s  %s", f.GoName, f.GoType, f.Doc))
		}
		choice, err := w.choose("  which field", labels, 0)
		if err != nil {
			return err
		}
		if choice == 0 {
			return nil
		}
		field := resp.Fields[choice-1]

		target, err := w.askValid("  context field", GoFieldName(ColumnName(field.GoName)), ValidGoIdent,
			"a context field is a Go field name on db.Context, e.g. AgentUuid")
		if err != nil {
			return err
		}
		r.Response.Persist = append(r.Response.Persist, RecipePersist{
			From:    resp.Field + ":" + field.GoName,
			Context: target,
			Column:  ColumnName(target),
		})
		w.noteColumn(target)
	}
}

// singleChar reports whether s is exactly one character, the shape of a cobra
// flag shorthand.
func singleChar(s string) bool { return len([]rune(s)) == 1 }

func (w *Wizard) noteColumn(field string) {
	if _, known := w.Schema.Field(field); !known {
		w.printf("  %s is new; a migration will add the %s column\n", field, ColumnName(field))
	}
}
