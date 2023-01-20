package cmd

import (
	_ "embed"
	"fmt"
	"hcloud-api-client/internal/view"

	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:    "tui",
	Short:  "run terminal user interface",
	Run:    initTui,
	Hidden: true,
}

var debug bool

func init() {
	tuiCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "output debug information")
	rootCmd.AddCommand(tuiCmd)
}

func initTui(cmd *cobra.Command, args []string) {

	app := view.NewApp(debug)

	if err := app.Init(); err != nil {
		panic(fmt.Sprintf("app init failed -- %v", err))
	}
	if err := app.Run(); err != nil {
		panic(fmt.Sprintf("app run failed %v", err))
	}
	if view.ExitStatus != "" {
		panic(fmt.Sprintf("view exit status %s", view.ExitStatus))
	}
}
