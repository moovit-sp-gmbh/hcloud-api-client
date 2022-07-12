package app

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/spf13/cobra"
)

var listAppsCmd = &cobra.Command{
	Use:   "list",
	Short: "list all apps of the active organization",
	Run:   listApps,
}

func init() {
	listAppsCmd.PersistentFlags().IntVarP(&page, "page", "p", 0, "the page of a paginated request (0 equals first page)")
	listAppsCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 500, "the amount of results")
	appCmd.AddCommand(listAppsCmd)
}

func listApps(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.ClientConfig{Api: ctx.Server, Token: ctx.Token}))

	apps, err := high5.GetApps(page, limit)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(apps)
}
