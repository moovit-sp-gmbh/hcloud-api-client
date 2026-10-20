package codegen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ApplyOptions controls how a recipe is applied.
type ApplyOptions struct {
	DryRun bool
	// Now stamps generated migration filenames; zero means time.Now.
	Now time.Time
	// SkipSQLC leaves `sqlc generate` to the caller. Only useful in tests: a
	// new context column does not exist on db.Context until sqlc has run, so
	// the tree will not build without it.
	SkipSQLC bool
	// Index, when set, is reused instead of reparsing the generated clients.
	// Those are 4.4MB of Go, so a bulk run that reloaded them per command would
	// spend almost all of its time there.
	Index *Index
}

// ApplyResult records what applying a recipe did.
type ApplyResult struct {
	Created  []string
	Edited   []string
	Warnings []string
	Recipe   string
	Columns  []ContextField
}

// Apply generates a command from a recipe and wires it into the tree.
func Apply(root string, r *Recipe, opts ApplyOptions) (*ApplyResult, error) {
	if opts.Now.IsZero() {
		opts.Now = time.Now()
	}

	if err := r.Validate(); err != nil {
		return nil, err
	}

	ix := opts.Index
	if ix == nil {
		var err error
		if ix, err = LoadIndex(root); err != nil {
			return nil, err
		}
	}
	op, ok := ix.Op(r.Client, r.Operation)
	if !ok {
		return nil, fmt.Errorf("%s has no operation %s", r.Client, r.Operation)
	}
	if !op.Callable() {
		return nil, fmt.Errorf("%s %s cannot be called: %s", op.Method, op.Path, op.Blocked)
	}
	// The tree is reloaded every time: a previous command may have added the
	// group this one hangs under.
	tree, err := LoadTree(root)
	if err != nil {
		return nil, err
	}
	schema, err := LoadContextSchema(root)
	if err != nil {
		return nil, err
	}
	ctors, err := loadConstructors(root)
	if err != nil {
		return nil, err
	}

	a := &apply{root: root, recipe: r, op: op, tree: tree, schema: schema, ctors: ctors, opts: opts, res: &ApplyResult{}}

	// Context field names are canonicalized before anything else reads them:
	// the migration, sqlc's model and the generated code all have to agree on
	// one spelling, and sqlc's is the one that wins.
	a.canonicalizeContext()

	// New context columns come first: the generated command references fields
	// that only exist on db.Context once sqlc has regenerated the model.
	if err := a.contextColumns(); err != nil {
		return nil, err
	}
	if err := a.servicePackage(); err != nil {
		return nil, err
	}
	parentVar, parentFile, err := a.groups()
	if err != nil {
		return nil, err
	}
	if err := a.clientConstructor(); err != nil {
		return nil, err
	}
	if err := a.command(parentVar, parentFile); err != nil {
		return nil, err
	}
	if err := a.manifest(); err != nil {
		return nil, err
	}
	return a.res, nil
}

type apply struct {
	root   string
	recipe *Recipe
	op     *Operation
	tree   *Tree
	schema *ContextSchema
	ctors  map[string]string
	opts   ApplyOptions
	res    *ApplyResult
	// pending records files this run has already written, so a dry run — where
	// nothing lands on disk — still knows a file it just created exists.
	pending map[string]bool
}

// canonicalizeContext rewrites every context field name in the recipe to the
// spelling sqlc generates, so "agent_id", "agentId" and "AgentId" all resolve
// to db.Context.AgentID rather than to a field that does not exist.
func (a *apply) canonicalizeContext() {
	canonical := func(name string) string {
		if name == "" {
			return ""
		}
		if f, ok := a.schema.Field(name); ok {
			return f.GoName
		}
		return GoFieldName(ColumnName(name))
	}

	for i := range a.recipe.Args {
		a.recipe.Args[i].Context = canonical(a.recipe.Args[i].Context)
	}
	for i := range a.recipe.Body.Fields {
		a.recipe.Body.Fields[i].Context = canonical(a.recipe.Body.Fields[i].Context)
	}
	for i := range a.recipe.Response.Persist {
		p := &a.recipe.Response.Persist[i]
		p.Context = canonical(p.Context)
		if f, ok := a.schema.Field(p.Context); ok {
			p.Column = f.Column
		} else if p.Column == "" {
			p.Column = ColumnName(p.Context)
		}
	}
}

