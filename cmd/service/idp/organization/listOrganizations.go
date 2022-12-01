package organization

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/spf13/cobra"
)

var listOrganizationCmd = &cobra.Command{
	Use:   "list",
	Short: "list all organizations of the active user",
	Run:   listOrganizations,
}

func init() {
	listOrganizationCmd.PersistentFlags().IntVarP(&page, "page", "p", 0, "the page of a paginated request (0 equals first page)")
	listOrganizationCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 100, "the amount of results")
	organizationCmd.AddCommand(listOrganizationCmd)
}

func listOrganizations(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	idp := idp.New(hcloud.New(&hcloud.ClientConfig{Api: ctx.Server, Token: ctx.Token}))
	organizations, err := idp.ListOrganizations(limit, page)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(organizations)
}
