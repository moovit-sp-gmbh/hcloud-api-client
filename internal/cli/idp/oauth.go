package idp

import (
	"context"
	"encoding/json"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	oauthCommand.AddCommand(oauthWellKnownCommand)
}

var oauthCommand = &cobra.Command{
	Use:   "oauth",
	Short: "Commands for oauth",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var oauthWellKnownCommand = &cobra.Command{
	Use: "well-known",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
		if err != nil {
			return err
		}

		res, err := client.GetV1WellKnownOauthAuthorizationServerWithResponse(context.TODO())
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
