package codegen

import (
	"fmt"
	"sort"
	"testing"
)

func TestEmitEveryOperation(t *testing.T) {
	ix, schema, specs, ctors := loadFixtures(t)

	fails := map[string]int{}
	warns := map[string]int{}
	ok, bad, blocked := 0, 0, 0
	for _, name := range ix.ClientNames() {
		client := ix.Clients[name]
		for _, op := range client.Sorted() {
			if !op.Callable() {
				// Unreachable through the generated client whatever a command
				// does; coverage reports these separately.
				blocked++
				continue
			}
			r := Derive(op, client, specs, schema)
			g, err := Emit(r, op, schema, ctors)
			if err != nil {
				bad++
				fails[fmt.Sprint(err)]++
				t.Errorf("%s.%s (%s %s): %v", op.Client, op.Name, op.Method, op.Path, err)
				continue
			}
			ok++
			for _, w := range g.Warnings {
				warns[w]++
			}
		}
	}
	fmt.Printf("emitted %d ops, %d failed, %d not callable\n", ok, bad, blocked)

	type kv struct {
		k string
		n int
	}
	dump := func(label string, m map[string]int) {
		var rows []kv
		for k, n := range m {
			rows = append(rows, kv{k, n})
		}
		sort.Slice(rows, func(i, j int) bool { return rows[i].n > rows[j].n })
		if len(rows) == 0 {
			return
		}
		fmt.Printf("\n%s:\n", label)
		for i, r := range rows {
			if i >= 12 {
				fmt.Printf("  ... and %d more kinds\n", len(rows)-i)
				break
			}
			fmt.Printf("  %4d  %s\n", r.n, r.k)
		}
	}
	dump("failures", fails)
	dump("warnings", warns)
}
