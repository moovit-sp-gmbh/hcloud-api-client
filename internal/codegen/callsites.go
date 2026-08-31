package codegen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// CallSite is one place in the command tree that calls a client operation.
type CallSite struct {
	Client  string // client package name, e.g. "mothershipClient"
	Op      string // operation name, e.g. "GetV1OrgOrgNameStats"
	Pkg     string // package directory relative to the repo root
	File    string
	Line    int
	Command string // the cobra command variable the call sits in
}

// Key matches Operation.Key.
func (c CallSite) Key() string { return c.Client + "." + c.Op }

// LoadCallSites finds every client operation the command tree actually calls.
//
// This is not the source of truth for coverage — codegen/manifest.yaml is — but
// it is how the manifest gets seeded, and comparing the two catches a manifest
// that has drifted from the code.
func LoadCallSites(root string) ([]CallSite, error) {
	ctors, err := loadConstructors(root)
	if err != nil {
		return nil, err
	}

	cliRoot := filepath.Join(root, "internal", "cli")
	var files []string
	err = filepath.WalkDir(cliRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)

	var out []CallSite
	for _, path := range files {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		rel, _ := filepath.Rel(root, path)
		pkg := filepath.Dir(rel)

		// Which client packages does this file import? All current files import
		// at most one, which alone resolves the call. When a file imports more
		// than one, fall back to tracing the receiver back to its constructor.
		var imported []string
		aliases := map[string]string{}
		for _, imp := range f.Imports {
			p := strings.Trim(imp.Path.Value, `"`)
			if !strings.Contains(p, "/internal/clients/") {
				continue
			}
			name := filepath.Base(p)
			if imp.Name != nil {
				aliases[imp.Name.Name] = name
			}
			imported = append(imported, name)
		}

		// Map each command variable to the source range of its declaration, so
		// a call can be attributed to the command it belongs to.
		ranges := commandRanges(f)

		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || !strings.HasSuffix(sel.Sel.Name, "WithResponse") {
				return true
			}
			op := strings.TrimSuffix(sel.Sel.Name, "WithResponse")
			for _, suffix := range []string{"WithBody", "WithFormdataBody"} {
				op = strings.TrimSuffix(op, suffix)
			}

			client := ""
			switch {
			case len(imported) == 1:
				client = imported[0]
			default:
				client = resolveReceiver(f, sel.X, ctors, aliases)
			}
			if client == "" {
				return true
			}

			pos := fset.Position(call.Pos())
			out = append(out, CallSite{
				Client:  client,
				Op:      op,
				Pkg:     pkg,
				File:    rel,
				Line:    pos.Line,
				Command: commandAt(ranges, call.Pos()),
			})
			return true
		})
	}
	return out, nil
}

// loadConstructors reads internal/clients/clients.go and maps each
// New<Svc>Client function to the client package it returns, so the mapping
// cannot drift from the code.
func loadConstructors(root string) (map[string]string, error) {
	path := filepath.Join(root, "internal", "clients", "clients.go")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}

	out := map[string]string{}
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil || fn.Type.Results == nil || len(fn.Type.Results.List) == 0 {
			continue
		}
		star, ok := fn.Type.Results.List[0].Type.(*ast.StarExpr)
		if !ok {
			continue
		}
		sel, ok := star.X.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "ClientWithResponses" {
			continue
		}
		if id, ok := sel.X.(*ast.Ident); ok {
			out[fn.Name.Name] = id.Name
		}
	}
	return out, nil
}

// resolveReceiver traces a client variable back to the constructor it came from.
func resolveReceiver(f *ast.File, recv ast.Expr, ctors, aliases map[string]string) string {
	id, ok := recv.(*ast.Ident)
	if !ok {
		return ""
	}

	found := ""
	ast.Inspect(f, func(n ast.Node) bool {
		as, ok := n.(*ast.AssignStmt)
		if !ok || found != "" {
			return found == ""
		}
		assigned := false
		for _, lhs := range as.Lhs {
			if lid, ok := lhs.(*ast.Ident); ok && lid.Name == id.Name {
				assigned = true
			}
		}
		if !assigned || len(as.Rhs) != 1 {
			return true
		}
		call, ok := as.Rhs[0].(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		switch {
		case sel.Sel.Name == "NewClientWithResponses":
			// idpClient.NewClientWithResponses(server)
			if pid, ok := sel.X.(*ast.Ident); ok {
				name := pid.Name
				if real, ok := aliases[name]; ok {
					name = real
				}
				found = name
			}
		default:
			// clients.NewMothershipClient(cmd)
			if pkg, ok := ctors[sel.Sel.Name]; ok {
				found = pkg
			}
		}
		return found == ""
	})
	return found
}

// commandRange is the source span of one command variable's declaration.
type commandRange struct {
	name       token.Pos
	start, end token.Pos
	varName    string
}

func commandRanges(f *ast.File) []commandRange {
	var out []commandRange
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok || len(vs.Names) != 1 || len(vs.Values) != 1 {
				continue
			}
			if _, _, ok := cobraCommandLit(vs.Values[0]); !ok {
				continue
			}
			out = append(out, commandRange{
				start:   vs.Values[0].Pos(),
				end:     vs.Values[0].End(),
				varName: vs.Names[0].Name,
			})
		}
	}
	return out
}

func commandAt(ranges []commandRange, pos token.Pos) string {
	for _, r := range ranges {
		if pos >= r.start && pos <= r.end {
			return r.varName
		}
	}
	return ""
}
