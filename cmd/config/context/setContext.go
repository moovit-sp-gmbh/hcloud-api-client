package context

import (
	"fmt"
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/spf13/cobra"
)

var setContextCmd = &cobra.Command{
	Use:   "set",
	Short: "set the active context by identifier",
	Run:   setContext,
}

func init() {
	setContextCmd.PersistentFlags().StringVarP(&identifier, "identifier", "i", "", "the identifier of the context to be set")
	setContextCmd.MarkPersistentFlagRequired("identifier")
	contextCmd.AddCommand(setContextCmd)
}

func setContext(cmd *cobra.Command, args []string) {
	config.SetContext(identifier)
	pkg.Print(pkg.OkResponse{Result: fmt.Sprintf("context set to %s", identifier)})
}
