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

func setContext(cmd *cobra.Command, args []string) {
	config.SetContext(identifier)
	pkg.Print(pkg.OkResponse{Result: fmt.Sprintf("context set to %s", identifier)})
}
