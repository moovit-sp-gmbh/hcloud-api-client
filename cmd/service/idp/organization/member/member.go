package member

import (
	"github.com/spf13/cobra"
)

var organizationMemberCmd = &cobra.Command{
	Use:   "member",
	Short: "manage organizations members at the helmut.cloud identity provider",
}

var id, userid string
var page, limit int

func Init(organizationCmd *cobra.Command) {
	organizationCmd.AddCommand(organizationMemberCmd)

	listOrganizationMembersCmd.PersistentFlags().IntVarP(&page, "page", "p", 0, "the page of a paginated request (0 equals first page)")
	listOrganizationMembersCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 100, "the amount of results")
	listOrganizationMembersCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "the id of the organization")
	listOrganizationMembersCmd.MarkPersistentFlagRequired("id")
	organizationMemberCmd.AddCommand(listOrganizationMembersCmd)

	addOrganizationMemberCmd.PersistentFlags().StringVarP(&id, "organization-id", "o", "", "the id of the organization")
	addOrganizationMemberCmd.PersistentFlags().StringVarP(&userid, "user-id", "u", "", "the id of the user")
	addOrganizationMemberCmd.MarkPersistentFlagRequired("organization-id")
	addOrganizationMemberCmd.MarkPersistentFlagRequired("user-id")
	organizationMemberCmd.AddCommand(addOrganizationMemberCmd)

	delOrganizationMemberCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "the id of the organization")
	delOrganizationMemberCmd.MarkPersistentFlagRequired("id")
	organizationMemberCmd.AddCommand(delOrganizationMemberCmd)
}
