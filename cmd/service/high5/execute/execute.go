package execute

import (
	"github.com/spf13/cobra"
)

var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "execute a stream or event",
}

var id, name, target, data, dataType string
var wait, dataStdin bool
var timeout int

func Init(high5Cmd *cobra.Command) {
	high5Cmd.AddCommand(executeCmd)
}
