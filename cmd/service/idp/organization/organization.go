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
