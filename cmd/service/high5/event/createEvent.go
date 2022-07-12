package event

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/spf13/cobra"
)

var createEventCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new event in the given app",
	Run:   createEvent,
}

func init() {
	createEventCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "the name of the new organization")
	createEventCmd.MarkPersistentFlagRequired("name")

	createEventCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "the id of the app")
	createEventCmd.MarkPersistentFlagRequired("id")
	eventCmd.AddCommand(createEventCmd)
}

func createEvent(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.ClientConfig{Api: ctx.Server, Token: ctx.Token}))

	apps, err := high5.CreateEvent(name, id)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(apps)
}
