package codegen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// ContextField is one column of the context table — the CLI's whole notion of
// session state.
type ContextField struct {
	GoName string // field name on db.Context, e.g. "AgentUuid"
	Column string // SQL column name, e.g. "agent_uuid"
}

// ContextSchema is the context table as it stands, paired up between the sqlc
// model and the insert statement.
type ContextSchema struct {
	Fields []ContextField
}

// Field finds a field by Go name or column name, whichever the caller has.
func (s *ContextSchema) Field(name string) (ContextField, bool) {
	key := columnKey(name)
	for _, f := range s.Fields {
		if columnKey(f.GoName) == key || columnKey(f.Column) == key {
			return f, true
		}
	}
	return ContextField{}, false
}

// Names returns the Go field names, for prompting.
func (s *ContextSchema) Names() []string {
	out := make([]string, 0, len(s.Fields))
	for _, f := range s.Fields {
		out = append(out, f.GoName)
	}
	return out
}

// LoadContextSchema reads the context table's columns.
//
// The Go side comes from the sqlc model and the SQL side from the CreateContext
// insert, matched on a normalized name rather than on position: the first
// column is spelled JWT in SQL and Jwt in Go, so nothing may assume the two
// spellings agree.
func LoadContextSchema(root string) (*ContextSchema, error) {
	goNames, err := contextModelFields(filepath.Join(root, "internal", "db", "models.go"))
	if err != nil {
		return nil, err
	}
	columns, err := createContextColumns(filepath.Join(root, "internal", "db", "context.sql"))
	if err != nil {
		return nil, err
	}

	byKey := map[string]string{}
	for _, c := range columns {
		byKey[columnKey(c)] = c
	}

	s := &ContextSchema{}
	for _, goName := range goNames {
		column, ok := byKey[columnKey(goName)]
		if !ok {
			// A model field with no column in the insert is one CreateContext
			// does not write — id and timestamp are handled separately, so this
			// is worth surfacing rather than silently dropping.
			continue
		}
		s.Fields = append(s.Fields, ContextField{GoName: goName, Column: column})
	}
	return s, nil
}

// contextModelFields returns the exported fields of db.Context, minus the ones
// the service manages itself.
func contextModelFields(path string) ([]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}

	var out []string
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok || ts.Name.Name != "Context" {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}
			for _, field := range st.Fields.List {
				for _, n := range field.Names {
					if n.Name == "ID" || n.Name == "Timestamp" || !n.IsExported() {
						continue
					}
					out = append(out, n.Name)
				}
			}
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("%s: no db.Context struct fields found", path)
	}
	return out, nil
}

var createContextRe = regexp.MustCompile(`(?s)--\s*name:\s*CreateContext\b.*?INSERT INTO context\s*\((.*?)\)\s*VALUES`)

// createContextColumns returns the columns CreateContext inserts, in order.
func createContextColumns(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	match := createContextRe.FindSubmatch(data)
	if match == nil {
		return nil, fmt.Errorf("%s: no CreateContext insert found", path)
	}

	var out []string
	for _, raw := range strings.Split(string(match[1]), ",") {
		if col := strings.TrimSpace(raw); col != "" {
			out = append(out, col)
		}
	}
	return out, nil
}

// columnKey normalizes a column or field name so the two spellings of one
// column compare equal: JWT, jwt and Jwt all key to "jwt".
func columnKey(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, "_", ""))
}

// sqlcInitialisms are the name parts sqlc upper-cases wholesale. Its default
// list holds only "id", which is why the generated model says AgentUuid and Jwt
// but ID — anything else is title-cased.
var sqlcInitialisms = map[string]string{"id": "ID"}

// GoFieldName turns a snake_case column into the field name sqlc will generate
// for it. Getting this wrong does not fail quietly: the generated command
// references a field that does not exist and the tree stops compiling.
func GoFieldName(column string) string {
	var sb strings.Builder
	for _, part := range strings.Split(column, "_") {
		if part == "" {
			continue
		}
		lower := strings.ToLower(part)
		if initialism, ok := sqlcInitialisms[lower]; ok {
			sb.WriteString(initialism)
			continue
		}
		sb.WriteString(strings.ToUpper(part[:1]))
		sb.WriteString(lower[1:])
	}
	return sb.String()
}

var camelBoundary = regexp.MustCompile(`([a-z0-9])([A-Z])`)

