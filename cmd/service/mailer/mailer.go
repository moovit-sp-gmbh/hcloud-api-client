package mailer

import (
	"github.com/spf13/cobra"
)

var mailerCmd = &cobra.Command{
	Use:   "mailer",
	Short: "use the helmut.cloud mailer",
}

func Init(serviceCmd *cobra.Command) {
	serviceCmd.AddCommand(mailerCmd)
}
