package app

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/spf13/cobra"
)

var createAppCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new app in the active organization",
	Run:   createApp,
}

func init() {
	createAppCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "the name of the new organization")
	createAppCmd.MarkPersistentFlagRequired("name")
	appCmd.AddCommand(createAppCmd)
}

func createApp(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.ClientConfig{Api: ctx.Server, Token: ctx.Token}))

	apps, err := high5.CreateApp(name)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(apps)
}
