package idp

import (
	"fmt"
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/spf13/cobra"
)

var authenticateCmd = &cobra.Command{
	Use:   "authenticate",
	Short: "authenticate against the helmut.cloud identity provider",
	Run:   authenticate,
}

func authenticate(cmd *cobra.Command, args []string) {
	idp := idp.NewFromConfig(&hcloud.Config{Api: server})

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
