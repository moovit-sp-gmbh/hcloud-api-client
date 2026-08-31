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
