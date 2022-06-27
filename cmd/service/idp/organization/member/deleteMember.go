package member

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/spf13/cobra"
)

var delOrganizationMemberCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a member of an organization",
	Run:   delOrganizationMember,
}

func init() {
	delOrganizationMemberCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "the id of the organization")
	delOrganizationMemberCmd.MarkPersistentFlagRequired("id")
	organizationMemberCmd.AddCommand(delOrganizationMemberCmd)
}

func delOrganizationMember(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	idp := idp.New(hcloud.New(&hcloud.ClientConfig{Api: ctx.Server, Token: ctx.Token}))
	err := idp.DeleteOrganizationMemberById(id)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(pkg.OkResponse{Result: "member deleted from organization"})
}
