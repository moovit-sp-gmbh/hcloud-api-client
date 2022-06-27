package idp

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/spf13/cobra"
)

var authorizeCmd = &cobra.Command{
	Use:   "authorize",
	Short: "authorize against the helmut.cloud identity provider with the current context",
	Run:   authorize,
}

func init() {
	idpCmd.AddCommand(authorizeCmd)
}

func authorize(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	idp := idp.New(hcloud.New(&hcloud.ClientConfig{Api: ctx.Server, Token: ctx.Token}))
	user, err := idp.Authorize()
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(user)
}
