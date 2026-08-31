package codegen

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

// copyRepo makes a working copy of the repository so an Apply test can build
// the result without touching the real tree.
func copyRepo(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	dst := t.TempDir()

	cmd := exec.Command("rsync", "-a",
		"--exclude=.git", "--exclude=bin", "--exclude=hcloud-cli",
		"--exclude=db.sqlite", "--exclude=*/db.sqlite", "--exclude=debug.log",
		root+"/", dst+"/")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("rsync: %v\n%s", err, out)
	}
	return dst
}

// goBuild builds and vets the tree. Vet matters as much as the build here: a
// generated fmt.Errorf with the wrong number of arguments compiles fine and is
// exactly the kind of mistake a code generator makes at scale.
func goBuild(t *testing.T, root string) {
	t.Helper()
	build := exec.Command("go", "build", "./...")
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build ./... failed after generation:\n%s", out)
	}
	vet := exec.Command("go", "vet", "./...")
	vet.Dir = root
	if out, err := vet.CombinedOutput(); err != nil {
		t.Fatalf("go vet ./... failed after generation:\n%s", out)
	}
}

// matches reports whether src contains a match for pattern.
func matches(src, pattern string) bool {
	return regexp.MustCompile(pattern).MatchString(src)
}

func mustRead(t *testing.T, root, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// TestApplyNewServicePackage generates a command for dali, which has no command
// package at all. That exercises the whole chain: a new package, its
// registration on the root command, a client call, and the manifest entry.
func TestApplyNewServicePackage(t *testing.T) {
	root := copyRepo(t)

	ix, err := LoadIndex(root)
	if err != nil {
		t.Fatal(err)
	}
	schema, err := LoadContextSchema(root)
	if err != nil {
		t.Fatal(err)
	}
	op, ok := ix.Op("daliClient", "PostV1AvatarOrgOrgName")
	if !ok {
		t.Fatal("PostV1AvatarOrgOrgName not in the dali client")
	}

	r := Derive(op, ix.Clients["daliClient"], LoadSpecs(root), schema)
	r.Command.Path = []string{"dali", "avatar", "org", "set"}
	r.Command.Var = CommandVar(r.Command.Path)
	r.Command.File = "avatar.go"
	r.Command.Short = "Set an organization's avatar"

	res, err := Apply(root, r, ApplyOptions{Now: time.Date(2026, 9, 1, 9, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatal(err)
	}
	goBuild(t, root)

	// The service package and both intermediate groups must exist.
	group := mustRead(t, root, "internal/cli/dali/dali.go")
	if !strings.Contains(group, "func Command() *cobra.Command") {
		t.Error("dali.go has no Command factory")
	}
	if !strings.Contains(group, `PersistentFlags().String("server"`) {
		t.Error("dali.go does not declare the --server flag")
	}

	rootGo := mustRead(t, root, "cmd/root.go")
	if !strings.Contains(rootGo, "rootCmd.AddCommand(dali.Command())") {
		t.Error("dali not registered on the root command")
	}

	avatar := mustRead(t, root, "internal/cli/dali/avatar.go")
	for _, want := range []string{
		"var avatarCommand = &cobra.Command{",
		"var avatarOrgCommand = &cobra.Command{",
		"var avatarOrgSetCommand = &cobra.Command{",
		"clients.NewDaliClient(cmd)",
		"PostV1AvatarOrgOrgNameWithResponse(",
		"command.SetResult(cmd, text)",
	} {
		if !strings.Contains(avatar, want) {
			t.Errorf("avatar.go missing %q\n%s", want, avatar)
		}
	}

	// The tree must now resolve the new path end to end.
	tree, err := LoadTree(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := tree.Find([]string{"dali", "avatar", "org", "set"}); !ok {
		t.Errorf("new command not reachable in the tree; have:\n%s", strings.Join(tree.Paths(), "\n"))
	}

	m, err := LoadManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	entry, ok := m.Lookup("daliClient", "PostV1AvatarOrgOrgName")
	if !ok {
		t.Fatal("manifest has no entry for the generated command")
	}
	if entry.Command != "dali avatar org set" || entry.Source != SourceGenerated {
		t.Errorf("manifest entry = %+v", entry)
	}
	if entry.Recipe == "" {
		t.Error("manifest entry does not name its recipe")
	}
	if _, err := os.Stat(filepath.Join(root, entry.Recipe)); err != nil {
		t.Errorf("recipe not saved: %v", err)
	}
	if len(res.Columns) != 0 {
		t.Errorf("unexpected new context columns: %+v", res.Columns)
	}
}

// TestApplyPersistsNewContextColumn covers the other half: keeping part of a
// response in the context table, which needs a migration, a regenerated sqlc
// model and a new field in ContextService.Push.
func TestApplyPersistsNewContextColumn(t *testing.T) {
	if _, err := exec.LookPath("sqlc"); err != nil {
		t.Skip("sqlc not installed")
	}
	root := copyRepo(t)

	ix, err := LoadIndex(root)
	if err != nil {
		t.Fatal(err)
	}
	schema, err := LoadContextSchema(root)
	if err != nil {
		t.Fatal(err)
	}
	op, ok := ix.Op("mothershipClient", "PostV1Hello")
	if !ok {
		t.Fatal("PostV1Hello not in the mothership client")
	}

	r := Derive(op, ix.Clients["mothershipClient"], LoadSpecs(root), schema)
	r.Command.Path = []string{"mothership", "hello", "status"}
	r.Command.Var = "helloStatusCommand"
	r.Command.File = "hello_status.go"
	r.Command.Short = "Say hello and remember the agent id"
	// AgentDto.UnderscoreId is the agent's id; keep it under a new column. The
	// field is deliberately spelled AgentId here: Apply has to canonicalize it
	// to AgentID, because that is what sqlc will generate for agent_id.
	r.Response.Persist = []RecipePersist{{
		From: "JSON200:UnderscoreId", Context: "AgentId", Column: "agent_id",
	}}

	res, err := Apply(root, r, ApplyOptions{Now: time.Date(2026, 9, 1, 10, 15, 0, 0, time.UTC)})
	if err != nil {
		t.Fatal(err)
	}
	goBuild(t, root)

	if len(res.Columns) != 1 || res.Columns[0].Column != "agent_id" || res.Columns[0].GoName != "AgentID" {
		t.Fatalf("columns = %+v", res.Columns)
	}

	// The saved recipe must record the canonical spelling, or a later regen
	// would reintroduce the mismatch.
	saved, err := LoadRecipe(filepath.Join(root, r.RecipePath()))
	if err != nil {
		t.Fatal(err)
	}
	if got := saved.Response.Persist[0].Context; got != "AgentID" {
		t.Errorf("saved recipe context field = %q, want AgentID", got)
	}

	migration := mustRead(t, root, "internal/db/migrations/20260901101500_agent_id_ctx.sql")
	if !strings.Contains(migration, "ADD COLUMN agent_id TEXT;") || !strings.Contains(migration, "DROP COLUMN agent_id;") {
		t.Errorf("migration:\n%s", migration)
	}

	// sqlc must have picked the column up, and Push must copy it, or the value
	// is written to the struct and then dropped on the way to the database.
	// gofmt aligns struct fields, so these are matched loosely rather than on
	// exact spacing.
	if !matches(mustRead(t, root, "internal/db/models.go"), `AgentID\s+sql\.NullString`) {
		t.Error("sqlc did not regenerate db.Context with AgentID")
	}
	if !matches(mustRead(t, root, "internal/service.go"), `AgentID:\s+arg\.AgentID,`) {
		t.Error("ContextService.Push does not carry AgentID")
	}

	cmdSrc := mustRead(t, root, "internal/cli/mothership/hello_status.go")
	for _, want := range []string{
		"if res.JSON200 != nil {",
		"latest.AgentID = sql.NullString{Valid: true, String: res.JSON200.UnderscoreId}",
		"service.Push(latest)",
	} {
		if !strings.Contains(cmdSrc, want) {
			t.Errorf("generated command missing %q\n%s", want, cmdSrc)
		}
	}

	// The migration has to actually apply against sqlite.
	if _, err := exec.LookPath("goose"); err == nil {
		db := filepath.Join(t.TempDir(), "test.sqlite")
		cmd := exec.Command("goose", "sqlite3", db, "up", "-dir", "internal/db/migrations")
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Errorf("goose up failed:\n%s", out)
		}
	}
}

// TestApplyDryRunWritesNothing makes sure the preview is really a preview.
func TestApplyDryRunWritesNothing(t *testing.T) {
	root := copyRepo(t)

	ix, err := LoadIndex(root)
	if err != nil {
		t.Fatal(err)
	}
	schema, err := LoadContextSchema(root)
	if err != nil {
		t.Fatal(err)
	}
	op, _ := ix.Op("mothershipClient", "PostV1Hello")

	r := Derive(op, ix.Clients["mothershipClient"], LoadSpecs(root), schema)
	r.Command.Path = []string{"mothership", "hello", "status"}
	r.Command.Var = "helloStatusCommand"
	r.Command.File = "hello_status.go"

	before := mustRead(t, root, "internal/cli/mothership/mothership.go")
	res, err := Apply(root, r, ApplyOptions{DryRun: true, Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Created) == 0 {
		t.Error("dry run reported nothing")
	}
	if _, err := os.Stat(filepath.Join(root, "internal/cli/mothership/hello_status.go")); err == nil {
		t.Error("dry run wrote the command file")
	}
	if after := mustRead(t, root, "internal/cli/mothership/mothership.go"); after != before {
		t.Error("dry run edited the parent file")
	}
	if _, err := os.Stat(filepath.Join(root, ManifestPath)); err == nil {
		if m, _ := LoadManifest(root); len(m.Endpoints) > 0 {
			if _, ok := m.Lookup("mothershipClient", "PostV1Hello"); ok {
				// PostV1Hello is already implemented by hand, so an entry may
				// legitimately exist; only a generated source would be wrong.
				e, _ := m.Lookup("mothershipClient", "PostV1Hello")
				if e.Source == SourceGenerated {
					t.Error("dry run rewrote the manifest")
				}
			}
		}
	}
}

// TestRegenIsIdempotent replays a recipe over the command it already generated.
// Recipes only earn their keep if replaying one after an emitter change rewrites
// the command in place, so a second pass has to be a no-op rather than a second
// declaration.
func TestRegenIsIdempotent(t *testing.T) {
	root := copyRepo(t)

	ix, err := LoadIndex(root)
	if err != nil {
		t.Fatal(err)
	}
	schema, err := LoadContextSchema(root)
	if err != nil {
		t.Fatal(err)
	}
	op, ok := ix.Op("mothershipClient", "DeleteV1Agent")
	if !ok {
		t.Fatal("DeleteV1Agent not in the mothership client")
	}

	r := Derive(op, ix.Clients["mothershipClient"], LoadSpecs(root), schema)
	r.Command.Path = []string{"mothership", "agent", "retire"}
	r.Command.Var = "agentRetireCommand"
	r.Command.File = "agent_retire.go"
	// DeleteV1Agent declares no authorization parameter even though it needs
	// credentials, so the header has to be added explicitly and sent through a
	// RequestEditorFn — the same thing internal/cli/mothership/agent.go does by
	// hand.
	r.Args = append(r.Args, AuthorizationArg())

	if _, err := Apply(root, r, ApplyOptions{Now: time.Now()}); err != nil {
		t.Fatal(err)
	}
	goBuild(t, root)

	const rel = "internal/cli/mothership/agent_retire.go"
	first := mustRead(t, root, rel)
	for _, want := range []string{
		`agentRetireCommand.Flags().String("token"`,
		"func(_ context.Context, req *http.Request) error {",
		`req.Header.Set("authorization", jwt)`,
	} {
		if !strings.Contains(first, want) {
			t.Errorf("generated command missing %q\n%s", want, first)
		}
	}
	parentFirst := mustRead(t, root, "internal/cli/mothership/agent.go")
	manifestFirst := mustRead(t, root, ManifestPath)

	// Replay from the saved recipe, the way `codegen regen` does.
	saved, err := LoadRecipe(filepath.Join(root, r.RecipePath()))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Apply(root, saved, ApplyOptions{Now: time.Now()}); err != nil {
		t.Fatalf("regen failed: %v", err)
	}
	goBuild(t, root)

	if second := mustRead(t, root, rel); second != first {
		t.Errorf("regeneration changed the command file\n--- first ---\n%s\n--- second ---\n%s", first, second)
	}
	if second := mustRead(t, root, "internal/cli/mothership/agent.go"); second != parentFirst {
		t.Error("regeneration re-registered the subcommand")
	}
	if second := mustRead(t, root, ManifestPath); second != manifestFirst {
		t.Error("regeneration duplicated the manifest entry")
	}

	// Exactly one declaration and one flag registration each.
	if n := strings.Count(first, "var agentRetireCommand ="); n != 1 {
		t.Errorf("%d declarations of agentRetireCommand", n)
	}

	// The same has to hold when a new command group shares the file, and so the
	// init(), with the command: that is where an append would keep reordering
	// the block.
	shared := Derive(op, ix.Clients["mothershipClient"], LoadSpecs(root), schema)
	shared.Command.Path = []string{"mothership", "fleet", "retire"}
	shared.Command.Var = "fleetRetireCommand"
	shared.Command.File = "fleet.go"
	shared.Args = append(shared.Args, AuthorizationArg())
	if _, err := Apply(root, shared, ApplyOptions{Now: time.Now()}); err != nil {
		t.Fatal(err)
	}
	goBuild(t, root)
	sharedFirst := mustRead(t, root, "internal/cli/mothership/fleet.go")
	if !strings.Contains(sharedFirst, "fleetCommand.AddCommand(fleetRetireCommand)") {
		t.Errorf("group and command not wired in a shared file:\n%s", sharedFirst)
	}
	if _, err := Apply(root, shared, ApplyOptions{Now: time.Now()}); err != nil {
		t.Fatal(err)
	}
	goBuild(t, root)
	if sharedSecond := mustRead(t, root, "internal/cli/mothership/fleet.go"); sharedSecond != sharedFirst {
		t.Errorf("regeneration reordered a shared init()\n--- first ---\n%s\n--- second ---\n%s", sharedFirst, sharedSecond)
	}

	// A changed recipe must actually take effect on replay.
	saved.Command.Short = "Retire an agent from the mothership"
	if _, err := Apply(root, saved, ApplyOptions{Now: time.Now()}); err != nil {
		t.Fatal(err)
	}
	goBuild(t, root)
	if third := mustRead(t, root, rel); !strings.Contains(third, "Retire an agent from the mothership") {
		t.Errorf("edited recipe had no effect on replay:\n%s", third)
	}
}
