package idp

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	loginOauthAppCommand.Flags().String("clientId", "", "Client ID of the OAuth application")
	loginOauthCommand.AddCommand(loginOauthAppCommand)
	loginOauthCommand.AddCommand(loginOauthAuthorizeCommand)
	loginOauthCommand.AddCommand(loginOauthDeviceCommand)
}

var loginOauthCommand = &cobra.Command{
	Use:   "oauth",
	Short: "Commands for oauth",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var loginOauthAppCommand = &cobra.Command{
	Use:   "app",
	Short: "OAuth App public info by clientId",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		clientId, err := command.GetArg(cmd, command.StringFlagIfSet("clientId"), command.ContextIfValidString(latest.OauthAppClientID), command.PromptString("OAuth Client ID"))
		if err != nil {
			return err
		}
		clientIdParsed, err := uuid.Parse(clientId)
		if err != nil {
			return err
		}
		latest.OauthAppClientID = sql.NullString{Valid: true, String: fmt.Sprint(clientId)}

		res, err := client.GetV1LoginOauthAppClientIdWithResponse(context.TODO(), clientIdParsed)
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

		if res.StatusCode() >= 200 && res.StatusCode() < 300 {
			service.Push(latest)
		}

		command.SetResult(cmd, text)

		return nil
	},
}
