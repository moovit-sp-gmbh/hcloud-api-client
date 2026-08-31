package codegen

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Resolution sources for an argument, in the order the generated GetArg chain
// tries them.
const (
	FromFlag    = "flag"
	FromContext = "context"
	FromPrompt  = "prompt"
)

// Body modes.
const (
	BodyNone  = "none"
	BodyTyped = "typed"
	BodyRaw   = "raw"
)

// How a header parameter reaches the request.
const (
	ViaParams = "params" // the operation's typed Params field
	ViaEditor = "editor" // a RequestEditorFn closure
)

// RecipeCommand says where in the tree the command goes and what it is called.
type RecipeCommand struct {
	Path  []string `yaml:"path"` // full path under the root command
	File  string   `yaml:"file"` // file under the service package
	Var   string   `yaml:"var"`  // the cobra.Command variable name
	Short string   `yaml:"short,omitempty"`
	Long  string   `yaml:"long,omitempty"`
}

// Use is the command's own name, the last segment of its path.
func (c RecipeCommand) Use() string {
	if len(c.Path) == 0 {
		return ""
	}
	return c.Path[len(c.Path)-1]
}

// Service is the service package the command lives in, the first path segment.
func (c RecipeCommand) Service() string {
	if len(c.Path) == 0 {
		return ""
	}
	return c.Path[0]
}

// RecipeArg maps one API parameter or body field onto a CLI argument.
type RecipeArg struct {
	Source      Location `yaml:"source"` // path, query or header
	Field       string   `yaml:"field,omitempty"`
	Name        string   `yaml:"name"`             // API-side name
	GoName      string   `yaml:"goName,omitempty"` // Params field or positional arg
	CLI         string   `yaml:"cli"`              // flag name
	Shorthand   string   `yaml:"shorthand,omitempty"`
	Var         string   `yaml:"var,omitempty"` // local variable in the generated code
	Type        string   `yaml:"type"`          // Go type as the client declares it
	Description string   `yaml:"description,omitempty"`
	Default     string   `yaml:"default,omitempty"`
	Resolve     []string `yaml:"resolve"`           // flag, context, prompt
	Context     string   `yaml:"context,omitempty"` // db.Context field to read
	Persist     bool     `yaml:"persist,omitempty"` // write the resolved value back
	Prompt      string   `yaml:"prompt,omitempty"`
	Via         string   `yaml:"via,omitempty"` // header only
	Skip        bool     `yaml:"skip,omitempty"`
}

// RecipeBodyField maps one request-body field onto a flag.
type RecipeBodyField struct {
	Field       string   `yaml:"field"` // Go field on the DTO
	Name        string   `yaml:"name"`  // JSON name
	CLI         string   `yaml:"cli"`
	Shorthand   string   `yaml:"shorthand,omitempty"`
	Var         string   `yaml:"var,omitempty"`
	Type        string   `yaml:"type"`
	Description string   `yaml:"description,omitempty"`
	Optional    bool     `yaml:"optional,omitempty"`
	Resolve     []string `yaml:"resolve"`
	Context     string   `yaml:"context,omitempty"`
	Persist     bool     `yaml:"persist,omitempty"`
	Prompt      string   `yaml:"prompt,omitempty"`
	Skip        bool     `yaml:"skip,omitempty"`
}

// RecipeBody says how the request body is built.
type RecipeBody struct {
	Mode   string            `yaml:"mode"` // none, typed or raw
	Flag   string            `yaml:"flag,omitempty"`
	Fields []RecipeBodyField `yaml:"fields,omitempty"`
}

// RecipePersist saves one piece of a successful response into the context table.
type RecipePersist struct {
	// From names the source: "header:<name>" for a response header, or
	// "<responseField>:<DtoField>" such as "JSON200:Id" for a body field.
	From    string `yaml:"from"`
	Context string `yaml:"context"`          // db.Context field name
	Column  string `yaml:"column,omitempty"` // SQL column; set when it is new
}

// Header reports whether the value comes from a response header, and its name.
func (p RecipePersist) Header() (string, bool) {
	name, ok := strings.CutPrefix(p.From, "header:")
	return name, ok
}

// Body returns the response field and DTO field a body persist reads.
func (p RecipePersist) Body() (respField, dtoField string, ok bool) {
	respField, dtoField, ok = strings.Cut(p.From, ":")
	if !ok || respField == "header" {
		return "", "", false
	}
	return respField, dtoField, true
}

// RecipeResponse says how the response is rendered and what is kept.
type RecipeResponse struct {
	// Render is "typed" to print the typed success field, or "body" to
	// unmarshal and pretty-print res.Body whatever the status.
	Render  string          `yaml:"render"`
	Field   string          `yaml:"field,omitempty"` // e.g. "JSON200"
	Persist []RecipePersist `yaml:"persist,omitempty"`
}

// Render modes.
const (
	RenderTyped = "typed"
	RenderBody  = "body"
)

// Recipe is the complete input needed to generate one command. It is written
// next to the generated code so a template change or a client regen can be
// replayed instead of re-answered.
type Recipe struct {
	Client    string         `yaml:"client"`
	Operation string         `yaml:"operation"`
	Command   RecipeCommand  `yaml:"command"`
	Args      []RecipeArg    `yaml:"args,omitempty"`
	Body      RecipeBody     `yaml:"body"`
	Response  RecipeResponse `yaml:"response"`
}

// RecipesDir is where recipes live, relative to the repo root.
const RecipesDir = "codegen/recipes"

// RecipePath is the file a recipe is stored in, derived from its command path so
// two commands cannot collide.
func (r *Recipe) RecipePath() string {
	name := strings.Join(r.Command.Path, "-")
	if name == "" {
		name = strings.ToLower(r.Operation)
	}
	return filepath.Join(RecipesDir, r.Command.Service(), name+".yaml")
}

// Save writes the recipe.
func (r *Recipe) Save(root string) (string, error) {
	rel := r.RecipePath()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}

	var buf strings.Builder
	buf.WriteString("# Generation input for `" + strings.Join(r.Command.Path, " ") + "`.\n")
	buf.WriteString("# Edit and re-run `codegen regen` to regenerate the command.\n")
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(r); err != nil {
		return "", err
	}
	if err := enc.Close(); err != nil {
		return "", err
	}
	return rel, os.WriteFile(path, []byte(buf.String()), 0o644)
}

// LoadRecipe reads one recipe file.
func LoadRecipe(path string) (*Recipe, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var r Recipe
	if err := yaml.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if r.Client == "" || r.Operation == "" {
		return nil, fmt.Errorf("%s: recipe names no client operation", path)
	}
	return &r, nil
}

// LoadRecipes reads every recipe under codegen/recipes.
func LoadRecipes(root string) ([]*Recipe, []string, error) {
	pattern := filepath.Join(root, RecipesDir, "*", "*.yaml")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, nil, err
	}
	sort.Strings(files)

	var out []*Recipe
	var paths []string
	for _, path := range files {
		r, err := LoadRecipe(path)
		if err != nil {
			return nil, nil, err
		}
		rel, _ := filepath.Rel(root, path)
		out = append(out, r)
		paths = append(paths, rel)
	}
	return out, paths, nil
}