func (a *apply) created(path string) {
	for _, p := range a.res.Created {
		if p == path {
			return
		}
	}
	a.res.Created = append(a.res.Created, path)
}

func (a *apply) edited(path string) {
	// A file created by this run is reported as created, not edited.
	for _, p := range a.res.Created {
		if p == path {
			return
		}
	}
	for _, p := range a.res.Edited {
		if p == path {
			return
		}
	}
	a.res.Edited = append(a.res.Edited, path)
}

func (a *apply) warn(format string, args ...any) {
	a.res.Warnings = append(a.res.Warnings, fmt.Sprintf(format, args...))
}

// write writes a file unless this is a dry run.
func (a *apply) write(rel string, src []byte, isNew bool) error {
	if isNew {
		a.created(rel)
	} else {
		a.edited(rel)
	}
	a.note(rel)
	if a.opts.DryRun {
		return nil
	}
	return WriteFormatted(filepath.Join(a.root, rel), src)
}

// onlyPending reports that a file exists only because this run claimed it and
// nothing was actually written, which happens during a dry run.
func (a *apply) onlyPending(rel string) bool {
	if !a.pending[rel] {
		return false
	}
	_, err := os.Stat(filepath.Join(a.root, rel))
	return err != nil
}

// note records that this run has produced a file.
func (a *apply) note(rel string) {
	if a.pending == nil {
		a.pending = map[string]bool{}
	}
	a.pending[rel] = true
}

