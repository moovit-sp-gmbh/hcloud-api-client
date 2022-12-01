package event

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/spf13/cobra"
)

var deleteEventCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete an event by it's ID",
	Run:   deleteEventById,
}

func init() {
	deleteEventCmd.PersistentFlags().StringVarP(&id, "eventId", "i", "", "the id of the event")
	deleteEventCmd.MarkPersistentFlagRequired("eventId")
	eventCmd.AddCommand(deleteEventCmd)
}

func deleteEventById(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.Config{Api: ctx.Server, Token: ctx.Token}))

	err := high5.DeleteEventById(id)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(pkg.OkResponse{Result: "deleted"})
}
