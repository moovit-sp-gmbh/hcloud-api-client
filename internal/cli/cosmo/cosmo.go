package cosmo

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	cosmoCommand := &cobra.Command{
		Use:   "cosmo",
		Short: "Entrypoint to Cosmo commands",
		Long:  "Entrypoint to helmut.cloud Cosmo's commands",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cosmoCommand.PersistentFlags().String("server", "https://eu-central-1.app.helmut.cloud", "Set the target server/cluster")

	cosmoCommand.AddCommand(assetsCommand)
	cosmoCommand.AddCommand(eventsCommand)
	cosmoCommand.AddCommand(spacesCommand)
	cosmoCommand.AddCommand(statsCommand)
	cosmoCommand.AddCommand(permissionsCommand)
	cosmoCommand.AddCommand(userCommand)

	return cosmoCommand
}
