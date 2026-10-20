package dali

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	daliCommand := &cobra.Command{
		Use:   "dali",
		Short: "Entrypoint to Dali commands",
		Long:  "Entrypoint to helmut.cloud Dali's commands",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	daliCommand.PersistentFlags().String("server", "https://eu-central-1.app.helmut.cloud", "Set the target server/cluster")

	daliCommand.AddCommand(avatarCommand)

	return daliCommand
}
