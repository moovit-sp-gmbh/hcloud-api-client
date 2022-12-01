package cmd

import (
	_ "embed"
	"hcloud-api-client/ui"

	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:    "tui",
	Short:  "run terimal user interface",
	Run:    initTui,
	Hidden: true,
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

func initTui(cmd *cobra.Command, args []string) {
	ui.NewApp().Init()
}
