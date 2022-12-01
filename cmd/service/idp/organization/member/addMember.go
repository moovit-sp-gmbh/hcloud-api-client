package member

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"
	"strings"

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
	addOrganizationMemberCmd.PersistentFlags().StringVarP(&id, "organizationId", "i", "", "the id of the organization")
	addOrganizationMemberCmd.PersistentFlags().StringVarP(&email, "email", "e", "", "the email of the user")
	addOrganizationMemberCmd.PersistentFlags().StringVarP(&permission, "permission", "p", "READ", "the permission the new user should hold (READ, MANAGE, ADMIN, OWNER; default READ)")
	addOrganizationMemberCmd.MarkPersistentFlagRequired("organizationId")
	addOrganizationMemberCmd.MarkPersistentFlagRequired("email")
	organizationMemberCmd.AddCommand(addOrganizationMemberCmd)
}

func addOrganizationMember(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	idp := idp.New(hcloud.New(&hcloud.Config{Api: ctx.Server, Token: ctx.Token}))

	perm := hcloud.ORGANIZATION_READ

	switch strings.ToUpper(permission) {
	case "MANAGE":
		perm = hcloud.ORGANIZATION_MANAGE
		break
	case "ADMIN":
		perm = hcloud.ORGANIZATION_MANAGE
		break
	case "OWNER":
		perm = hcloud.ORGANIZATION_OWNER
		break
	}

	organizationMember, err := idp.AddOrganizationMemberById(id, email, perm)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(organizationMember)
}
