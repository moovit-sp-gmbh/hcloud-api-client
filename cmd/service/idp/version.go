package idp

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "request the version number of the helmut.cloud identity provider",
	Run:   version,
}

func init() {
	idpCmd.AddCommand(versionCmd)
}

func version(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	idp := idp.New(hcloud.New(&hcloud.Config{Api: ctx.Server}))
	version, err := idp.Version()
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(version)
}
