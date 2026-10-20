package cosmo

import (
	"context"
	"encoding/json"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	permissionsCommand.AddCommand(permissionsDiscoveryCommand)
}

var permissionsCommand = &cobra.Command{
	Use:   "permissions",
	Short: "Commands for permissions",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var permissionsDiscoveryCommand = &cobra.Command{
	Use:   "discovery",
	Short: "Obtain a list of possible permissions inside the Cosmo service",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewCosmoClient(cmd)
		if err != nil {
			return err
		}

		res, err := client.GetV1PermissionsDiscoveryWithResponse(context.TODO())
		if err != nil {
			return err
		}

		var data any
		err = json.Unmarshal(res.Body, &data)
		if err != nil {
			return err
		}
		text, err := command.JSON(data)
		if err != nil {
			return err
		}

		command.SetResult(cmd, text)

		return nil
	},
}
