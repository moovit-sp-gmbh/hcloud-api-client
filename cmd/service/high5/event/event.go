package event

import (
	"github.com/spf13/cobra"
)

var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "request, create or delete event(s)",
}

var page, limit int
var id, name string

func Init(high5Cmd *cobra.Command) {
	high5Cmd.AddCommand(eventCmd)
}
