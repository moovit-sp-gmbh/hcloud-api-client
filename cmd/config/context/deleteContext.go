package context

import (
	"fmt"
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/spf13/cobra"
)

var deleteContextCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a context by identifier",
	Run:   deleteContext,
}

func init() {
	deleteContextCmd.PersistentFlags().StringVarP(&identifier, "identifier", "i", "", "the identifier of the context to be deleted")
	deleteContextCmd.MarkPersistentFlagRequired("identifier")
	contextCmd.AddCommand(deleteContextCmd)
}

func deleteContext(cmd *cobra.Command, args []string) {
	config.DelContext(identifier)
	pkg.Print(pkg.OkResponse{Result: fmt.Sprintf("context %s deleted", identifier)})
}
