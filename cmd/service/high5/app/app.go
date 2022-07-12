package app

import (
	"github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "request, create or delete app(s)",
}

var page, limit int
var id, name, data string

func Init(high5Cmd *cobra.Command) {
	high5Cmd.AddCommand(appCmd)
}
