package app

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/spf13/cobra"
)

var getAppCmd = &cobra.Command{
	Use:   "get",
	Short: "get app by id",
	Run:   getAppById,
}

func init() {
	getAppCmd.PersistentFlags().StringVarP(&id, "appId", "i", "", "the id of the app")
	getAppCmd.MarkPersistentFlagRequired("appId")
	appCmd.AddCommand(getAppCmd)
}

func getAppById(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.ClientConfig{Api: ctx.Server, Token: ctx.Token}))

	apps, err := high5.GetAppById(id)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(apps)
}
