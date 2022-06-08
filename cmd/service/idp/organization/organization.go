package organization

import (
	"hcloud-api-client/cmd/service/idp/organization/member"

	"github.com/spf13/cobra"
)

var organizationCmd = &cobra.Command{
	Use:   "organization",
	Short: "manage organizations at the helmut.cloud identity provider",
}

var page, limit int
var id, name, company string

func Init(idpCmd *cobra.Command) {
	idpCmd.AddCommand(organizationCmd)

	listOrganizationCmd.PersistentFlags().IntVarP(&page, "page", "p", 0, "the page of a paginated request (0 equals first page)")
	listOrganizationCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 100, "the amount of results")
	organizationCmd.AddCommand(listOrganizationCmd)

	createOrganizationCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "the name of the new organization")
	createOrganizationCmd.PersistentFlags().StringVarP(&company, "company", "c", "", "the company of the new organization")
	createOrganizationCmd.MarkPersistentFlagRequired("name")
	organizationCmd.AddCommand(createOrganizationCmd)

	getOrganizationCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "the id of the organization")
	getOrganizationCmd.MarkPersistentFlagRequired("id")
	organizationCmd.AddCommand(getOrganizationCmd)

	deleteOrganizationCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "the id of the organization")
	deleteOrganizationCmd.MarkPersistentFlagRequired("id")
	organizationCmd.AddCommand(deleteOrganizationCmd)

	// init organization members
	member.Init(organizationCmd)
}

// 	{
// 		Use:   "member",
// 		Short: "manage organization members",
// 		// Run:          members,
// 	},
// }

// func createOrganization(cmd *cobra.Command, args []string) {
// 	ctx := config.Config.GetActiveContext()
// 	idp := idp.NewFromConfig(&hcloud.Config{Api: ctx.Server, Token: ctx.Token})
// 	organizations, err := idp.CreateOrganization()
// 	if err != nil {
// 		pkg.PrintErr(err)
// 	}

// 	pkg.Print(organizations)
// }
