package stream

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/spf13/cobra"
)

var deleteStreamCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a stream by it's ID",
	Run:   deleteStream,
}

func init() {
	deleteStreamCmd.PersistentFlags().StringVarP(&id, "streamId", "i", "", "the id of the stream")
	deleteStreamCmd.MarkPersistentFlagRequired("streamId")

	streamCmd.AddCommand(deleteStreamCmd)
}

func deleteStream(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.ClientConfig{Api: ctx.Server, Token: ctx.Token}))

	err := high5.DeleteStreamById(id)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(pkg.OkResponse{Result: "deleted"})
}
