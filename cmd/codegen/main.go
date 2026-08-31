// Command codegen turns a generated OpenAPI client operation into a cobra
// command, and reports which operations still have no command.
//
// The clients expose several hundred operations and only a handful have
// commands, so the two jobs here are finding the gap and closing it without
// hand-writing the same sixty lines of plumbing every time.
//
//	codegen coverage            what is covered, what is left
//	codegen new                 wizard: pick an endpoint, generate its command
//	codegen regen               replay recipes non-interactively
//	codegen seed                seed the manifest from existing commands
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:           "codegen",
		Short:         "Generate hcloud commands from the OpenAPI clients",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(coverageCommand(), seedCommand(), newCommand(), regenCommand())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "codegen:", err)
		os.Exit(1)
	}
}
