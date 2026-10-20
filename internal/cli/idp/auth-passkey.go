package idp

import (
	"context"
	"encoding/json"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/idpClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	authPasskeyVerifyCommand.Flags().String("body", "{}", "request body as JSON")

	authPasskeyCommand.AddCommand(authPasskeyOptionsCommand)
	authPasskeyCommand.AddCommand(authPasskeyVerifyCommand)
	authCommand.AddCommand(authPasskeyCommand)
}

var authCommand = &cobra.Command{
	Use:   "auth",
	Short: "Commands for auth",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var authPasskeyCommand = &cobra.Command{
	Use:   "passkey",
	Short: "Commands for passkey",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var authPasskeyOptionsCommand = &cobra.Command{
	Use:   "options",
	Short: "Create WebAuthn authentication options to log in with a passkey (discoverable credential, no prior identification needed)",
	Long:  "No authentication required",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
		if err != nil {
			return err
		}

		res, err := client.PostV1AuthPasskeyOptionsWithResponse(context.TODO())
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

var authPasskeyVerifyCommand = &cobra.Command{
	Use:   "verify",
	Short: "Verify a WebAuthn authentication response and authorize the user against the system",
	Long:  "No authentication required. On success, sets the Authorization header and auth_token cookie",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
		if err != nil {
			return err
		}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body idpClient.PostV1AuthPasskeyVerifyJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1AuthPasskeyVerifyWithResponse(context.TODO(), body)
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
