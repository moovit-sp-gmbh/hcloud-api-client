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

// Node is one command in the cobra tree.
type Node struct {
	Pkg   string // package directory relative to the repo root
	Var   string // the variable the command is declared as, e.g. "orgCommand"
	Use   string // the command's Use field, first word only
	Short string
	File  string // file the declaration lives in, relative to the repo root

	// Factory marks the command a package's Command() function returns. Those
	// are declared as locals inside the factory rather than at package level, so
	// a subcommand is attached in the factory body, not in an init().
	Factory bool

	Parent   *Node
	Children []*Node
}

// ID identifies a node across packages.
func (n *Node) ID() string { return n.Pkg + "." + n.Var }

// Path is the command path a user types, without the root command.
func (n *Node) Path() []string {
	var out []string
	for cur := n; cur != nil && cur.Parent != nil; cur = cur.Parent {
		out = append([]string{cur.Use}, out...)
	}
	return out
}

// String renders the command path, e.g. "mothership connect search". The root
// itself has no path, so it names itself.
func (n *Node) String() string {
	if n.Parent == nil {
		return "the root command"
	}
	return strings.Join(n.Path(), " ")
}

// Tree is the CLI's command tree as it exists in the source.
type Tree struct {
	Root  *Node
	nodes map[string]*Node
}

// Walk calls fn for every node below the root, parents before children.
func (t *Tree) Walk(fn func(*Node)) {
	var walk func(*Node)
	walk = func(n *Node) {
		fn(n)
		for _, c := range n.Children {
			walk(c)
		}
	}
	for _, c := range t.Root.Children {
		walk(c)
	}
}

// Find resolves a command path such as ["mothership", "connect"] to its node.
func (t *Tree) Find(path []string) (*Node, bool) {
	cur := t.Root
	for _, want := range path {
		var next *Node
		for _, c := range cur.Children {
			if c.Use == want {
				next = c
				break
			}
		}
		if next == nil {
			return cur, false
		}
		cur = next
	}
	return cur, true
}

// Resolve splits a requested command path into the deepest existing ancestor
// and the segments that do not exist yet.
func (t *Tree) Resolve(path []string) (parent *Node, missing []string) {
	cur := t.Root
	for i, want := range path {
		var next *Node
		for _, c := range cur.Children {
			if c.Use == want {
				next = c
				break
			}
		}
		if next == nil {
			return cur, path[i:]
		}
		cur = next
	}
	return cur, nil
}

// Node looks a node up by package and variable name.
func (t *Tree) Node(pkg, varName string) (*Node, bool) {
	n, ok := t.nodes[pkg+"."+varName]
	return n, ok
}

// Paths returns every command path in the tree, sorted.
func (t *Tree) Paths() []string {
	var out []string
	t.Walk(func(n *Node) { out = append(out, n.String()) })
	sort.Strings(out)
	return out
}