// ColumnName turns a Go-ish or camelCase name into a snake_case column, e.g.
// "agentPoolId" and "AgentPoolId" both become "agent_pool_id".
func ColumnName(name string) string {
	s := camelBoundary.ReplaceAllString(name, "${1}_${2}")
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ToLower(s)
	for strings.Contains(s, "__") {
		s = strings.ReplaceAll(s, "__", "_")
	}
	return strings.Trim(s, "_")
}

// MigrationContent renders a goose migration adding columns to the context
// table, in the style of internal/db/migrations/20260812091305_agent_keys_ctx.sql.
func MigrationContent(columns []string) string {
	var up, down strings.Builder
	for _, c := range columns {
		fmt.Fprintf(&up, "ALTER TABLE context\nADD COLUMN %s TEXT;\n", c)
		fmt.Fprintf(&down, "ALTER TABLE context\nDROP COLUMN %s;\n", c)
	}
	return fmt.Sprintf("-- +goose Up\n%s\n-- +goose Down\n%s", up.String(), down.String())
}

// MigrationFilename is the goose-ordered name for a context migration.
func MigrationFilename(name string, now time.Time) string {
	return fmt.Sprintf("%s_%s_ctx.sql", now.Format("20060102150405"), name)
}

// MigrationFor writes a goose migration adding columns to the context table and
// returns the path it wrote.
func MigrationFor(root string, columns []string, name string, now time.Time) (string, error) {
	if len(columns) == 0 {
		return "", fmt.Errorf("no columns to add")
	}
	if name == "" {
		name = columns[0]
	}
	path := filepath.Join(root, "internal", "db", "migrations", MigrationFilename(name, now))
	if err := os.WriteFile(path, []byte(MigrationContent(columns)), 0o644); err != nil {
		return "", err
	}
	return path, nil
}

// AddCreateContextColumns adds columns to the CreateContext insert in
// internal/db/context.sql, so sqlc regenerates a CreateContextParams that can
// carry them.
func AddCreateContextColumns(root string, columns []string) error {
	path := filepath.Join(root, "internal", "db", "context.sql")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	text := string(data)

	loc := createContextRe.FindStringSubmatchIndex(text)
	if loc == nil {
		return fmt.Errorf("%s: no CreateContext insert found", path)
	}
	colStart, colEnd := loc[2], loc[3]

	existing, err := createContextColumns(path)
	if err != nil {
		return err
	}
	have := map[string]bool{}
	for _, c := range existing {
		have[columnKey(c)] = true
	}

	var toAdd []string
	for _, c := range columns {
		if have[columnKey(c)] {
			continue
		}
		have[columnKey(c)] = true
		toAdd = append(toAdd, c)
	}
	if len(toAdd) == 0 {
		return nil
	}

	// The placeholder list has to grow by exactly as much as the column list,
	// so locate it before touching either.
	valStart, valEnd, err := valuesListRange(text, colEnd)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	indent := insertIndent(text[colStart:colEnd])
	columnList := strings.TrimRight(text[colStart:colEnd], " \t\n")
	placeholders := strings.TrimRight(text[valStart:valEnd], " \t\n")
	for _, c := range toAdd {
		columnList += ",\n" + indent + c
		placeholders += ",\n" + indent + "?"
	}
	columnList += "\n"
	placeholders += "\n"

	// Splice the later range first; rewriting the column list would invalidate
	// the placeholder offsets.
	out := text[:valStart] + placeholders + text[valEnd:]
	out = out[:colStart] + columnList + out[colEnd:]
	return os.WriteFile(path, []byte(out), 0o644)
}

// valuesListRange finds the span between "VALUES (" and its closing paren,
// starting the search after the column list.
func valuesListRange(text string, after int) (start, end int, err error) {
	rest := text[after:]
	kw := strings.Index(rest, "VALUES")
	if kw < 0 {
		return 0, 0, fmt.Errorf("no VALUES keyword after the column list")
	}
	open := strings.Index(rest[kw:], "(")
	if open < 0 {
		return 0, 0, fmt.Errorf("no opening paren after VALUES")
	}
	open += kw + 1

	depth := 1
	for i := open; i < len(rest); i++ {
		switch rest[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return after + open, after + i, nil
			}
		}
	}
	return 0, 0, fmt.Errorf("unterminated VALUES list")
}

// insertIndent recovers the indentation the existing column list uses, so an
// added column lines up with the ones already there.
func insertIndent(list string) string {
	for _, line := range strings.Split(list, "\n") {
		if trimmed := strings.TrimLeft(line, " \t"); trimmed != "" {
			return line[:len(line)-len(trimmed)]
		}
	}
	return "        "
}
