package cmd

import (
	"hcloud-api-client/cmd/service/idp"

	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "use a helmut.cloud service",
}

func init() {
	rootCmd.AddCommand(serviceCmd)

	// init idp subcommand
	idp.Init(serviceCmd)
}
