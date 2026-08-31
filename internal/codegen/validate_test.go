package codegen

import (
	"strings"
	"testing"
)

// validRecipe is the shape Derive produces for GetV1OrgOrgNameStats, used as the
// baseline each case breaks in one way.
func validRecipe() *Recipe {
	return &Recipe{
		Client:    "mothershipClient",
		Operation: "GetV1OrgOrgNameStats",
		Command: RecipeCommand{
			Path:  []string{"mothership", "stats"},
			Var:   "statsCommand",
			File:  "stats.go",
			Short: "Retrieve usage statistics",
		},
		Args: []RecipeArg{
			{
				Source: InPath, Name: "orgName", CLI: "orgName", Var: "orgName", Type: "string",
				Resolve: []string{FromFlag, FromContext, FromPrompt}, Context: "OrgName",
				Persist: true, Prompt: "orgName",
			},
			{
				Source: InHeader, Name: "authorization", CLI: "token", Var: "jwt", Type: "string",
				Resolve: []string{FromFlag, FromContext, FromPrompt}, Context: "Jwt",
				Persist: true, Prompt: "authorization token", Via: ViaParams,
			},
		},
		Body:     RecipeBody{Mode: BodyNone},
		Response: RecipeResponse{Render: RenderBody},
	}
}

func TestValidateAcceptsDerivedRecipe(t *testing.T) {
	if err := validRecipe().Validate(); err != nil {
		t.Fatalf("a well-formed recipe was rejected: %v", err)
	}

	// And every recipe Derive produces for a callable operation must pass, or
	// regen would reject its own output.
	ix, schema, specs, _ := loadFixtures(t)
	for _, name := range ix.ClientNames() {
		client := ix.Clients[name]
		for _, op := range client.Sorted() {
			if !op.Callable() {
				continue
			}
			r := Derive(op, client, specs, schema)
			r.Command.Var = CommandVar(r.Command.Path)
			r.Command.File = r.Command.Use() + ".go"
			if err := r.Validate(); err != nil {
				t.Errorf("Derive produced an invalid recipe for %s.%s: %v", name, op.Name, err)
			}
		}
	}
}

func TestValidateRejects(t *testing.T) {
	cases := map[string]struct {
		mutate func(*Recipe)
		want   string
	}{
		// The failure that motivated this: a free-text answer landing in the
		// context field would otherwise become a migration column.
		"prose context field": {
			func(r *Recipe) { r.Args[0].Context = "Name of the space" },
			"not an identifier",
		},
		"command path with spaces in a segment": {
			func(r *Recipe) { r.Command.Path = []string{"cosmo", "address book"} },
			"lowercase letters, digits and dashes",
		},
		"uppercase command segment": {
			func(r *Recipe) { r.Command.Path = []string{"cosmo", "Space"} },
			"lowercase",
		},
		"single-segment path": {
			func(r *Recipe) { r.Command.Path = []string{"mothership"} },
			"at least a service",
		},
		"non-identifier command variable": {
			func(r *Recipe) { r.Command.Var = "stats-command" },
			"not an identifier",
		},
		"file with a directory": {
			func(r *Recipe) { r.Command.File = "sub/stats.go" },
			"bare .go file name",
		},
		"duplicate flag": {
			func(r *Recipe) { r.Args[1].CLI = "orgName" },
			"both use the flag --orgName",
		},
		"multi-character shorthand": {
			func(r *Recipe) { r.Args[0].Shorthand = "org" },
			"single character",
		},
		"context resolution without a field": {
			func(r *Recipe) { r.Args[0].Context = "" },
			"names no context field",
		},
		"unknown resolution source": {
			func(r *Recipe) { r.Args[0].Resolve = []string{"guess"} },
			"unknown resolution source",
		},
		"unknown body mode": {
			func(r *Recipe) { r.Body.Mode = "inline" },
			"unknown body mode",
		},
		"typed render without a field": {
			func(r *Recipe) { r.Response.Render = RenderTyped },
			"names no response field",
		},
		"malformed persist source": {
			func(r *Recipe) {
				r.Response.Persist = []RecipePersist{{From: "JSON200", Context: "AgentID"}}
			},
			"header:<name> or",
		},
		"unknown via": {
			func(r *Recipe) { r.Args[1].Via = "middleware" },
			"must be params or editor",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			r := validRecipe()
			tc.mutate(r)
			err := r.Validate()
			if err == nil {
				t.Fatalf("%s was accepted", name)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("error %q does not mention %q", err, tc.want)
			}
		})
	}
}

// TestValidateIgnoresSkipped makes sure an argument left out is not held to the
// same rules as one that is used.
func TestValidateIgnoresSkipped(t *testing.T) {
	r := validRecipe()
	r.Args = append(r.Args, RecipeArg{
		Source: InQuery, Name: "namespace", CLI: "", Skip: true,
	})
	if err := r.Validate(); err != nil {
		t.Errorf("a skipped argument was validated: %v", err)
	}
}
