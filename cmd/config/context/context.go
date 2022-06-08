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

	contextCmd.AddCommand(showContextCmd)

	contextCmd.AddCommand(listContextCmd)

	deleteContextCmd.PersistentFlags().StringVarP(&identifier, "identifier", "i", "", "the identifier of the context to be deleted")
	deleteContextCmd.MarkPersistentFlagRequired("identifier")
	contextCmd.AddCommand(deleteContextCmd)

	setContextCmd.PersistentFlags().StringVarP(&identifier, "identifier", "i", "", "the identifier of the context to be set")
	setContextCmd.MarkPersistentFlagRequired("identifier")
	contextCmd.AddCommand(setContextCmd)
}
