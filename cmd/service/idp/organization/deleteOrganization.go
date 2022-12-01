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
	Run:   deleteOrganization,
}

func init() {
	deleteOrganizationCmd.PersistentFlags().StringVarP(&id, "organizationId", "i", "", "the id of the organization")
	deleteOrganizationCmd.MarkPersistentFlagRequired("organizationId")
	organizationCmd.AddCommand(deleteOrganizationCmd)
}

func deleteOrganization(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	idp := idp.New(hcloud.New(&hcloud.Config{Api: ctx.Server, Token: ctx.Token}))

	err := idp.DeleteOrganizationById(id)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(pkg.OkResponse{Result: "organization deleted"})
}
