package event

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/spf13/cobra"
)

var getEventCmd = &cobra.Command{
	Use:   "get",
	Short: "get an event by it's ID",
	Run:   getEventById,
}

func init() {
	getEventCmd.PersistentFlags().StringVarP(&id, "eventId", "i", "", "the id of the event")
	getEventCmd.MarkPersistentFlagRequired("eventId")
	eventCmd.AddCommand(getEventCmd)
}

func getEventById(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.Config{Api: ctx.Server, Token: ctx.Token}))

	apps, err := high5.GetEventById(id)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(apps)
}
