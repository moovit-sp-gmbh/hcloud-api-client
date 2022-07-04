package auditor

import (
	"github.com/spf13/cobra"
)

var auditorCmd = &cobra.Command{
	Use:   "auditor",
	Short: "use the helmut.cloud auditor",
}

var organization string
var limit, page int

func Init(serviceCmd *cobra.Command) {
	serviceCmd.AddCommand(auditorCmd)
}
