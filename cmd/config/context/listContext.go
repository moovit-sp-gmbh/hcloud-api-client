package context

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/spf13/cobra"
)

var listContextCmd = &cobra.Command{
	Use:   "list",
	Short: "list all available contexts",
	Run:   listContext,
}

func listContext(cmd *cobra.Command, args []string) {
	pkg.Print(config.Config.Contexts)
}
