package context

import (
	"github.com/spf13/cobra"
)

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "manage the local configuration context",
}

var identifier string

func Init(configCmd *cobra.Command) {
	configCmd.AddCommand(contextCmd)
}
