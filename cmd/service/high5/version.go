package high5

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "request the version number of helmut.cloud high5",
	Run:   version,
}

func init() {
	high5Cmd.AddCommand(versionCmd)
}

func version(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.Config{Api: ctx.Server}))
	version, err := high5.Version()
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(version)
}
