package idp

import (
	"fmt"
	"hcloud-api-client/pkg"
	"runtime"
	"syscall"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "register a new account at the helmut.cloud identity provider",
	Run:   register,
}

func register(cmd *cobra.Command, args []string) {
	idp := idp.NewFromConfig(&hcloud.Config{Api: server})

	if runtime.GOOS != "windows" {
		if passwordStdin {
			fmt.Scan(&password)
		}

		if passwordPrompt {
			fmt.Print("Please enter password:\n")
			bytepw, err := term.ReadPassword(syscall.Stdin)
			if err != nil {
				pkg.PrintErr(&hcloud.ErrorResponse{Code: -1, Message: err.Error()})
			}
			password = string(bytepw)
		}
	}

	user, err := idp.Register(name, company, email, password)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(user)
}
