package member

import (
	"github.com/spf13/cobra"
)

var organizationMemberCmd = &cobra.Command{
	Use:   "member",
	Short: "manage organizations members at the helmut.cloud identity provider",
}

var id, userid, email, permission string
var page, limit int

func Init(organizationCmd *cobra.Command) {
	organizationCmd.AddCommand(organizationMemberCmd)
}
