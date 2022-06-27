package member

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/spf13/cobra"
)

var listOrganizationMembersCmd = &cobra.Command{
	Use:   "list",
	Short: "list all organization members",
	Run:   listOrganizationMembers,
}

func init() {
	listOrganizationMembersCmd.PersistentFlags().IntVarP(&page, "page", "p", 0, "the page of a paginated request (0 equals first page)")
	listOrganizationMembersCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 100, "the amount of results")
	listOrganizationMembersCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "the id of the organization")
	listOrganizationMembersCmd.MarkPersistentFlagRequired("id")
	organizationMemberCmd.AddCommand(listOrganizationMembersCmd)
}

func listOrganizationMembers(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	idp := idp.New(hcloud.New(&hcloud.ClientConfig{Api: ctx.Server, Token: ctx.Token}))
	members, err := idp.ListOrganizationMembersById(id, page, limit)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(members)
}
