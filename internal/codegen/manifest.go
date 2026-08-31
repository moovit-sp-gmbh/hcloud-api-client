package codegen

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Status values for a manifest entry.
const (
	StatusImplemented = "implemented"
	StatusSkipped     = "skipped"
)

// Source values for a manifest entry.
const (
	SourceHandwritten = "handwritten"
	SourceGenerated   = "generated"
)

// Entry records one endpoint the CLI has dealt with — either because a command
// covers it, or because it was deliberately left alone.
type Entry struct {
	Client    string `yaml:"client"`
	Operation string `yaml:"operation"`
	Method    string `yaml:"method,omitempty"`
	Path      string `yaml:"path,omitempty"`
	Command   string `yaml:"command,omitempty"`
	Status    string `yaml:"status,omitempty"` // defaults to implemented
	Source    string `yaml:"source,omitempty"`
	Recipe    string `yaml:"recipe,omitempty"`
	Reason    string `yaml:"reason,omitempty"` // why it is skipped
}

// Key matches Operation.Key.
func (e Entry) Key() string { return e.Client + "." + e.Operation }

// Skipped reports whether the endpoint was deliberately left unimplemented.
func (e Entry) Skipped() bool { return e.Status == StatusSkipped }

// Manifest is the committed list of endpoints the CLI has handled. It is the
// authoritative answer to "what is covered": the code is scanned only to check
// the manifest against reality, never instead of reading it.
type Manifest struct {
	Endpoints []Entry `yaml:"endpoints"`
}

// ManifestPath is where the manifest lives, relative to the repo root.
const ManifestPath = "codegen/manifest.yaml"

// LoadManifest reads the manifest. A missing file is not an error — it is an
// empty manifest, which is what a first run sees.
func LoadManifest(root string) (*Manifest, error) {
	path := filepath.Join(root, ManifestPath)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Manifest{}, nil
	}
	if err != nil {
		return nil, err
	}

	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	for i := range m.Endpoints {
		if m.Endpoints[i].Status == "" {
			m.Endpoints[i].Status = StatusImplemented
		}
	}
	return &m, nil
}

// Save writes the manifest back, sorted by client then path then method so the
// file stays reviewable and a regenerated entry does not reshuffle the diff.
func (m *Manifest) Save(root string) error {
	m.sort()

	path := filepath.Join(root, ManifestPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	var buf strings.Builder
	buf.WriteString("# Endpoints this CLI has handled, and the ones deliberately left alone.\n")
	buf.WriteString("# Written by cmd/codegen; `make coverage` reads it. Hand edits are fine —\n")
	buf.WriteString("# add a `status: skipped` entry with a `reason` to retire an endpoint from\n")
	buf.WriteString("# the missing list.\n")

	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(m); err != nil {
		return err
	}
	if err := enc.Close(); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(buf.String()), 0o644)
}

func (m *Manifest) sort() {
	sort.SliceStable(m.Endpoints, func(i, j int) bool {
		a, b := m.Endpoints[i], m.Endpoints[j]
		if a.Client != b.Client {
			return a.Client < b.Client
		}
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		if a.Method != b.Method {
			return a.Method < b.Method
		}
		return a.Operation < b.Operation
	})
}

// Lookup finds the entry for a client operation.
func (m *Manifest) Lookup(client, op string) (Entry, bool) {
	for _, e := range m.Endpoints {
		if e.Client == client && e.Operation == op {
			return e, true
		}
	}
	return Entry{}, false
}

// Upsert adds an entry, replacing any existing one for the same operation.
func (m *Manifest) Upsert(e Entry) {
	if e.Status == "" {
		e.Status = StatusImplemented
	}
	for i, existing := range m.Endpoints {
		if existing.Client == e.Client && existing.Operation == e.Operation {
			m.Endpoints[i] = e
			return
		}
	}
	m.Endpoints = append(m.Endpoints, e)
}
