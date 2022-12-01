package stream

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/spf13/cobra"
)

var createStreamCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new stream for the given event",
	Run:   createStream,
}

func init() {
	createStreamCmd.PersistentFlags().StringVarP(&id, "eventId", "i", "", "the id of the event")
	createStreamCmd.MarkPersistentFlagRequired("eventId")

	createStreamCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "the name of the new stream")
	createStreamCmd.MarkPersistentFlagRequired("name")

	streamCmd.AddCommand(createStreamCmd)
}

func createStream(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.Config{Api: ctx.Server, Token: ctx.Token}))

	stream, err := high5.CreateStream(name, id)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(stream)
}
