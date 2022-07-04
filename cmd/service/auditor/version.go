package auditor

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/auditor"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "request the version number of the helmut.cloud auditor",
	Run:   version,
}

func init() {
	auditorCmd.AddCommand(versionCmd)
}

func version(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	auditor := auditor.New(hcloud.New(&hcloud.ClientConfig{Api: ctx.Server}))
	version, err := auditor.Version()
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(version)
}
