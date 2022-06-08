package organization

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/spf13/cobra"
)

var listOrganizationCmd = &cobra.Command{
	Use:   "list",
	Short: "list all organizations of the active user",
	Run:   listOrganizations,
}

func listOrganizations(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	idp := idp.NewFromConfig(&hcloud.Config{Api: ctx.Server, Token: ctx.Token})
	organizations, err := idp.ListOrganizations(page, limit)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(organizations)
}
