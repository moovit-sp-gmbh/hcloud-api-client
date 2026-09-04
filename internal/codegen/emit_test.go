package codegen

import (
	"os"
	"strings"
	"testing"
)

// loadFixtures gives a test the real clients, specs and context schema.
func loadFixtures(t *testing.T) (*Index, *ContextSchema, *Specs, map[string]string) {
	t.Helper()
	ix, err := LoadIndex("../..")
	if err != nil {
		t.Fatal(err)
	}
	schema, err := LoadContextSchema("../..")
	if err != nil {
		t.Fatal(err)
	}
	ctors, err := loadConstructors("../..")
	if err != nil {
		t.Fatal(err)
	}
	return ix, schema, LoadSpecs("../.."), ctors
}

// TestEmitReproducesStats is the emitter's acceptance test: generating
// mothership stats from its derived recipe must produce the RunE body that was
// written by hand in internal/cli/mothership/stats.go. That file is the
// cleanest instance of the full idiom, so matching it means the templates agree
// with the convention rather than with themselves.
func TestEmitReproducesStats(t *testing.T) {
	ix, schema, specs, ctors := loadFixtures(t)
	op, ok := ix.Op("mothershipClient", "GetV1OrgOrgNameStats")
	if !ok {
		t.Fatal("GetV1OrgOrgNameStats not in the index")
	}

	r := Derive(op, ix.Clients["mothershipClient"], specs, schema)
	// Derive proposes the path from the endpoint; the real command is one level
	// shallower because "stats" is already unambiguous.
	r.Command.Path = []string{"mothership", "stats"}
	r.Command.Var = "statsCommand"
	r.Command.File = "stats.go"
	r.Command.Short = "Retrieve usage statistics for a specific organization"

	g, err := Emit(r, op, schema, ctors)
	if err != nil {
		t.Fatal(err)
	}

	original, err := os.ReadFile("../cli/mothership/stats.go")
	if err != nil {
		t.Fatal(err)
	}

	gotBody := runEBody(t, string(g.Source))
	wantBody := runEBody(t, string(original))
	if gotBody != wantBody {
		t.Errorf("generated RunE body differs from the hand-written stats.go\n--- generated ---\n%s\n--- stats.go ---\n%s", gotBody, wantBody)
	}

	// stats.go registers no flags at all even though it reads two, which is the
	// bug the generator closes by construction.
	if len(g.Flags) != 2 {
		t.Errorf("flags = %v, want two registrations", g.Flags)
	}
	joined := strings.Join(g.Flags, "\n")
	for _, want := range []string{`statsCommand.Flags().String("orgName"`, `statsCommand.Flags().String("token"`} {
		if !strings.Contains(joined, want) {
			t.Errorf("missing flag registration %q in:\n%s", want, joined)
		}
	}
}

// runEBody extracts the statements inside a command's RunE literal, normalized
// for comparison.
func runEBody(t *testing.T, source string) string {
	t.Helper()
	const marker = "RunE: func(cmd *cobra.Command, args []string) error {"
	start := strings.Index(source, marker)
	if start < 0 {
		t.Fatalf("no RunE literal found in:\n%s", source)
	}
	rest := source[start+len(marker):]

	depth := 1
	end := -1
	for i, r := range rest {
		switch r {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				end = i
			}
		}
		if end >= 0 {
			break
		}
	}
	if end < 0 {
		t.Fatal("unterminated RunE literal")
	}

	var out []string
	for _, line := range strings.Split(rest[:end], "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return strings.Join(out, "\n")
}

