package idp

import (
	"fmt"
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"
	"runtime"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/spf13/cobra"
)

var authenticateCmd = &cobra.Command{
	Use:   "authenticate",
	Short: "authenticate against the helmut.cloud identity provider",
	Run:   authenticate,
}

func init() {
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
}

func authenticate(cmd *cobra.Command, args []string) {
	idp := idp.New(hcloud.New(&hcloud.ClientConfig{Api: server}))
	if passwordStdin {
		fmt.Scan(&password)
	}

	token, err := idp.Authenticate(email, password)
	if err != nil {
		pkg.PrintErr(err)
	}

	if identifier == "" {
		identifier = server
	}

	config.AddContext(identifier, server, email, token.Token)
	if setContext {
		config.SetContext(identifier)
	}

	pkg.Print(pkg.OkResponse{Result: "authenticated"})
}
