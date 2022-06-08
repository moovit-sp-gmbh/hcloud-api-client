package context

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/spf13/cobra"
)

var showContextCmd = &cobra.Command{
	Use:   "show",
	Short: "show the active context",
	Run:   showContext,
}

func showContext(cmd *cobra.Command, args []string) {
	pkg.Print(config.Config.GetActiveContext())
}
