package stream

import (
	"github.com/spf13/cobra"
)

var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "request, create or delete stream(s)",
}

var page, limit int
var id, name string

func Init(high5Cmd *cobra.Command) {
	high5Cmd.AddCommand(streamCmd)
}
