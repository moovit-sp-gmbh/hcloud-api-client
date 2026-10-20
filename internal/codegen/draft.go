package codegen

import (
	"fmt"
	"sort"
	"strings"
)

// DraftPaths assigns every operation a command path that is unique across the
// batch and against the commands that already exist.
//
// The single-operation heuristic in SuggestPath is enough for the wizard, where
// a person sees the proposal and corrects it. Generating hundreds of commands at
// once needs more: two endpoints routinely reduce to the same path because the
// heuristic drops path parameters, and a path that is also another command's
// group would collide on the variable name. Both are resolved here, by adding
// back the distinction the heuristic threw away.
//
// taken lists the command paths already in the tree. The result is keyed by
// Operation.Key.
func DraftPaths(ops []*Operation, clients map[string]*Client, taken []string) map[string][]string {
	reserved := map[string]bool{}
	for _, t := range taken {
		reserved[t] = true
		// A command's ancestors are groups, and a leaf may not sit on one.
		for _, prefix := range prefixesOf(strings.Fields(t)) {
			reserved[prefix] = true
		}
	}

	candidates := make([]*candidate, 0, len(ops))
	for _, op := range ops {
		// Start from the nouns alone rather than from SuggestPath: the wizard's
		// heuristic bakes in a verb, and a verb added before the clash is known
		// leaves no room to say which of two endpoints is the collection.
		candidates = append(candidates, &candidate{op: op, path: CommandNouns(op)})
	}
	// A stable order keeps the result reproducible run to run.
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].op.Client != candidates[j].op.Client {
			return candidates[i].op.Client < candidates[j].op.Client
		}
		if candidates[i].op.Path != candidates[j].op.Path {
			return candidates[i].op.Path < candidates[j].op.Path
		}
		return candidates[i].op.Method < candidates[j].op.Method
	})

	// Refine until nothing collides. Each round adds one more distinguishing
	// piece to every candidate still in a clash, and the pieces are chosen by
	// comparing the clashing endpoints rather than each one on its own.
	const maxRounds = 6
	for round := 0; round < maxRounds; round++ {
		byPath := map[string][]*candidate{}
		groups := map[string]bool{}
		for _, c := range candidates {
			key := strings.Join(c.path, " ")
			byPath[key] = append(byPath[key], c)
			for _, prefix := range prefixesOf(c.path) {
				groups[prefix] = true
			}
		}

		changed := false
		for key, group := range byPath {
			if len(group) == 1 && !reserved[key] && !groups[key] {
				continue
			}
			for _, c := range group {
				if next, ok := refine(c, group); ok {
					c.path, c.level = next, c.level+1
					changed = true
				}
			}
		}
		if !changed {
			break
		}
	}

	// Anything still colliding after refinement gets a numeric suffix, so the
	// draft always produces a buildable tree rather than failing late.
	seen := map[string]bool{}
	for k := range reserved {
		seen[k] = true
	}
	out := map[string][]string{}
	for _, c := range candidates {
		path := c.path
		for n := 2; seen[strings.Join(path, " ")]; n++ {
			path = append(append([]string{}, c.path[:len(c.path)-1]...),
				fmt.Sprintf("%s-%d", c.path[len(c.path)-1], n))
		}
		seen[strings.Join(path, " ")] = true
		out[c.op.Key()] = path
	}
	return out
}

// candidate is one operation and the command path assigned to it so far.
type candidate struct {
	op    *Operation
	path  []string
	level int // how much disambiguation has been applied
}

// refine adds the next piece of disambiguation to a clashing command path.
// group is every candidate currently sharing that path, because what
// distinguishes two endpoints is how they differ from each other.
func refine(c *candidate, group []*candidate) ([]string, bool) {
	switch c.level {
	case 0:
		verb := draftVerb(c, group)
		if verb == "" || (len(c.path) > 0 && c.path[len(c.path)-1] == verb) {
			return nil, false
		}
		return append(append([]string{}, c.path...), verb), true

	case 1:
		// Two endpoints differing only by a path parameter in the middle —
		// DELETE /assets/subscribe against DELETE /assets/{assetId}/subscribe.
		// The one taking fewer parameters acts on the collection.
		if len(c.op.PathParams) != fewestParams(group) {
			return nil, false
		}
		leaf := c.path[len(c.path)-1]
		if strings.HasSuffix(leaf, "-all") {
			return nil, false
		}
		out := append([]string{}, c.path...)
		out[len(out)-1] = leaf + "-all"
		return out, true

	case 2:
		// Fall back to naming the resource the endpoint is scoped to.
		if params := c.op.PathParams; len(params) > 0 {
			last := KebabCase(params[len(params)-1].Name)
			if last != "" && (len(c.path) == 0 || c.path[len(c.path)-1] != "by-"+last) {
				return append(append([]string{}, c.path...), "by-"+last), true
			}
		}
		return nil, false
	}
	return nil, false
}

// draftVerb names the action, telling a collection apart from one resource by
// how the clashing endpoints differ: of two GETs on the same nouns, the one
// taking fewer path parameters is the list.
func draftVerb(c *candidate, group []*candidate) string {
	if strings.ToUpper(c.op.Method) != "GET" {
		return MethodVerb(c.op.Method)
	}
	if len(group) > 1 {
		if len(c.op.PathParams) == fewestParams(group) {
			return "list"
		}
		return "get"
	}
	if AddressesItem(c.op) {
		return "get"
	}
	return "list"
}

func fewestParams(group []*candidate) int {
	fewest := -1
	for _, c := range group {
		if fewest < 0 || len(c.op.PathParams) < fewest {
			fewest = len(c.op.PathParams)
		}
	}
	return fewest
}

