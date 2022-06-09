package idp

import (
	"hcloud-api-client/cmd/service/idp/organization"

	"github.com/spf13/cobra"
)

var idpCmd = &cobra.Command{
	Use:   "idp",
	Short: "use the helmut.cloud identity provider",
}

var server, email, password, identifier, name, company string
var passwordStdin, passwordPrompt, setContext bool

func Init(serviceCmd *cobra.Command) {
	serviceCmd.AddCommand(idpCmd)

	// init organization subcommand
	organization.Init(idpCmd)
}
