package main

import (
	"encoding/json"
	"fmt"
	"text/tabwriter"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/codegen"
	"github.com/spf13/cobra"
)

func coverageCommand() *cobra.Command {
	var (
		clients []string
		missing bool
		verify  bool
		asJSON  bool
	)

	cmd := &cobra.Command{
		Use:   "coverage",
		Short: "Compare the endpoint manifest against the generated clients",
		Long: "Reports, per client, how many operations have a command (implemented), how " +
			"many were deliberately left alone (skipped), how many have neither (missing), " +
			"and how many the manifest claims but the client no longer exposes (stale).\n\n" +
			"Exits non-zero on a stale entry, since that is a command whose operation has " +
			"disappeared from underneath it.",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := codegen.Root()
			if err != nil {
				return err
			}
			ix, err := codegen.LoadIndex(root)
			if err != nil {
				return err
			}
			m, err := codegen.LoadManifest(root)
			if err != nil {
				return err
			}
			specs := codegen.LoadSpecs(root)

			cov := codegen.Compare(ix, m, clients)

			if asJSON {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				if err := enc.Encode(cov); err != nil {
					return err
				}
			} else {
				printCoverage(cmd.OutOrStdout(), cov, specs, missing)
			}

			var discrepancies []codegen.Discrepancy
			if verify {
				tree, err := codegen.LoadTree(root)
				if err != nil {
					return err
				}
				sites, err := codegen.LoadCallSites(root)
				if err != nil {
					return err
				}
				discrepancies = codegen.Verify(ix, m, tree, sites)
				printDiscrepancies(cmd.OutOrStdout(), discrepancies)
			}

			_, _, _, _, stale := cov.Totals()
			if stale > 0 {
				return fmt.Errorf("%s name operations the clients no longer expose", plural(stale, "manifest entry", "manifest entries"))
			}
			if n := len(discrepancies); n > 0 {
				return fmt.Errorf("%s between the manifest and the command tree", plural(n, "discrepancy", "discrepancies"))
			}
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&clients, "client", "c", nil, "limit to these clients (cosmo, high5, idp, dali, mothership)")
	cmd.Flags().BoolVarP(&missing, "missing", "m", false, "list the operations that have no command")
	cmd.Flags().BoolVar(&verify, "verify", false, "cross-check the manifest against the command tree's call sites")
	cmd.Flags().BoolVar(&asJSON, "json", false, "emit the report as JSON")
	return cmd
}

// plural renders a count with the right noun.
func plural(n int, one, many string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, one)
	}
	return fmt.Sprintf("%d %s", n, many)
}

func printCoverage(w interface{ Write([]byte) (int, error) }, cov *codegen.Coverage, specs *codegen.Specs, listMissing bool) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "CLIENT\tOPS\tIMPL\tSKIPPED\tMISSING\tSTALE")
	for _, cc := range cov.Clients {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\t%d\t%d\n",
			cc.Client, cc.Total, len(cc.Implemented), len(cc.Skipped), len(cc.Missing), len(cc.Stale))
	}
	total, impl, skipped, missing, stale := cov.Totals()
	fmt.Fprintf(tw, "\t%d\t%d\t%d\t%d\t%d\n", total, impl, skipped, missing, stale)
	tw.Flush()

	for _, cc := range cov.Clients {
		if len(cc.Stale) == 0 {
			continue
		}
		fmt.Fprintf(w, "\nstale in %s — the manifest claims these but the client has no such operation:\n", cc.Client)
		for _, e := range cc.Stale {
			fmt.Fprintf(w, "  %s  %s\n", e.Operation, e.Command)
		}
	}

	if !listMissing {
		return
	}
	for _, cc := range cov.Clients {
		if len(cc.Missing) == 0 {
			continue
		}
		blocked := 0
		for _, op := range cc.Missing {
			if !op.Callable() {
				blocked++
			}
		}
		if blocked > 0 {
			fmt.Fprintf(w, "\nmissing in %s (%d, %d not callable):\n", cc.Client, len(cc.Missing), blocked)
		} else {
			fmt.Fprintf(w, "\nmissing in %s (%d):\n", cc.Client, len(cc.Missing))
		}
		mtw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		for _, op := range cc.Missing {
			summary := ""
			if specs != nil {
				if info, ok := specs.Lookup(op.Client, op.Method, op.Path); ok {
					summary = info.Summary
					if summary == "" {
						summary = info.Tag
					}
				}
			}
			if !op.Callable() {
				summary = "NOT CALLABLE: " + op.Blocked
			}
			fmt.Fprintf(mtw, "  %s\t%s\t%s\t%s\n", op.Method, op.Path, op.Name, summary)
		}
		mtw.Flush()
	}
}

func printDiscrepancies(w interface{ Write([]byte) (int, error) }, ds []codegen.Discrepancy) {
	if len(ds) == 0 {
		fmt.Fprintln(w, "\nverify: manifest agrees with the command tree")
		return
	}
	fmt.Fprintf(w, "\nverify: %d discrepancies\n", len(ds))
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, d := range ds {
		fmt.Fprintf(tw, "  %s\t%s.%s\t%s\n", d.Kind, d.Client, d.Op, d.Detail)
	}
	tw.Flush()
}
