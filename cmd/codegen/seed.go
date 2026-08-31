package main

import (
	"fmt"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/codegen"
	"github.com/spf13/cobra"
)

func seedCommand() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "seed",
		Short: "Add the hand-written commands' endpoints to the manifest",
		Long: "Scans the command tree for client calls and records what it finds. Run once to " +
			"bootstrap the manifest; afterwards `codegen new` maintains it and " +
			"`codegen coverage --verify` checks it.",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := codegen.Root()
			if err != nil {
				return err
			}
			ix, err := codegen.LoadIndex(root)
			if err != nil {
				return err
			}
			tree, err := codegen.LoadTree(root)
			if err != nil {
				return err
			}
			sites, err := codegen.LoadCallSites(root)
			if err != nil {
				return err
			}
			m, err := codegen.LoadManifest(root)
			if err != nil {
				return err
			}

			added := 0
			for _, e := range codegen.Seed(ix, tree, sites) {
				if _, exists := m.Lookup(e.Client, e.Operation); exists {
					continue
				}
				m.Upsert(e)
				added++
				fmt.Fprintf(cmd.OutOrStdout(), "+ %-18s %-6s %-42s %s\n", e.Client, e.Method, e.Path, e.Command)
			}

			if added == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "manifest already covers every call site")
				return nil
			}
			if dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "\n%d entries would be added (dry run)\n", added)
				return nil
			}
			if err := m.Save(root); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\n%d entries added to %s\n", added, codegen.ManifestPath)
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "report what would be added without writing")
	return cmd
}
