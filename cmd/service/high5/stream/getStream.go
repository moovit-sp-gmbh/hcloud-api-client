package stream

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/spf13/cobra"
)

var getStreamCmd = &cobra.Command{
	Use:   "get",
	Short: "get a stream by it's ID",
	Run:   getStream,
}

func init() {
	getStreamCmd.PersistentFlags().StringVarP(&id, "streamId", "i", "", "the id of the stream")
	getStreamCmd.MarkPersistentFlagRequired("streamId")

	streamCmd.AddCommand(getStreamCmd)
}

func getStream(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.Config{Api: ctx.Server, Token: ctx.Token}))

	apps, err := high5.GetStreamById(id)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(apps)
}
