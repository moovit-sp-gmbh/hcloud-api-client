package codegen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSpecPrefillsDescriptions covers the reason specs/ is vendored at all: the
// generated clients carry no doc comment for a path parameter and no operation
// summary, so without the spec every description a command shows would have to
// be typed by hand.
func TestSpecPrefillsDescriptions(t *testing.T) {
	if _, err := os.Stat("../../specs"); err != nil {
		t.Skip("specs/ not vendored; run make specs")
	}
	ix, schema, specs, _ := loadFixtures(t)
	if specs.Empty() {
		t.Fatal("specs/ exists but nothing was parsed out of it")
	}

	op, ok := ix.Op("high5Client", "PostV1OrgOrgNameSpaces")
	if !ok {
		t.Fatal("PostV1OrgOrgNameSpaces not in the high5 client")
	}
	r := Derive(op, ix.Clients["high5Client"], specs, schema)

	if r.Command.Short != "Create a new space" {
		t.Errorf("Short = %q, want the spec's summary", r.Command.Short)
	}
	if !strings.HasPrefix(r.Command.Long, "Creates a new space in the organization") {
		t.Errorf("Long = %q, want the spec's description", r.Command.Long)
	}

	// The path parameter is the case that has no other source at all.
	var orgName *RecipeArg
	for i := range r.Args {
		if r.Args[i].Name == "orgName" {
			orgName = &r.Args[i]
		}
	}
	if orgName == nil {
		t.Fatal("no orgName argument")
	}
	if !strings.HasPrefix(orgName.Description, "Name of the organization") {
		t.Errorf("orgName description = %q, want the spec's", orgName.Description)
	}
	// cmd/patcher joins the two org-name union variants with " | "; only the
	// first belongs on a flag help line.
	if strings.Contains(orgName.Description, " | ") {
		t.Errorf("the joined union variants leaked into the flag help: %q", orgName.Description)
	}

	// Body field descriptions come from the client's own doc comments, which
	// oapi-codegen does emit, and must survive too.
	for _, f := range r.Body.Fields {
		if f.Field == "EnableDev" && !strings.Contains(f.Description, "versions under development") {
			t.Errorf("enableDev description = %q", f.Description)
		}
	}
}

// TestSpecDescriptionPrecedence pins the extraction rules against a spec written
// the way these APIs write them: prose on the schema, not on the parameter, and
// union variants that each describe themselves.
func TestSpecDescriptionPrecedence(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, SpecsDir), 0o755); err != nil {
		t.Fatal(err)
	}
	spec := `{
	  "paths": {
	    "/v1/org/{orgName}/things/{thingId}": {
	      "get": {
	        "summary": "List things",
	        "description": "Lists the things in an organization",
	        "tags": ["Organization.Things"],
	        "security": [{"bearerAuth": []}],
	        "parameters": [
	          {"name": "orgName", "in": "path",
	           "schema": {"anyOf": [
	             {"type": "string", "description": "Name of the organization"},
	             {"type": "string", "description": "Name of the private organization"}
	           ]}},
	          {"name": "thingId", "in": "path", "description": "on the parameter",
	           "schema": {"type": "string", "description": "on the schema"}},
	          {"name": "limit", "in": "query",
	           "schema": {"type": "number", "description": "Page size"}},
	          {"name": "authorization", "in": "header", "schema": {"type": "string"}}
	        ]
	      }
	    }
	  }
	}`
	if err := os.WriteFile(filepath.Join(dir, SpecsDir, "widget.json"), []byte(spec), 0o644); err != nil {
		t.Fatal(err)
	}

	specs := LoadSpecs(dir)
	info, ok := specs.Lookup("widgetClient", "GET", "/v1/org/{orgName}/things/{thingId}")
	if !ok {
		t.Fatal("the operation was not indexed")
	}

	if info.Summary != "List things" || info.Tag != "Organization.Things" || !info.Secured {
		t.Errorf("operation info = %+v", info)
	}
	for _, tc := range []struct {
		in   Location
		name string
		want string
		why  string
	}{
		{InPath, "orgName", "Name of the organization", "the first union variant, not the joined pair"},
		{InPath, "thingId", "on the parameter", "the parameter's own description wins over the schema's"},
		{InQuery, "limit", "Page size", "a plain schema description"},
		{InHeader, "authorization", "", "no description in the spec"},
	} {
		if got := info.Param(tc.in, tc.name); got != tc.want {
			t.Errorf("%s %s = %q, want %q (%s)", tc.in, tc.name, got, tc.want, tc.why)
		}
	}
}

// TestLoadSpecsToleratesAbsence keeps generation working before anyone has run
// make specs.
func TestLoadSpecsToleratesAbsence(t *testing.T) {
	specs := LoadSpecs(t.TempDir())
	if !specs.Empty() {
		t.Error("an empty tree produced spec metadata")
	}
	if got := specs.For(&Operation{Client: "x", Method: "GET", Path: "/v1/y"}); got.Summary != "" || got.Known {
		t.Errorf("lookup against no specs returned %+v", got)
	}
	if got := (OpInfo{}).Param(InPath, "orgName"); got != "" {
		t.Errorf("param lookup on a zero OpInfo returned %q", got)
	}
}