// AddressesItem reports whether an endpoint acts on one named resource rather
// than on a collection, which is the difference between `spaces get` and
// `spaces list`.
//
// Every endpoint here is org-scoped, so the organization's own parameter says
// nothing about which of the two it is.
func AddressesItem(op *Operation) bool {
	for _, p := range op.PathParams {
		if !strings.EqualFold(p.Name, "orgName") {
			return true
		}
	}
	return false
}

// prefixesOf returns every strict prefix of a command path, which are the
// groups the command hangs under.
func prefixesOf(path []string) []string {
	var out []string
	for i := 1; i < len(path); i++ {
		out = append(out, strings.Join(path[:i], " "))
	}
	return out
}

// DraftVars assigns each drafted command a variable name that does not collide
// inside its package.
//
// A unique command path does not imply a unique variable: CommandVar drops the
// service segment, so `idp patch` reduces to patchCommand, which internal/cli/idp
// already declares by hand for `idp user patch`. A clash there is a compile
// error in a package nobody touched, so the service goes back into the name.
func DraftVars(paths map[string][]string, taken map[string]map[string]bool) map[string]string {
	// Deterministic order, so the same input always names things the same way.
	keys := make([]string, 0, len(paths))
	for key := range paths {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	used := map[string]map[string]bool{}
	for pkg, names := range taken {
		used[pkg] = map[string]bool{}
		for name := range names {
			used[pkg][name] = true
		}
	}

	out := make(map[string]string, len(paths))
	for _, key := range keys {
		path := paths[key]
		pkg := path[0]
		if used[pkg] == nil {
			used[pkg] = map[string]bool{}
		}

		name := CommandVar(path)
		if used[pkg][name] {
			// Put the service back: idp patch -> idpPatchCommand.
			name = CommandVar(append([]string{""}, path...))
		}
		for n := 2; used[pkg][name]; n++ {
			name = fmt.Sprintf("%s%dCommand", strings.TrimSuffix(CommandVar(path), "Command"), n)
		}
		used[pkg][name] = true
		out[key] = name
	}
	return out
}

// DraftFile is the file a drafted command is written to: one file per command
// group, the way internal/cli/mothership/agent.go holds agentCommand and its
// subcommand.
func DraftFile(path []string) string {
	// path[0] is the service, which the package already names.
	group := path[1 : len(path)-1]
	if len(group) == 0 {
		return path[len(path)-1] + ".go"
	}
	return strings.Join(group, "-") + ".go"
}

// DraftResult records what drafting one endpoint did.
type DraftResult struct {
	Op      *Operation
	Command string
	Err     error
	Result  *ApplyResult
}

// DraftOptions controls a bulk draft.
type DraftOptions struct {
	// Clients limits the draft to these client packages; empty means all.
	Clients []string
	// Limit caps how many commands are generated; zero means no cap.
	Limit  int
	DryRun bool
	// Progress is called before each command is generated.
	Progress func(i, total int, command string)
}

// Draft generates a command for every endpoint that has none.
//
// This is the non-interactive counterpart of the wizard: every answer is taken
// from Derive's defaults and from DraftPaths, so the result is a starting point
// to edit rather than a considered design. Each command still lands with its
// recipe, so a command whose shape is wrong is fixed by editing the recipe and
// running regen, not by hand-patching generated code.
func Draft(root string, opts DraftOptions) ([]DraftResult, error) {
	ix, err := LoadIndex(root)
	if err != nil {
		return nil, err
	}
	manifest, err := LoadManifest(root)
	if err != nil {
		return nil, err
	}
	schema, err := LoadContextSchema(root)
	if err != nil {
		return nil, err
	}
	tree, err := LoadTree(root)
	if err != nil {
		return nil, err
	}
	specs := LoadSpecs(root)

	var taken []string
	// Variable names are tracked per package, including the groups a factory
	// declares, because a drafted command lands beside them.
	takenVars := map[string]map[string]bool{}
	tree.Walk(func(n *Node) {
		taken = append(taken, n.String())
		service := n.Pkg[strings.LastIndex(n.Pkg, "/")+1:]
		if takenVars[service] == nil {
			takenVars[service] = map[string]bool{}
		}
		takenVars[service][n.Var] = true
	})

	var ops []*Operation
	for _, cc := range Compare(ix, manifest, opts.Clients).Clients {
		for _, op := range cc.Missing {
			if op.Callable() {
				ops = append(ops, op)
			}
		}
	}
	sort.Slice(ops, func(i, j int) bool {
		if ops[i].Client != ops[j].Client {
			return ops[i].Client < ops[j].Client
		}
		if ops[i].Path != ops[j].Path {
			return ops[i].Path < ops[j].Path
		}
		return ops[i].Method < ops[j].Method
	})

	paths := DraftPaths(ops, ix.Clients, taken)
	vars := DraftVars(paths, takenVars)
	if opts.Limit > 0 && len(ops) > opts.Limit {
		ops = ops[:opts.Limit]
	}

	results := make([]DraftResult, 0, len(ops))
	for i, op := range ops {
		path := paths[op.Key()]
		command := strings.Join(path, " ")
		if opts.Progress != nil {
			opts.Progress(i+1, len(ops), command)
		}

		r := Derive(op, ix.Clients[op.Client], specs, schema)
		r.Command.Path = path
		r.Command.Var = vars[op.Key()]
		r.Command.File = DraftFile(path)

		res, err := Apply(root, r, ApplyOptions{DryRun: opts.DryRun, Index: ix})
		// One endpoint that cannot be drafted must not stop the other 300, so
		// the failure is recorded and the run continues.
		results = append(results, DraftResult{Op: op, Command: command, Err: err, Result: res})
	}
	return results, nil
}
