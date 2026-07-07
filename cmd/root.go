package cmd

import (
	"os"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/cli/context"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/cli/idp"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(idp.Command())
	rootCmd.AddCommand(context.Command())
}

var rootCmd = &cobra.Command{
	Use:   "hcloud",
	Short: "Command-Line Interface for helmut.cloud",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