// writeRaw writes a non-Go file.
func (a *apply) writeRaw(rel string, data []byte, isNew bool) error {
	if isNew {
		a.created(rel)
	} else {
		a.edited(rel)
	}
	a.note(rel)
	if a.opts.DryRun {
		return nil
	}
	path := filepath.Join(a.root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (a *apply) read(rel string) ([]byte, error) {
	return os.ReadFile(filepath.Join(a.root, rel))
}

func (a *apply) exists(rel string) bool {
	if a.pending[rel] {
		return true
	}
	_, err := os.Stat(filepath.Join(a.root, rel))
	return err == nil
}

// contextColumns adds any context column the recipe persists into that the
// schema does not have yet: a goose migration, the CreateContext insert, sqlc,
// and the field ContextService.Push copies across.
func (a *apply) contextColumns() error {
	wanted := map[string]string{} // Go field name -> column
	add := func(goName, column string) {
		if goName == "" {
			return
		}
		if _, ok := a.schema.Field(goName); ok {
			return
		}
		if column == "" {
			column = ColumnName(goName)
		}
		wanted[GoFieldName(column)] = column
	}

	for _, p := range a.recipe.Response.Persist {
		add(p.Context, p.Column)
	}
	for _, arg := range a.recipe.Args {
		if !arg.Skip && arg.Persist {
			add(arg.Context, "")
		}
	}
	for _, f := range a.recipe.Body.Fields {
		if !f.Skip && f.Persist {
			add(f.Context, "")
		}
	}
	if len(wanted) == 0 {
		return nil
	}

	// Deterministic order so a migration's contents do not depend on map
	// iteration.
	var goNames []string
	for goName := range wanted {
		goNames = append(goNames, goName)
	}
	sort.Strings(goNames)

	var columns []string
	for _, goName := range goNames {
		columns = append(columns, wanted[goName])
		a.res.Columns = append(a.res.Columns, ContextField{GoName: goName, Column: wanted[goName]})
	}

	name := columns[0]
	if len(columns) > 1 {
		name = a.recipe.Command.Use()
	}
	rel := filepath.Join("internal", "db", "migrations", MigrationFilename(name, a.opts.Now))
	if err := a.writeRaw(rel, []byte(MigrationContent(columns)), true); err != nil {
		return err
	}

	if !a.opts.DryRun {
		if err := AddCreateContextColumns(a.root, columns); err != nil {
			return err
		}
	}
	a.edited(filepath.Join("internal", "db", "context.sql"))

	if a.opts.DryRun || a.opts.SkipSQLC {
		a.edited(filepath.Join("internal", "db", "models.go"))
		a.edited(filepath.Join("internal", "db", "context.sql.go"))
	} else if err := a.runSQLC(); err != nil {
		return err
	}

	// Push enumerates every field by hand, so a new column is silently dropped
	// until it is listed there.
	src, err := a.read("internal/service.go")
	if err != nil {
		return err
	}
	for _, goName := range goNames {
		src, err = AddPushField(src, goName)
		if err != nil {
			return fmt.Errorf("internal/service.go: %w", err)
		}
	}
	return a.write("internal/service.go", src, false)
}

func (a *apply) runSQLC() error {
	cmd := exec.Command("sqlc", "generate")
	cmd.Dir = a.root
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("sqlc generate failed (the new column will not exist on db.Context until it runs): %w\n%s", err, out)
	}
	a.edited(filepath.Join("internal", "db", "models.go"))
	a.edited(filepath.Join("internal", "db", "context.sql.go"))
	return nil
}

// servicePackage creates internal/cli/<service> and registers it on the root
// command when the service has no commands yet.
func (a *apply) servicePackage() error {
	service := a.recipe.Command.Service()
	if _, ok := a.tree.Find([]string{service}); ok {
		return nil
	}

	rel := filepath.Join("internal", "cli", service, service+".go")
	if !a.exists(rel) {
		group, err := ServiceGroupFile(service, "", "")
		if err != nil {
			return err
		}
		if err := a.write(rel, group, true); err != nil {
			return err
		}
	}

	src, err := a.read("cmd/root.go")
	if err != nil {
		return err
	}
	src, err = RegisterService(src, service)
	if err != nil {
		return err
	}
	return a.write("cmd/root.go", src, false)
}

// groups creates any intermediate command group the path needs and returns the
// variable and file the new command should be registered against.
func (a *apply) groups() (parentVar, parentFile string, err error) {
	path := a.recipe.Command.Path
	service := a.recipe.Command.Service()
	pkgDir := filepath.Join("internal", "cli", service)

	// The service group itself, whether it existed or was just created.
	parentVar = service + "Command"
	parentFile = filepath.Join(pkgDir, service+".go")
	if n, ok := a.tree.Find([]string{service}); ok {
		parentVar, parentFile = n.Var, n.File
	}

	// Everything between the service and the leaf. New groups go in the same
	// file as the command being generated: internal/cli/mothership/agent.go
	// holds both agentCommand and its subcommand, and splitting a group away
	// from its children would read worse than the convention it copies.
	targetFile := filepath.Join(pkgDir, a.commandFile())
	for i := 2; i < len(path); i++ {
		segment := path[i-1]
		if n, ok := a.tree.Find(path[:i]); ok {
			parentVar, parentFile = n.Var, n.File
			continue
		}

		groupVar := CommandVar(path[:i])
		decl := GroupDecl(groupVar, segment, fmt.Sprintf("Commands for %s", segment))

		if a.exists(targetFile) && !a.onlyPending(targetFile) {
			src, err := a.read(targetFile)
			if err != nil {
				return "", "", err
			}
			src, err = AddImports(AppendDecl(src, decl), []string{"github.com/spf13/cobra"})
			if err != nil {
				return "", "", err
			}
			if err := a.write(targetFile, src, false); err != nil {
				return "", "", err
			}
		} else {
			src := fmt.Sprintf("package %s\n\nimport \"github.com/spf13/cobra\"\n\n%s", service, decl)
			if err := a.write(targetFile, []byte(src), true); err != nil {
				return "", "", err
			}
		}

		if err := a.registerChild(parentFile, parentVar, groupVar); err != nil {
			return "", "", err
		}
		parentVar, parentFile = groupVar, targetFile
	}
	return parentVar, parentFile, nil
}

// commandFile is the file the command's declaration goes in.
func (a *apply) commandFile() string {
	if a.recipe.Command.File != "" {
		return a.recipe.Command.File
	}
	return a.recipe.Command.Use() + ".go"
}

// registerChild adds child.AddCommand into the file the parent is declared in.
func (a *apply) registerChild(parentFile, parentVar, childVar string) error {
	if !a.exists(parentFile) || a.onlyPending(parentFile) {
		// The parent file is itself part of this run and has not landed on disk,
		// which only happens during a dry run.
		a.edited(parentFile)
		return nil
	}
	src, err := a.read(parentFile)
	if err != nil {
		return err
	}
	out, err := AddSubcommand(src, parentVar, childVar)
	if err != nil {
		return fmt.Errorf("%s: %w", parentFile, err)
	}
	return a.write(parentFile, out, false)
}

// clientConstructor adds the clients.New*Client function when the service does
// not have one — cosmoClient is generated with its constructor commented out.
func (a *apply) clientConstructor() error {
	ctor := constructorName(a.op.Client, a.ctors)
	if a.ctors[ctor] == a.op.Client {
		return nil
	}

	src, err := a.read("internal/clients/clients.go")
	if err != nil {
		return err
	}
	out, err := AddClientConstructor(src, a.op.Client, ctor)
	if err != nil {
		return err
	}
	if len(out) == len(src) {
		return nil
	}
	return a.write("internal/clients/clients.go", out, false)
}

// command emits the command itself and registers it under its parent.
func (a *apply) command(parentVar, parentFile string) error {
	g, err := Emit(a.recipe, a.op, a.schema, a.ctors)
	if err != nil {
		return err
	}
	a.res.Warnings = append(a.res.Warnings, g.Warnings...)

	// Register first. When a new group and its command share a file they also
	// share an init(), and inserting the AddCommand call after the flags on the
	// first run but before them on a regeneration would leave the file
	// reordering itself forever.
	if err := a.registerChild(parentFile, parentVar, a.recipe.Command.Var); err != nil {
		return err
	}

	if !a.exists(g.Path) || a.onlyPending(g.Path) {
		return a.write(g.Path, g.Source, !a.exists(g.Path))
	}

	src, err := a.read(g.Path)
	if err != nil {
		return err
	}
	declared, err := DeclaresVar(src, a.recipe.Command.Var)
	if err != nil {
		return fmt.Errorf("%s: %w", g.Path, err)
	}
	// Replacing a declaration is how regeneration works, but only when it is
	// this command's own. Two commands deriving the same variable name would
	// otherwise silently overwrite each other.
	if declared {
		if n, ok := a.tree.Node(g.Dir, a.recipe.Command.Var); ok {
			if want := strings.Join(a.recipe.Command.Path, " "); n.String() != want {
				return fmt.Errorf("%s already declares %s for %q; %q needs a different command variable",
					g.Path, a.recipe.Command.Var, n.String(), want)
			}
		}
	}

	if src, err = AddImports(src, g.Imports); err != nil {
		return err
	}
	// SetFlagsFor replaces this command's registrations where they stand, so a
	// file shared by several commands keeps its order across regenerations.
	if src, err = SetFlagsFor(src, a.recipe.Command.Var, g.Flags); err != nil {
		return err
	}
	if declared {
		// Regenerating: the recipe owns this declaration, so replace it rather
		// than appending a second copy.
		if src, err = ReplaceVarDecl(src, a.recipe.Command.Var, g.VarBlock); err != nil {
			return err
		}
	} else {
		src = AppendDecl(src, g.VarBlock)
	}
	return a.write(g.Path, src, false)
}

// manifest records the endpoint and saves the recipe for replay.
func (a *apply) manifest() error {
	rel := a.recipe.RecipePath()
	a.res.Recipe = rel

	if !a.opts.DryRun {
		if _, err := a.recipe.Save(a.root); err != nil {
			return err
		}
	}
	if a.exists(rel) {
		a.edited(rel)
	} else {
		a.created(rel)
	}

	m, err := LoadManifest(a.root)
	if err != nil {
		return err
	}
	m.Upsert(Entry{
		Client:    a.op.Client,
		Operation: a.op.Name,
		Method:    a.op.Method,
		Path:      a.op.Path,
		Command:   strings.Join(a.recipe.Command.Path, " "),
		Status:    StatusImplemented,
		Source:    SourceGenerated,
		Recipe:    rel,
	})
	a.edited(ManifestPath)
	if a.opts.DryRun {
		return nil
	}
	return m.Save(a.root)
}

// GroupDecl renders a command group: a parent that prints its own help, the way
// internal/cli/mothership/agent.go does.
func GroupDecl(varName, use, short string) string {
	return fmt.Sprintf(`var %s = &cobra.Command{
	Use:   %q,
	Short: %q,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}`, varName, use, short)
}
