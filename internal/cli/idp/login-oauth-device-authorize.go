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
	loginOauthDeviceAuthorizeCreateCommand.Flags().String("client_id", "", "Client ID of the OAuth application")
	loginOauthDeviceAuthorizeCreateCommand.Flags().String("user_code", "", "User code to validate")
	loginOauthDeviceAuthorizeCreateCommand.Flags().String("status", "", "Status to set")
}

var loginOauthDeviceAuthorizeCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Device flow - Accept/reject authorization",
	Long:  "Endpoint used in the OAuth 2.0 Device Authorization Flow to store the decision made by the user\nRequired scopes: HCLOUD_FULL",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
		if err != nil {
			return err
		}

		clientId, err := command.FlagOrPromptString(cmd, "client_id")
		if err != nil {
			return err
		}

		userCode, err := command.FlagOrPromptString(cmd, "user_code")
		if err != nil {
			return err
		}

		status, err := command.FlagOrPromptString(cmd, "status")
		if err != nil {
			return err
		}

		body := idpClient.PostV1LoginOauthDeviceAuthorizeClientIdUserCodeJSONRequestBody{
			Status: idpClient.OAuthDeviceDecisionStatus(status),
		}

		res, err := client.PostV1LoginOauthDeviceAuthorizeClientIdUserCodeWithResponse(context.TODO(), clientId, userCode, body)
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
