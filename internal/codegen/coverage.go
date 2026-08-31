package codegen

import (
	"fmt"
	"sort"
	"strings"
)

// ClientCoverage is the coverage of one client package.
type ClientCoverage struct {
	Client string `json:"client"`
	Total  int    `json:"total"`

	Implemented []Entry      `json:"implemented,omitempty"`
	Skipped     []Entry      `json:"skipped,omitempty"`
	Missing     []*Operation `json:"missing,omitempty"`
	// Stale entries name an operation the client no longer exposes. Clients are
	// regenerated from live specs, so an endpoint can disappear underneath a
	// command; that is a build break waiting to happen, not a coverage gap.
	Stale []Entry `json:"stale,omitempty"`
}

// Coverage is the comparison of the manifest against the generated clients.
type Coverage struct {
	Clients []ClientCoverage `json:"clients"`
}

// Totals sums the per-client counts.
func (c *Coverage) Totals() (total, implemented, skipped, missing, stale int) {
	for _, cc := range c.Clients {
		total += cc.Total
		implemented += len(cc.Implemented)
		skipped += len(cc.Skipped)
		missing += len(cc.Missing)
		stale += len(cc.Stale)
	}
	return
}

// Compare joins the manifest against the operation index. When clients is
// non-empty only those client packages are reported.
func Compare(ix *Index, m *Manifest, clients []string) *Coverage {
	want := map[string]bool{}
	for _, c := range clients {
		want[ClientPackage(c)] = true
	}

	entries := map[string][]Entry{}
	for _, e := range m.Endpoints {
		entries[e.Client] = append(entries[e.Client], e)
	}

	cov := &Coverage{}
	for _, name := range ix.ClientNames() {
		if len(want) > 0 && !want[name] {
			continue
		}
		client := ix.Clients[name]
		cc := ClientCoverage{Client: name, Total: len(client.Ops)}

		covered := map[string]bool{}
		for _, e := range entries[name] {
			if _, ok := client.Ops[e.Operation]; !ok {
				cc.Stale = append(cc.Stale, e)
				continue
			}
			covered[e.Operation] = true
			if e.Skipped() {
				cc.Skipped = append(cc.Skipped, e)
			} else {
				cc.Implemented = append(cc.Implemented, e)
			}
		}

		for _, op := range client.Sorted() {
			if !covered[op.Name] {
				cc.Missing = append(cc.Missing, op)
			}
		}
		cov.Clients = append(cov.Clients, cc)
	}
	return cov
}

// Discrepancy is one disagreement between the manifest and the source.
type Discrepancy struct {
	Kind   string `json:"kind"`
	Client string `json:"client"`
	Op     string `json:"op"`
	Detail string `json:"detail"`
}

// Verify cross-checks the manifest against the operations the command tree
// actually calls. The manifest stays authoritative; this reports where it has
// drifted so a stale row cannot quietly claim coverage that no longer exists.
func Verify(ix *Index, m *Manifest, tree *Tree, sites []CallSite) []Discrepancy {
	var out []Discrepancy

	called := map[string][]CallSite{}
	for _, s := range sites {
		called[s.Key()] = append(called[s.Key()], s)
	}

	for _, e := range m.Endpoints {
		if e.Skipped() {
			if sites, ok := called[e.Key()]; ok {
				out = append(out, Discrepancy{
					Kind:   "skipped-but-called",
					Client: e.Client, Op: e.Operation,
					Detail: fmt.Sprintf("marked skipped but called at %s:%d", sites[0].File, sites[0].Line),
				})
			}
			continue
		}

		sites, ok := called[e.Key()]
		if !ok {
			out = append(out, Discrepancy{
				Kind:   "claimed-but-uncalled",
				Client: e.Client, Op: e.Operation,
				Detail: "manifest claims a command, but no call site calls the operation",
			})
			continue
		}
		if e.Command == "" {
			continue
		}
		var paths []string
		for _, s := range sites {
			if n, ok := tree.Node(s.Pkg, s.Command); ok {
				paths = append(paths, n.String())
			}
		}
		if len(paths) > 0 && !contains(paths, e.Command) {
			out = append(out, Discrepancy{
				Kind:   "command-moved",
				Client: e.Client, Op: e.Operation,
				Detail: fmt.Sprintf("manifest says %q, code says %s", e.Command, strings.Join(paths, ", ")),
			})
		}
	}

	for key, sites := range called {
		client, op, _ := strings.Cut(key, ".")
		if _, ok := m.Lookup(client, op); ok {
			continue
		}
		if _, known := ix.Op(client, op); !known {
			out = append(out, Discrepancy{
				Kind:   "unknown-operation",
				Client: client, Op: op,
				Detail: fmt.Sprintf("called at %s:%d but no such operation in the client", sites[0].File, sites[0].Line),
			})
			continue
		}
		out = append(out, Discrepancy{
			Kind:   "unlisted",
			Client: client, Op: op,
			Detail: fmt.Sprintf("called at %s:%d but absent from the manifest", sites[0].File, sites[0].Line),
		})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Kind != out[j].Kind {
			return out[i].Kind < out[j].Kind
		}
		return out[i].Client+out[i].Op < out[j].Client+out[j].Op
	})
	return out
}

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

// Seed builds manifest entries from the operations the command tree already
// calls. It is how the hand-written commands get into the manifest without
// anyone transcribing sixteen rows by hand.
func Seed(ix *Index, tree *Tree, sites []CallSite) []Entry {
	seen := map[string]bool{}
	var out []Entry
	for _, s := range sites {
		if seen[s.Key()] {
			continue
		}
		op, ok := ix.Op(s.Client, s.Op)
		if !ok {
			continue
		}
		seen[s.Key()] = true

		command := ""
		if n, ok := tree.Node(s.Pkg, s.Command); ok {
			command = n.String()
		}
		out = append(out, Entry{
			Client:    s.Client,
			Operation: s.Op,
			Method:    op.Method,
			Path:      op.Path,
			Command:   command,
			Status:    StatusImplemented,
			Source:    SourceHandwritten,
		})
	}
	return out
}
