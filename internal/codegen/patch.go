package codegen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// The edits in this file locate an insertion point with go/ast and then splice
// text, rather than reprinting the file from its AST. Reprinting would reformat
// and drop comments across files the generator only means to add a line to.

// AddImports adds import paths to a file, leaving everything else alone. Only
// the import block is rewritten, and only when something is actually missing.
func AddImports(src []byte, paths []string) ([]byte, error) {
	if len(paths) == 0 {
		return src, nil
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	have := map[string]bool{}
	for _, imp := range f.Imports {
		have[strings.Trim(imp.Path.Value, `"`)] = true
	}
	var missing []string
	for _, p := range paths {
		if !have[p] {
			missing = append(missing, p)
			have[p] = true
		}
	}
	if len(missing) == 0 {
		return src, nil
	}

	all := make([]string, 0, len(have))
	for p := range have {
		all = append(all, p)
	}
	sort.Strings(all)

	var block strings.Builder
	block.WriteString("import (\n")
	for i, group := range groupImports(all) {
		if i > 0 {
			block.WriteString("\n")
		}
		for _, imp := range group {
			fmt.Fprintf(&block, "\t%q\n", imp)
		}
	}
	block.WriteString(")")

	// Replace the existing import declaration, or insert one after the package
	// clause when the file has none.
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.IMPORT {
			continue
		}
		start := fset.Position(gd.Pos()).Offset
		end := fset.Position(gd.End()).Offset
		return splice(src, start, end, block.String()), nil
	}

	pkgEnd := fset.Position(f.Name.End()).Offset
	return splice(src, pkgEnd, pkgEnd, "\n\n"+block.String()), nil
}

// AddFlagsToInit adds flag registrations to a file's init() function, creating
// one after the imports if the file has none.
func AddFlagsToInit(src []byte, lines []string) ([]byte, error) {
	if len(lines) == 0 {
		return src, nil
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	body := "\n"
	for _, l := range lines {
		body += "\t" + l + "\n"
	}

	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil || fn.Name.Name != "init" || fn.Body == nil {
			continue
		}
		// Insert right after the opening brace rather than at the end. When a
		// group and its subcommand share a file they share an init(), and
		// regeneration removes the flag lines and adds them back — appending
		// would put them on the far side of the AddCommand call each time, so
		// the file would never stop changing.
		opening := fset.Position(fn.Body.Lbrace).Offset + 1
		insert := "\n"
		for _, l := range lines {
			insert += "\t" + l + "\n"
		}
		return splice(src, opening, opening, strings.TrimSuffix(insert, "\n")), nil
	}

	// No init() yet. Put one after the import block, which is where every file
	// in internal/cli keeps it.
	anchor := fset.Position(f.Name.End()).Offset
	for _, decl := range f.Decls {
		if gd, ok := decl.(*ast.GenDecl); ok && gd.Tok == token.IMPORT {
			anchor = fset.Position(gd.End()).Offset
		}
	}
	return splice(src, anchor, anchor, "\n\nfunc init() {"+body+"}"), nil
}

// AppendDecl appends a declaration to the end of a file.
func AppendDecl(src []byte, decl string) []byte {
	out := bytes.TrimRight(src, "\n")
	return append(out, []byte("\n\n"+strings.TrimRight(decl, "\n")+"\n")...)
}

// DeclaresVar reports whether a file declares a package-level variable.
func DeclaresVar(src []byte, name string) (bool, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)
	if err != nil {
		return false, err
	}
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, n := range vs.Names {
				if n.Name == name {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

// ReplaceVarDecl swaps the declaration of a package-level variable for a new
// one. This is what makes `codegen regen` idempotent: the recipe owns that
// declaration, so replaying it rewrites the block instead of appending a second
// one.
func ReplaceVarDecl(src []byte, name, decl string) ([]byte, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, d := range f.Decls {
		gd, ok := d.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok || len(vs.Names) != 1 || vs.Names[0].Name != name {
				continue
			}
			start := fset.Position(gd.Pos()).Offset
			// Take any doc comment with the declaration, so a replacement does
			// not leave the old comment stranded above the new block.
			if gd.Doc != nil {
				start = fset.Position(gd.Doc.Pos()).Offset
			}
			end := fset.Position(gd.End()).Offset
			return splice(src, start, end, strings.TrimRight(decl, "\n")), nil
		}
	}
	return nil, fmt.Errorf("no declaration of %s to replace", name)
}

// RemoveFlagsFor drops every flag registration in init() that targets one
// command variable, so regeneration re-registers rather than accumulating.
func RemoveFlagsFor(src []byte, varName string) ([]byte, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil || fn.Name.Name != "init" || fn.Body == nil {
			continue
		}

		// Collect the source lines to drop, then rebuild the file without them.
		drop := map[int]bool{}
		for _, stmt := range fn.Body.List {
			expr, ok := stmt.(*ast.ExprStmt)
			if !ok {
				continue
			}
			if !isFlagRegistration(expr.X, varName) {
				continue
			}
			from := fset.Position(stmt.Pos()).Line
			to := fset.Position(stmt.End()).Line
			for l := from; l <= to; l++ {
				drop[l] = true
			}
		}
		if len(drop) == 0 {
			return src, nil
		}

		var kept []string
		for i, line := range strings.Split(string(src), "\n") {
			if !drop[i+1] {
				kept = append(kept, line)
			}
		}
		return []byte(strings.Join(kept, "\n")), nil
	}
	return src, nil
}

