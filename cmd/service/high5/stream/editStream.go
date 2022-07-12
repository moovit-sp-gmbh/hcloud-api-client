package stream

import (
	"encoding/base64"
	"fmt"
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"
	"net/url"

	"github.com/spf13/cobra"
)

var editStreamCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit a stream by it's ID in the Stream Designer Studio",
	Run:   editStream,
}

func init() {
	editStreamCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "the id of the stream")
	editStreamCmd.MarkPersistentFlagRequired("id")

	streamCmd.AddCommand(editStreamCmd)
}

func editStream(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()

	pkg.Print(pkg.OkResponse{Result: fmt.Sprintf("Link: %s/high5/designer/#/%s?b64jwt=%s", ctx.Server, id, base64.StdEncoding.EncodeToString([]byte(url.QueryEscape(ctx.Token))))})
}
