package codegen

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// SpecsDir is where the patched specs are vendored, relative to the repo root.
const SpecsDir = "specs"

// OpInfo is the spec metadata for one operation. The generated clients carry no
// doc comments — these specs declare no operationIds and oapi-codegen drops the
// summaries — so this is the only source of human-readable text for a
// generated command's Short and Long, and of the tag that suggests where the
// command belongs in the tree.
type OpInfo struct {
	Summary     string
	Description string
	Tag         string // first tag, e.g. "Organization.Statistics"
	Deprecated  bool
	// Secured records that the operation declares a security requirement.
	// Several endpoints need credentials without declaring an authorization
	// header parameter, so the generated client has no field for it and only
	// the spec can say that a command must send one.
	Secured bool
	// Known distinguishes "the spec says no auth" from "there is no spec".
	Known bool
	// Params holds each parameter's description, keyed by "<in>|<name>". The
	// generated clients carry no doc comment for a path parameter at all, so
	// for those this is the only description that exists.
	Params map[string]string
}

// Param returns the spec's description for one parameter.
func (o OpInfo) Param(in Location, name string) string {
	return o.Params[string(in)+"|"+name]
}

// Specs is the vendored spec metadata, keyed by client package.
type Specs struct {
	ops map[string]OpInfo // "<client> <METHOD> <path>" -> info
}

// LoadSpecs reads specs/*.json. The specs are an optional enrichment: coverage
// and generation both work without them, just with less readable output, so a
// missing directory returns an empty set rather than an error.
func LoadSpecs(root string) *Specs {
	s := &Specs{ops: map[string]OpInfo{}}

	files, err := filepath.Glob(filepath.Join(root, SpecsDir, "*.json"))
	if err != nil {
		return s
	}
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var doc struct {
			Paths map[string]map[string]struct {
				Summary     string           `json:"summary"`
				Description string           `json:"description"`
				Tags        []string         `json:"tags"`
				Deprecated  bool             `json:"deprecated"`
				Security    []map[string]any `json:"security"`
				Parameters  []specParameter  `json:"parameters"`
			} `json:"paths"`
		}
		if err := json.Unmarshal(data, &doc); err != nil {
			continue
		}

		client := ClientPackage(strings.TrimSuffix(filepath.Base(path), ".json"))
		for p, methods := range doc.Paths {
			for method, op := range methods {
				switch strings.ToLower(method) {
				case "get", "put", "post", "delete", "patch", "head", "options", "trace":
				default:
					continue
				}
				// Many descriptions start with a newline. YAML strips that on
				// the way back out of a recipe, so leaving it in would make a
				// regenerated command differ from the one just drafted — and a
				// Long that opens with a blank line reads badly anyway.
				info := OpInfo{
					Summary:     strings.TrimSpace(op.Summary),
					Description: strings.TrimSpace(op.Description),
					Deprecated:  op.Deprecated,
					Secured:     len(op.Security) > 0,
					Known:       true,
				}
				if len(op.Tags) > 0 {
					info.Tag = op.Tags[0]
				}
				for _, param := range op.Parameters {
					if desc := param.describe(); desc != "" {
						if info.Params == nil {
							info.Params = map[string]string{}
						}
						info.Params[param.In+"|"+param.Name] = desc
					}
				}
				s.ops[specKey(client, method, p)] = info
			}
		}
	}
	return s
}

// Empty reports whether any spec metadata was found.
func (s *Specs) Empty() bool { return s == nil || len(s.ops) == 0 }

// Lookup returns the spec metadata for an operation.
func (s *Specs) Lookup(client, method, path string) (OpInfo, bool) {
	if s == nil {
		return OpInfo{}, false
	}
	info, ok := s.ops[specKey(client, method, path)]
	return info, ok
}

// For returns the spec metadata for an operation, ignoring whether it was found.
func (s *Specs) For(op *Operation) OpInfo {
	info, _ := s.Lookup(op.Client, op.Method, op.Path)
	return info
}

// specParameter is one entry of an operation's parameter list.
type specParameter struct {
	Name   string     `json:"name"`
	In     string     `json:"in"`
	Desc   string     `json:"description"`
	Schema specSchema `json:"schema"`
}

// specSchema is the part of a parameter's schema that carries prose.
type specSchema struct {
	Desc  string       `json:"description"`
	AnyOf []specSchema `json:"anyOf"`
	OneOf []specSchema `json:"oneOf"`
}

// describe returns the parameter's description.
//
// These specs put it on the schema rather than on the parameter, and an
// org-name parameter is a union of two string variants that each describe
// themselves. cmd/patcher joins those variants with " | " so nothing is lost
// during generation, but a joined pair runs past 200 characters and reads badly
// as flag help, so only the first variant is offered.
func (p specParameter) describe() string {
	desc := p.Desc
	if desc == "" {
		desc = p.Schema.describe()
	}
	if first, _, joined := strings.Cut(desc, " | "); joined {
		desc = first
	}
	return strings.TrimSpace(desc)
}

func (s specSchema) describe() string {
	if s.Desc != "" {
		return s.Desc
	}
	for _, variants := range [][]specSchema{s.AnyOf, s.OneOf} {
		for _, v := range variants {
			if desc := v.describe(); desc != "" {
				return desc
			}
		}
	}
	return ""
}

func specKey(client, method, path string) string {
	return client + " " + strings.ToUpper(method) + " " + path
}

// ClientPackage turns a service shorthand into its client package name.
// "cosmo" and "cosmoClient" both mean cosmoClient.
func ClientPackage(service string) string {
	if strings.HasSuffix(service, "Client") {
		return service
	}
	return service + "Client"
}

// Service is the inverse of ClientPackage: "mothershipClient" -> "mothership".
// It is also the package name under internal/cli for that service's commands.
func Service(client string) string {
	return strings.TrimSuffix(client, "Client")
}
