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

func listOrganizationMembers(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	idp := idp.NewFromConfig(&hcloud.Config{Api: ctx.Server, Token: ctx.Token})
	members, err := idp.ListOrganizationMembersById(id, page, limit)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(members)
}
