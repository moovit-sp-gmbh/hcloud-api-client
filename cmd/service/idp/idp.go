package idp

import (
	"hcloud-api-client/cmd/service/idp/organization"
	"runtime"

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

	idpCmd.AddCommand(authorizeCmd)

	idpCmd.AddCommand(versionCmd)

	registerCmd.PersistentFlags().StringVarP(&server, "server", "s", "http://app.helmut.cloud", "the helmut.cloud api url")
	registerCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "the name to use to register at the identity provider")
	registerCmd.PersistentFlags().StringVarP(&email, "email", "e", "", "the email to use to register at the identity provider")
	registerCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "the password to use register at the identity provider")
	if runtime.GOOS != "windows" {
		registerCmd.PersistentFlags().BoolVarP(&passwordStdin, "password-stdin", "", false, "the password to use register at the identity provider read from stdin")
	}
	registerCmd.PersistentFlags().StringVarP(&company, "company", "c", "", "the company to use to register at the identity provider")
	registerCmd.MarkPersistentFlagRequired("server")
	registerCmd.MarkPersistentFlagRequired("email")
	registerCmd.MarkPersistentFlagRequired("name")
	idpCmd.AddCommand(registerCmd)

	authenticateCmd.PersistentFlags().StringVarP(&server, "server", "s", "http://app.helmut.cloud", "the helmut.cloud api url")
	authenticateCmd.PersistentFlags().StringVarP(&email, "email", "e", "", "the email to use to authenticate against the identity provider")
	authenticateCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "the password to use to authenticate against the identity provider")
	if runtime.GOOS != "windows" {
		authenticateCmd.PersistentFlags().BoolVarP(&passwordStdin, "password-stdin", "", false, "the password to use register at the identity provider read from stdin")
	}
	authenticateCmd.PersistentFlags().StringVarP(&identifier, "identifier", "i", "", "an identifier to tag this authentication context")
	authenticateCmd.PersistentFlags().BoolVarP(&setContext, "set-context", "c", true, "set this authentication context as actice")
	authenticateCmd.MarkPersistentFlagRequired("server")
	authenticateCmd.MarkPersistentFlagRequired("email")
	idpCmd.AddCommand(authenticateCmd)

	// init organizations
	organization.Init(idpCmd)
}
