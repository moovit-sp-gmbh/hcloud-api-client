package main

import (
	"fmt"
	"os"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/codegen"
	"github.com/spf13/cobra"
)

func newCommand() *cobra.Command {
	var (
		client   string
		filter   string
		dryRun   bool
		skipSQLC bool
	)

	cmd := &cobra.Command{
		Use:   "new",
		Short: "Generate a command for an endpoint (interactive)",
		Long: "Walks through picking an endpoint, placing it in the command tree, naming each " +
			"argument and choosing what to keep from the response, then writes the command, " +
			"registers it, adds any context migration it needs and records it in the manifest.\n\n" +
			"Every answer is saved to a recipe under codegen/recipes so `codegen regen` can " +
			"reproduce the command without asking again.",
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
			schema, err := codegen.LoadContextSchema(root)
			if err != nil {
				return err
			}
			manifest, err := codegen.LoadManifest(root)
			if err != nil {
				return err
			}
			specs := codegen.LoadSpecs(root)
			if specs.Empty() {
				fmt.Fprintln(cmd.OutOrStdout(),
					"note: no vendored specs under specs/, so nothing can be prefilled from them.\n"+
						"      run `make specs` to fetch them (it does not touch the clients).")
			}

			out := cmd.OutOrStdout()
			w := codegen.NewWizard(os.Stdin, out, ix, tree, specs, schema)

			op, err := w.PickOperation(manifest, client, filter)
			if err != nil {
				return err
			}

			recipe := codegen.Derive(op, ix.Clients[op.Client], specs, schema)
			if err := w.Run(op, recipe); err != nil {
				return err
			}

			res, err := codegen.Apply(root, recipe, codegen.ApplyOptions{DryRun: dryRun, SkipSQLC: skipSQLC})
			if err != nil {
				return err
			}
			report(out, res, dryRun)

			if !dryRun {
				fmt.Fprintln(out, "\nnext: go build ./... && ./bin/codegen coverage --verify")
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&client, "client", "c", "", "service to pick from (cosmo, high5, idp, dali, mothership)")
	cmd.Flags().StringVarP(&filter, "filter", "f", "", "narrow the endpoint list by path, method or operation name")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "report what would be written without writing it")
	cmd.Flags().BoolVar(&skipSQLC, "skip-sqlc", false, "do not run sqlc generate for a new context column")
	return cmd
}

func regenCommand() *cobra.Command {
	var (
		dryRun   bool
		skipSQLC bool
		quiet    bool
	)

	cmd := &cobra.Command{
		Use:   "regen [recipe...]",
		Short: "Regenerate commands from their recipes",
		Long: "Replays saved recipes non-interactively. Use it after changing the emitter or " +
			"regenerating a client, to roll the change across every generated command.\n\n" +
			"With no arguments it replays every recipe under codegen/recipes.",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := codegen.Root()
			if err != nil {
				return err
			}

			var recipes []*codegen.Recipe
			var paths []string
			if len(args) > 0 {
				for _, path := range args {
					r, err := codegen.LoadRecipe(path)
					if err != nil {
						return err
					}
					recipes = append(recipes, r)
					paths = append(paths, path)
				}
			} else if recipes, paths, err = codegen.LoadRecipes(root); err != nil {
				return err
			}

			if len(recipes) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "no recipes under %s yet — run `codegen new` first\n", codegen.RecipesDir)
				return nil
			}

			// Reuse one index across the batch: reparsing 4.4MB of generated
			// client per recipe would dominate a full replay.
			ix, err := codegen.LoadIndex(root)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			for i, r := range recipes {
				if !quiet {
					fmt.Fprintf(out, "\n=== %s ===\n", paths[i])
				}
				res, err := codegen.Apply(root, r, codegen.ApplyOptions{DryRun: dryRun, SkipSQLC: skipSQLC, Index: ix})
				if err != nil {
					return fmt.Errorf("%s: %w", paths[i], err)
				}
				if !quiet {
					report(out, res, dryRun)
				}
			}
			fmt.Fprintf(out, "\nreplayed %d recipes\n", len(recipes))
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "report what would be written without writing it")
	cmd.Flags().BoolVar(&skipSQLC, "skip-sqlc", false, "do not run sqlc generate for a new context column")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "only report the total")
	return cmd
}

func report(w interface{ Write([]byte) (int, error) }, res *codegen.ApplyResult, dryRun bool) {
	verb := "wrote"
	if dryRun {
		verb = "would write"
	}
	fmt.Fprintln(w)
	for _, p := range res.Created {
		fmt.Fprintf(w, "  %s (new)  %s\n", verb, p)
	}
	for _, p := range res.Edited {
		fmt.Fprintf(w, "  %s        %s\n", verb, p)
	}
	for _, c := range res.Columns {
		fmt.Fprintf(w, "  context column %s (db.Context.%s)\n", c.Column, c.GoName)
	}
	if len(res.Warnings) > 0 {
		fmt.Fprintln(w, "\nwarnings:")
		for _, warning := range res.Warnings {
			fmt.Fprintf(w, "  - %s\n", warning)
		}
	}
}
