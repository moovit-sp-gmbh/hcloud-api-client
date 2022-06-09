package idp

import (
	"fmt"
	"hcloud-api-client/pkg"
	"runtime"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "register a new account at the helmut.cloud identity provider",
	Run:   register,
}

func init() {
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