// LoadTree parses cmd/ and internal/cli/ and reconstructs the cobra tree.
//
// The tree is assembled in three idioms and all three have to be understood:
// cmd/root.go attaches each service's Command() factory, a factory attaches its
// own direct subcommands from a local variable, and deeper levels attach in the
// child file's init().
func LoadTree(root string) (*Tree, error) {
	dirs := []string{filepath.Join(root, "cmd")}
	cliRoot := filepath.Join(root, "internal", "cli")
	entries, err := os.ReadDir(cliRoot)
	if err != nil {
		return nil, err
	}
	dirs = append(dirs, cliRoot)
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, filepath.Join(cliRoot, e.Name()))
		}
	}

	t := &Tree{nodes: map[string]*Node{}}
	// factory maps a package to the command its Command() function returns, and
	// edges are collected as (parentID, childID) so they can be linked once
	// every declaration is known.
	factory := map[string]*Node{}
	type edge struct{ parent, child string }
	var edges []edge
	// importAlias maps, per file, an imported package's local name to its
	// directory, so a cross-package AddCommand in root.go can be resolved.
	imports := map[string]map[string]string{}

	for _, dir := range dirs {
		pkg, err := filepath.Rel(root, dir)
		if err != nil {
			return nil, err
		}
		files, err := filepath.Glob(filepath.Join(dir, "*.go"))
		if err != nil {
			return nil, err
		}
		sort.Strings(files)

		for _, path := range files {
			if strings.HasSuffix(path, "_test.go") {
				continue
			}
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", path, err)
			}
			rel, _ := filepath.Rel(root, path)

			imports[rel] = map[string]string{}
			for _, imp := range f.Imports {
				p := strings.Trim(imp.Path.Value, `"`)
				name := filepath.Base(p)
				if imp.Name != nil {
					name = imp.Name.Name
				}
				imports[rel][name] = p
			}

			// Package-level command declarations.
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
					if use, short, ok := cobraCommandLit(vs.Values[0]); ok {
						n := &Node{Pkg: pkg, Var: vs.Names[0].Name, Use: use, Short: short, File: rel}
						t.nodes[n.ID()] = n
					}
				}
			}

			// Declarations and edges inside function bodies: the Command()
			// factories and the init() blocks.
			for _, decl := range f.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Body == nil {
					continue
				}
				isFactory := fn.Recv == nil && fn.Name.Name == "Command"

				ast.Inspect(fn.Body, func(nd ast.Node) bool {
					switch node := nd.(type) {
					case *ast.AssignStmt:
						if !isFactory || len(node.Lhs) != 1 || len(node.Rhs) != 1 {
							return true
						}
						id, ok := node.Lhs[0].(*ast.Ident)
						if !ok {
							return true
						}
						if use, short, ok := cobraCommandLit(node.Rhs[0]); ok {
							n := &Node{Pkg: pkg, Var: id.Name, Use: use, Short: short, File: rel, Factory: true}
							t.nodes[n.ID()] = n
							factory[pkg] = n
						}
					case *ast.CallExpr:
						sel, ok := node.Fun.(*ast.SelectorExpr)
						if !ok || sel.Sel.Name != "AddCommand" || len(node.Args) != 1 {
							return true
						}
						recv, ok := sel.X.(*ast.Ident)
						if !ok {
							return true
						}
						parentID := pkg + "." + recv.Name

						switch arg := node.Args[0].(type) {
						case *ast.Ident:
							edges = append(edges, edge{parentID, pkg + "." + arg.Name})
						case *ast.CallExpr:
							// rootCmd.AddCommand(idp.Command())
							cs, ok := arg.Fun.(*ast.SelectorExpr)
							if !ok || cs.Sel.Name != "Command" {
								return true
							}
							alias, ok := cs.X.(*ast.Ident)
							if !ok {
								return true
							}
							target, ok := imports[rel][alias.Name]
							if !ok {
								return true
							}
							childPkg := strings.TrimPrefix(target, modulePath+"/")
							edges = append(edges, edge{parentID, "factory:" + childPkg})
						}
					}
					return true
				})
			}
		}
	}

	rootNode, ok := t.nodes["cmd.rootCmd"]
	if !ok {
		return nil, fmt.Errorf("no rootCmd declaration in cmd/")
	}
	t.Root = rootNode

	for _, e := range edges {
		parent, ok := t.nodes[e.parent]
		if !ok {
			continue
		}
		var child *Node
		if strings.HasPrefix(e.child, "factory:") {
			child, ok = factory[strings.TrimPrefix(e.child, "factory:")]
			if !ok {
				continue
			}
		} else if child, ok = t.nodes[e.child]; !ok {
			continue
		}
		child.Parent = parent
		parent.Children = append(parent.Children, child)
	}
	t.Walk(func(n *Node) {
		sort.Slice(n.Children, func(i, j int) bool { return n.Children[i].Use < n.Children[j].Use })
	})
	return t, nil
}

// modulePath is this repository's module path, needed to turn an import path
// back into a directory relative to the repo root.
const modulePath = "github.com/moovit-sp-gmbh/hcloud-cli-repl"

// cobraCommandLit reports whether an expression is a &cobra.Command{...}
// literal and returns its Use and Short fields. Use is trimmed to its first
// word, since a Use may carry a usage hint ("create <name>").
func cobraCommandLit(e ast.Expr) (use, short string, ok bool) {
	unary, isUnary := e.(*ast.UnaryExpr)
	if !isUnary || unary.Op != token.AND {
		return "", "", false
	}
	lit, isLit := unary.X.(*ast.CompositeLit)
	if !isLit || !isSelector(lit.Type, "cobra", "Command") {
		return "", "", false
	}
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}
		v, _ := stringLit(kv.Value)
		switch key.Name {
		case "Use":
			if fields := strings.Fields(v); len(fields) > 0 {
				use = fields[0]
			}
		case "Short":
			short = v
		}
	}
	return use, short, use != ""
}
