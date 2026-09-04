package idp

import (
	"context"
	"encoding/json"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/oapi-codegen/runtime/types"
	"github.com/spf13/cobra"
)

func init() {
	guestEmailCommand.Flags().String("email", "", "")
	guestCommand.AddCommand(guestEmailCommand)
}

var guestCommand = &cobra.Command{
	Use:   "guest",
	Short: "Commands for guest",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var guestEmailCommand = &cobra.Command{
	Use: "email",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
		if err != nil {
			return err
		}

		email, err := command.FlagOrPromptString(cmd, "email")
		if err != nil {
			return err
		}

		res, err := client.GetV1GuestEmailEmailWithResponse(context.TODO(), types.Email(email))
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
