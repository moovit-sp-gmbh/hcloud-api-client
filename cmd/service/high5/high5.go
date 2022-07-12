package high5

import (
	"github.com/spf13/cobra"
)

var high5Cmd = &cobra.Command{
	Use:   "high5",
	Short: "use helmut.cloud high5 (stream engine)",
}

func Init(serviceCmd *cobra.Command) {
	serviceCmd.AddCommand(high5Cmd)
}