// TestOptionalUnsetBehaviour covers the two things an optional field can do
// when its flag is not given: be left out so the API applies its own default,
// or be sent carrying a default the command states itself.
func TestOptionalUnsetBehaviour(t *testing.T) {
	ix, schema, specs, ctors := loadFixtures(t)
	op, ok := ix.Op("high5Client", "PostV1OrgOrgNameSpaces")
	if !ok {
		t.Fatal("PostV1OrgOrgNameSpaces not in the high5 client")
	}

	// enableDev is *bool and optional, which is the case that prompted this.
	find := func(r *Recipe, field string) *RecipeBodyField {
		for i := range r.Body.Fields {
			if r.Body.Fields[i].Field == field {
				return &r.Body.Fields[i]
			}
		}
		t.Fatalf("no body field %s; the DTO must have changed", field)
		return nil
	}

	base := func() *Recipe {
		r := Derive(op, ix.Clients["high5Client"], specs, schema)
		r.Command.Path = []string{"high5", "space", "create"}
		r.Command.Var = "spaceCreateCommand"
		r.Command.File = "spaces.go"
		return r
	}

	t.Run("omit leaves the field out", func(t *testing.T) {
		r := base()
		enableDev := find(r, "EnableDev")
		if enableDev.Unset != UnsetOmit {
			t.Errorf("Derive should default to omit, got %q", enableDev.Unset)
		}

		g, err := Emit(r, op, schema, ctors)
		if err != nil {
			t.Fatal(err)
		}
		src := string(g.Source)
		if !strings.Contains(src, `if cmd.Flags().Changed("enableDev") {`) {
			t.Errorf("omit should guard on Changed:\n%s", src)
		}
		if !strings.Contains(src, `spaceCreateCommand.Flags().Bool("enableDev", false, `) {
			t.Errorf("flag not registered with the type zero:\n%s", strings.Join(g.Flags, "\n"))
		}
	})

	t.Run("default always sends the field", func(t *testing.T) {
		r := base()
		enableDev := find(r, "EnableDev")
		enableDev.Unset = UnsetDefault
		enableDev.Default = "true"

		g, err := Emit(r, op, schema, ctors)
		if err != nil {
			t.Fatal(err)
		}
		src := string(g.Source)
		if strings.Contains(src, `Changed("enableDev")`) {
			t.Errorf("a defaulted field must not be guarded on Changed:\n%s", src)
		}
		for _, want := range []string{
			`enableDev, err := cmd.Flags().GetBool("enableDev")`,
			"body.EnableDev = &enableDev",
		} {
			if !strings.Contains(src, want) {
				t.Errorf("missing %q:\n%s", want, src)
			}
		}
		// The default has to reach the flag registration, or --help lies and an
		// unset flag sends the wrong value.
		if !strings.Contains(strings.Join(g.Flags, "\n"), `Flags().Bool("enableDev", true, `) {
			t.Errorf("default not carried into the flag:\n%s", strings.Join(g.Flags, "\n"))
		}
	})

	t.Run("recipe defaults are plain values, not Go source", func(t *testing.T) {
		r := base()
		name := find(r, "Name")
		name.Resolve = []string{FromFlag}
		name.Prompt = ""
		name.Context = ""
		name.Persist = false
		name.Default = "my-space"

		g, err := Emit(r, op, schema, ctors)
		if err != nil {
			t.Fatal(err)
		}
		// The recipe says my-space; the generated code has to say "my-space".
		if !strings.Contains(strings.Join(g.Flags, "\n"), `Flags().String("name", "my-space", `) {
			t.Errorf("string default not quoted for Go:\n%s", strings.Join(g.Flags, "\n"))
		}
	})

	t.Run("a default that does not fit the type is an error", func(t *testing.T) {
		r := base()
		enableDev := find(r, "EnableDev")
		enableDev.Unset = UnsetDefault
		enableDev.Default = "yesplease"

		_, err := Emit(r, op, schema, ctors)
		if err == nil {
			t.Fatal("a non-boolean default for a *bool was accepted")
		}
		if !strings.Contains(err.Error(), "not a boolean") {
			t.Errorf("unhelpful error: %v", err)
		}
	})

	t.Run("optional query parameters honour the same choice", func(t *testing.T) {
		r := base()
		var noAvatar *RecipeArg
		for i := range r.Args {
			if r.Args[i].Name == "noAvatar" {
				noAvatar = &r.Args[i]
			}
		}
		if noAvatar == nil {
			t.Skip("noAvatar is no longer a parameter of this endpoint")
		}
		noAvatar.Unset = UnsetDefault
		noAvatar.Default = "true"

		g, err := Emit(r, op, schema, ctors)
		if err != nil {
			t.Fatal(err)
		}
		src := string(g.Source)
		if strings.Contains(src, `Changed("noAvatar")`) {
			t.Errorf("a defaulted query parameter must not be guarded on Changed:\n%s", src)
		}
		if !strings.Contains(src, "params.NoAvatar = &noAvatar") {
			t.Errorf("query parameter not assigned:\n%s", src)
		}
	})
}

func TestDefaultLiteral(t *testing.T) {
	cases := []struct {
		goType, value, want string
		wantErr             bool
	}{
		{goType: "string", value: "", want: `""`},
		{goType: "string", value: "my-space", want: `"my-space"`},
		{goType: "string", value: `say "hi"`, want: `"say \"hi\""`},
		{goType: "*bool", value: "true", want: "true"},
		{goType: "*bool", value: "false", want: "false"},
		{goType: "*bool", value: "yesplease", wantErr: true},
		{goType: "float32", value: "10", want: "10"},
		{goType: "float32", value: "2.5", want: "2.5"},
		{goType: "float32", value: "lots", wantErr: true},
		{goType: "int", value: "7", want: "7"},
		{goType: "int", value: "7.5", wantErr: true},
		{goType: "[]string", value: "a, b", want: `[]string{"a", "b"}`},
	}

	for _, tc := range cases {
		kind := KindFor(tc.goType, "", false, "client")
		got, err := kind.DefaultLiteral(tc.value)
		switch {
		case tc.wantErr && err == nil:
			t.Errorf("%s default %q was accepted as %s", tc.goType, tc.value, got)
		case !tc.wantErr && err != nil:
			t.Errorf("%s default %q: %v", tc.goType, tc.value, err)
		case !tc.wantErr && got != tc.want:
			t.Errorf("%s default %q = %s, want %s", tc.goType, tc.value, got, tc.want)
		}
	}
}
