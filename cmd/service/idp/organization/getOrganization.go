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

func init() {
	getOrganizationCmd.PersistentFlags().StringVarP(&id, "organizationId", "i", "", "the id of the organization")
	getOrganizationCmd.MarkPersistentFlagRequired("organizationId")
	organizationCmd.AddCommand(getOrganizationCmd)
}

func getOrganization(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	idp := idp.New(hcloud.New(&hcloud.ClientConfig{Api: ctx.Server, Token: ctx.Token}))

	organization, err := idp.GetOrganizationById(id)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(organization)
}
