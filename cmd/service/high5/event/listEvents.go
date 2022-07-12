package event

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/spf13/cobra"
)

var listEventsCmd = &cobra.Command{
	Use:   "list",
	Short: "list all events of the given app",
	Run:   listEvents,
}

func init() {
	listEventsCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "the id of the app")
	listEventsCmd.MarkPersistentFlagRequired("id")

	listEventsCmd.PersistentFlags().IntVarP(&page, "page", "p", 0, "the page of a paginated request (0 equals first page)")
	listEventsCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 500, "the amount of results")
	eventCmd.AddCommand(listEventsCmd)
}

func listEvents(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.ClientConfig{Api: ctx.Server, Token: ctx.Token}))

	apps, err := high5.GetEvents(id, page, limit)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(apps)
}
