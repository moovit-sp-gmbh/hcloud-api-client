package high5

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	high5Command := &cobra.Command{
		Use:   "high5",
		Short: "Entrypoint to High5 commands",
		Long:  "Entrypoint to helmut.cloud High5's commands",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	high5Command.PersistentFlags().String("server", "https://eu-central-1.app.helmut.cloud", "Set the target server/cluster")

	high5Command.AddCommand(spaceCommand)

	return high5Command
}