// isFlagRegistration matches <varName>.Flags().<Type>(...) and its
// PersistentFlags equivalent.
func isFlagRegistration(e ast.Expr, varName string) bool {
	call, ok := e.(*ast.CallExpr)
	if !ok {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	inner, ok := sel.X.(*ast.CallExpr)
	if !ok {
		return false
	}
	innerSel, ok := inner.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if innerSel.Sel.Name != "Flags" && innerSel.Sel.Name != "PersistentFlags" {
		return false
	}
	id, ok := innerSel.X.(*ast.Ident)
	return ok && id.Name == varName
}

// AddSubcommand registers child under parent.
//
// Where the call goes depends on how the parent was declared: a package-level
// command collects its children in an init(), while the command a Command()
// factory returns collects them in the factory body, before the return.
func AddSubcommand(src []byte, parentVar, childVar string) ([]byte, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	call := fmt.Sprintf("%s.AddCommand(%s)", parentVar, childVar)
	if bytes.Contains(src, []byte(call)) {
		return src, nil
	}

	// Prefer to sit with the parent's existing AddCommand calls.
	var lastCall *ast.CallExpr
	var inFactory *ast.FuncDecl
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			c, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := c.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "AddCommand" {
				return true
			}
			recv, ok := sel.X.(*ast.Ident)
			if !ok || recv.Name != parentVar {
				return true
			}
			lastCall = c
			inFactory = fn
			return true
		})
	}
	if lastCall != nil {
		end := fset.Position(lastCall.End()).Offset
		// Step past the statement's own newline so the new call lands on its
		// own line at the same indentation.
		indent := lineIndent(src, fset.Position(lastCall.Pos()).Offset)
		_ = inFactory
		return splice(src, end, end, "\n"+indent+call), nil
	}

	// The parent has no subcommands yet. If it is a factory local, insert
	// before the factory's return; otherwise use init().
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil || fn.Name.Name != "Command" || fn.Body == nil {
			continue
		}
		if !declaresLocal(fn, parentVar) {
			continue
		}
		for _, stmt := range fn.Body.List {
			ret, ok := stmt.(*ast.ReturnStmt)
			if !ok {
				continue
			}
			start := fset.Position(ret.Pos()).Offset
			indent := lineIndent(src, start)
			return splice(src, start, start, call+"\n\n"+indent), nil
		}
	}

	return AddFlagsToInit(src, []string{call})
}

// declaresLocal reports whether a function assigns to a local of this name.
func declaresLocal(fn *ast.FuncDecl, name string) bool {
	found := false
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		as, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}
		for _, lhs := range as.Lhs {
			if id, ok := lhs.(*ast.Ident); ok && id.Name == name {
				found = true
			}
		}
		return !found
	})
	return found
}

// AddPushField adds a field to the struct literal ContextService.Push builds,
// so a newly added context column is actually written on the next push.
func AddPushField(src []byte, goName string) ([]byte, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "Push" || fn.Body == nil {
			continue
		}

		var lit *ast.CompositeLit
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			cl, ok := n.(*ast.CompositeLit)
			if !ok {
				return true
			}
			// The literal is db.CreateContextParams{...}; take the outermost.
			if sel, ok := cl.Type.(*ast.SelectorExpr); ok && strings.HasSuffix(sel.Sel.Name, "Params") && lit == nil {
				lit = cl
			}
			return true
		})
		if lit == nil {
			return nil, fmt.Errorf("no CreateContextParams literal in Push")
		}

		for _, elt := range lit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			if id, ok := kv.Key.(*ast.Ident); ok && id.Name == goName {
				return src, nil // already pushed
			}
		}

		if len(lit.Elts) == 0 {
			return nil, fmt.Errorf("empty CreateContextParams literal in Push")
		}
		last := lit.Elts[len(lit.Elts)-1]
		end := fset.Position(last.End()).Offset
		indent := lineIndent(src, fset.Position(last.Pos()).Offset)
		// The literal is comma-separated with a trailing comma after each entry,
		// so step over that comma before inserting.
		if end < len(src) && src[end] == ',' {
			end++
		}
		return splice(src, end, end, fmt.Sprintf("\n%s%s: arg.%s,", indent, goName, goName)), nil
	}
	return nil, fmt.Errorf("no Push method found")
}

