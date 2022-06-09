package member

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/spf13/cobra"
)

var addOrganizationMemberCmd = &cobra.Command{
	Use:   "add",
	Short: "add a member to an organization",
	Run:   addOrganizationMember,
}

func init() {
	addOrganizationMemberCmd.PersistentFlags().StringVarP(&id, "organization-id", "o", "", "the id of the organization")
	addOrganizationMemberCmd.PersistentFlags().StringVarP(&userid, "user-id", "u", "", "the id of the user")
	addOrganizationMemberCmd.MarkPersistentFlagRequired("organization-id")
	addOrganizationMemberCmd.MarkPersistentFlagRequired("user-id")
	organizationMemberCmd.AddCommand(addOrganizationMemberCmd)
}

func addOrganizationMember(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	idp := idp.NewFromConfig(&hcloud.Config{Api: ctx.Server, Token: ctx.Token})
	organizationMember, err := idp.AddOrganizationMemberById(id, userid)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(organizationMember)
}
