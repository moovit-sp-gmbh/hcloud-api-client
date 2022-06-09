package organization

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/spf13/cobra"
)

var deleteOrganizationCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete an organizations by id",
	Run:   listOrganizations,
}

func init() {
	deleteOrganizationCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "the id of the organization")
	deleteOrganizationCmd.MarkPersistentFlagRequired("id")
	organizationCmd.AddCommand(deleteOrganizationCmd)
}

func deleteOrganization(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	idp := idp.NewFromConfig(&hcloud.Config{Api: ctx.Server, Token: ctx.Token})
	err := idp.DeleteOrganizationById(id)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(pkg.OkResponse{Result: "organization deleted"})
}
