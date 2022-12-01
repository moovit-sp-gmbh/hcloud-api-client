package mailer

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/mailer"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "request the version number of the helmut.cloud mailer",
	Run:   version,
}

func init() {
	mailerCmd.AddCommand(versionCmd)
}

func version(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	mailer := mailer.New(hcloud.New(&hcloud.Config{Api: ctx.Server}))
	version, err := mailer.Version()
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(version)
}
