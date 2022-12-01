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
