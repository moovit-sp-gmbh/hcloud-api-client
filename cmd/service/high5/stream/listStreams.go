package stream

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/spf13/cobra"
)

var listStreamsCmd = &cobra.Command{
	Use:   "list",
	Short: "list all stream for the given event",
	Run:   listStreams,
}

func init() {
	listStreamsCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "the id of the event")
	listStreamsCmd.MarkPersistentFlagRequired("id")

	listStreamsCmd.PersistentFlags().IntVarP(&page, "page", "p", 0, "the page of a paginated request (0 equals first page)")
	listStreamsCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 500, "the amount of results")
	streamCmd.AddCommand(listStreamsCmd)
}

func listStreams(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.ClientConfig{Api: ctx.Server, Token: ctx.Token}))

	apps, err := high5.GetStreams(id, page, limit)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(apps)
}
