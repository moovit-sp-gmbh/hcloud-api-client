package organization

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/spf13/cobra"
)

var getOrganizationCmd = &cobra.Command{
	Use:   "get",
	Short: "get an organization by id",
	Run:   getOrganization,
}

func getOrganization(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	idp := idp.NewFromConfig(&hcloud.Config{Api: ctx.Server, Token: ctx.Token})
	organization, err := idp.GetOrganizationById(id)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(organization)
}
