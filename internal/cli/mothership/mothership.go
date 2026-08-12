package mothership

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	mothershipCommand := &cobra.Command{
		Use:   "mothership",
		Short: "Entrypoint to Mothership commands",
		Long:  `Entrypoint to helmut.cloud Mothership's commands`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	mothershipCommand.PersistentFlags().String("server", "https://eu-central-1.app.helmut.cloud", "Set the target server/cluster")

	mothershipCommand.AddCommand(agentCommand)
	mothershipCommand.AddCommand(registerCommand)
	mothershipCommand.AddCommand(helloCommand)

	return mothershipCommand
}
