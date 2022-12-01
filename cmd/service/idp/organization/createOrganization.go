package organization

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/spf13/cobra"
)

var createOrganizationCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new organization as the active user",
	Run:   createOrganization,
}

func init() {
	createOrganizationCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "the name of the new organization")
	createOrganizationCmd.PersistentFlags().StringVarP(&company, "company", "c", "", "the company of the new organization")
	createOrganizationCmd.MarkPersistentFlagRequired("name")
	organizationCmd.AddCommand(createOrganizationCmd)
}

func createOrganization(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	idp := idp.New(hcloud.New(&hcloud.ClientConfig{Api: ctx.Server, Token: ctx.Token}))

	organization, err := idp.CreateOrganization(name, company)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(organization)
}