// AddClientConstructor adds a New<Svc>Client function to
// internal/clients/clients.go. cosmoClient is generated but has no constructor,
// so the first cosmo command needs one.
func AddClientConstructor(src []byte, client, ctorName string) ([]byte, error) {
	declared, err := declaresFunc(src, ctorName)
	if err != nil {
		return nil, err
	}
	if declared {
		return src, nil
	}

	// clients.go carries a commented-out NewCosmoClient whose body calls the
	// mothership client. Text-matching would mistake it for the real thing, and
	// leaving it behind next to a working constructor would be misleading.
	src = removeCommentedFunc(src, ctorName)

	service := Service(client)
	serverFn := "Server" + titleCase(service)
	if strings.EqualFold(service, "idp") {
		serverFn = "ServerIDP"
	}

	decl := fmt.Sprintf(`func %s(cmd *cobra.Command) (*%s.ClientWithResponses, error) {
	server, err := command.%s(cmd)
	if err != nil {
		return nil, err
	}

	client, err := %s.NewClientWithResponses(server)
	if err != nil {
		return nil, err
	}

	return client, nil
}`, ctorName, client, serverFn, client)

	out, err := AddImports(src, []string{modulePath + "/internal/clients/" + client})
	if err != nil {
		return nil, err
	}
	return AppendDecl(out, decl), nil
}

// declaresFunc reports whether a file declares a top-level function, checking
// the AST rather than the text so a commented-out stub does not count.
func declaresFunc(src []byte, name string) (bool, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)
	if err != nil {
		return false, err
	}
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Recv == nil && fn.Name.Name == name {
			return true, nil
		}
	}
	return false, nil
}

// removeCommentedFunc drops a run of comment lines that commented out a
// function of this name, along with any blank line left behind.
func removeCommentedFunc(src []byte, name string) []byte {
	lines := strings.Split(string(src), "\n")
	start := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "//") {
			continue
		}
		if strings.Contains(trimmed, "func "+name+"(") {
			start = i
			break
		}
	}
	if start < 0 {
		return src
	}

	// Walk back over any preceding comment lines, and forward to the end of the
	// commented block.
	for start > 0 && strings.HasPrefix(strings.TrimSpace(lines[start-1]), "//") {
		start--
	}
	end := start
	for end < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[end]), "//") {
		end++
	}
	for end < len(lines) && strings.TrimSpace(lines[end]) == "" {
		end++
	}
	for start > 0 && strings.TrimSpace(lines[start-1]) == "" {
		start--
	}

	kept := append(append([]string{}, lines[:start]...), lines[end:]...)
	return []byte(strings.Join(kept, "\n"))
}

// titleCase capitalizes the first letter, which is all strings.Title did that
// this needs.
func titleCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// ServiceGroupFile renders the <service>.go file that declares a new service's
// command group, matching internal/cli/mothership/mothership.go. high5, cosmo
// and dali have no command package at all, so the first command for one of them
// needs this.
func ServiceGroupFile(service, short, long string) ([]byte, error) {
	if short == "" {
		short = fmt.Sprintf("Entrypoint to %s commands", titleCase(service))
	}
	if long == "" {
		long = fmt.Sprintf("Entrypoint to helmut.cloud %s's commands", titleCase(service))
	}

	src := fmt.Sprintf(`package %s

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	%sCommand := &cobra.Command{
		Use:   %q,
		Short: %q,
		Long:  %q,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	%sCommand.PersistentFlags().String("server", "https://eu-central-1.app.helmut.cloud", "Set the target server/cluster")

	return %sCommand
}
`, service, service, service, short, long, service, service)

	return format.Source([]byte(src))
}

// RegisterService adds a service package to cmd/root.go.
func RegisterService(src []byte, service string) ([]byte, error) {
	out, err := AddImports(src, []string{modulePath + "/internal/cli/" + service})
	if err != nil {
		return nil, err
	}
	return AddSubcommand(out, "rootCmd", service+".Command()")
}

// splice replaces src[start:end] with text.
func splice(src []byte, start, end int, text string) []byte {
	out := make([]byte, 0, len(src)+len(text))
	out = append(out, src[:start]...)
	out = append(out, text...)
	return append(out, src[end:]...)
}

// lineIndent returns the leading whitespace of the line containing offset.
func lineIndent(src []byte, offset int) string {
	start := bytes.LastIndexByte(src[:offset], '\n') + 1
	line := src[start:offset]
	return string(line[:len(line)-len(bytes.TrimLeft(line, " \t"))])
}

// WriteFormatted gofmts and writes a file, creating parent directories.
func WriteFormatted(path string, src []byte) error {
	formatted, err := format.Source(src)
	if err != nil {
		return fmt.Errorf("%s: generated code does not parse: %w\n%s", path, err, numbered(string(src)))
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, formatted, 0o644)
}
