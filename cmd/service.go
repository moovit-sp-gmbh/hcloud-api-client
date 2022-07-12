package cmd

import (
	"hcloud-api-client/cmd/service/auditor"
	"hcloud-api-client/cmd/service/idp"
	"hcloud-api-client/cmd/service/mailer"

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

	// init auditor subcommand
	auditor.Init(serviceCmd)

	// init mailer subcommand
	mailer.Init(serviceCmd)
}
