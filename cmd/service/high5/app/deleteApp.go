package app

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/spf13/cobra"
)

var deleteAppCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an app by it's ID",
	Run:   deleteApp,
}

func init() {
	deleteAppCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "the id of the app")
	deleteAppCmd.MarkPersistentFlagRequired("id")
	appCmd.AddCommand(deleteAppCmd)
}

func deleteApp(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.ClientConfig{Api: ctx.Server, Token: ctx.Token}))

	err := high5.DeleteAppById(id)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print("")
}
