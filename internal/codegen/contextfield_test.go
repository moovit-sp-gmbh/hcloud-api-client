package codegen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadContextSchema(t *testing.T) {
	s, err := LoadContextSchema("../..")
	if err != nil {
		t.Fatal(err)
	}

	// The first column is spelled JWT in SQL and Jwt in Go; the pairing has to
	// survive that.
	want := map[string]string{
		"Jwt":             "JWT",
		"OrgName":         "org_name",
		"AgentUuid":       "agent_uuid",
		"AgentPrivateKey": "agent_private_key",
	}
	for goName, column := range want {
		f, ok := s.Field(goName)
		if !ok {
			t.Errorf("field %s not found; have %v", goName, s.Names())
			continue
		}
		if f.Column != column {
			t.Errorf("%s column = %q, want %q", goName, f.Column, column)
		}
	}
	if _, ok := s.Field("agent_hello_token"); !ok {
		t.Error("lookup by column name failed for agent_hello_token")
	}
	if _, ok := s.Field("NoSuchField"); ok {
		t.Error("unknown field resolved")
	}
}

func TestNameConversions(t *testing.T) {
	// GoFieldName has to agree with what sqlc actually generated, which is why
	// the model says AgentUuid and not AgentUUID.
	for column, goName := range map[string]string{
		"jwt":               "Jwt",
		"JWT":               "Jwt",
		"org_name":          "OrgName",
		"agent_uuid":        "AgentUuid",
		"agent_private_key": "AgentPrivateKey",
		"agent_hello_token": "AgentHelloToken",
		// sqlc upper-cases "id" wholesale but not "uuid", so a column named
		// agent_id becomes AgentID while agent_uuid becomes AgentUuid.
		"agent_id":   "AgentID",
		"id":         "ID",
		"pool_id":    "PoolID",
		"identifier": "Identifier",
	} {
		if got := GoFieldName(column); got != goName {
			t.Errorf("GoFieldName(%q) = %q, want %q", column, got, goName)
		}
	}

	for name, column := range map[string]string{
		"agentPoolId":    "agent_pool_id",
		"AgentPoolId":    "agent_pool_id",
		"spaceName":      "space_name",
		"jwt":            "jwt",
		"patch-password": "patch_password",
		"AgentUuid":      "agent_uuid",
	} {
		if got := ColumnName(name); got != column {
			t.Errorf("ColumnName(%q) = %q, want %q", name, got, column)
		}
	}
}

// TestAddCreateContextColumns edits a copy of the real query file: the column
// list and the placeholder list must grow together or sqlc generates an insert
// that fails at runtime.
func TestAddCreateContextColumns(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "internal", "db")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	original, err := os.ReadFile("../db/context.sql")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "context.sql")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatal(err)
	}

	// agent_uuid already exists, so only the two new columns should land.
	if err := AddCreateContextColumns(root, []string{"agent_pool_id", "agent_uuid", "space_name"}); err != nil {
		t.Fatal(err)
	}

	edited, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got := string(edited)

	columns, err := createContextColumns(path)
	if err != nil {
		t.Fatal(err)
	}
	before, _ := createContextColumnsFrom(string(original))
	if len(columns) != len(before)+2 {
		t.Fatalf("columns = %d, want %d", len(columns), len(before)+2)
	}
	for _, want := range []string{"agent_pool_id", "space_name"} {
		found := false
		for _, c := range columns {
			if c == want {
				found = true
			}
		}
		if !found {
			t.Errorf("column %q missing after edit; got %v", want, columns)
		}
	}
	if n := strings.Count(columns[0], "agent_uuid"); n != 0 {
		t.Error("unexpected duplicate")
	}

	// Count placeholders in the VALUES list and require a one-to-one match.
	valStart, valEnd, err := valuesListRange(got, strings.Index(got, ") VALUES"))
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(got[valStart:valEnd], "?"); n != len(columns) {
		t.Errorf("%d placeholders for %d columns\n%s", n, len(columns), got[valStart:valEnd])
	}

	// A second run must be a no-op.
	if err := AddCreateContextColumns(root, []string{"agent_pool_id"}); err != nil {
		t.Fatal(err)
	}
	again, _ := os.ReadFile(path)
	if string(again) != got {
		t.Error("re-adding an existing column changed the file")
	}
}

// createContextColumnsFrom is a test helper mirroring createContextColumns for
// in-memory text.
func createContextColumnsFrom(text string) ([]string, error) {
	match := createContextRe.FindStringSubmatch(text)
	if match == nil {
		return nil, os.ErrNotExist
	}
	var out []string
	for _, raw := range strings.Split(match[1], ",") {
		if col := strings.TrimSpace(raw); col != "" {
			out = append(out, col)
		}
	}
	return out, nil
}

func TestMigrationFor(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "internal", "db", "migrations"), 0o755); err != nil {
		t.Fatal(err)
	}

	stamp := time.Date(2026, 8, 31, 12, 30, 45, 0, time.UTC)
	path, err := MigrationFor(root, []string{"agent_pool_id", "space_name"}, "agent_pool", stamp)
	if err != nil {
		t.Fatal(err)
	}
	if base := filepath.Base(path); base != "20260831123045_agent_pool_ctx.sql" {
		t.Errorf("filename = %q", base)
	}

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := `-- +goose Up
ALTER TABLE context
ADD COLUMN agent_pool_id TEXT;
ALTER TABLE context
ADD COLUMN space_name TEXT;

-- +goose Down
ALTER TABLE context
DROP COLUMN agent_pool_id;
ALTER TABLE context
DROP COLUMN space_name;
`
	if string(body) != want {
		t.Errorf("migration body:\n%s\nwant:\n%s", body, want)
	}
}
