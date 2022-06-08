package idp

import (
	"fmt"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "register a new account at the helmut.cloud identity provider",
	Run:   register,
}

func register(cmd *cobra.Command, args []string) {
	idp := idp.NewFromConfig(&hcloud.Config{Api: server})

	if passwordStdin {
		fmt.Scan(&password)
	}

	user, err := idp.Register(name, company, email, password)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(user)
}
