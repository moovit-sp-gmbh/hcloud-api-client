package idp

import (
	"context"
	"encoding/json"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/idpClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/oapi-codegen/runtime/types"
	"github.com/spf13/cobra"
)

func init() {
	loginVerifyTotpCommand.Flags().String("email", "", "Email of the user")
	loginVerifyTotpCommand.Flags().String("token", "", "Token to match against a TOTP secret if 2FA is enabled")

	loginVerifyCommand.AddCommand(loginVerifyTotpCommand)
}

var loginVerifyCommand = &cobra.Command{
	Use:   "verify",
	Short: "Commands for verify",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var loginVerifyTotpCommand = &cobra.Command{
	Use:   "totp",
	Short: "Endpoint to validate user totp",
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

		email, err := command.GetArg(cmd, command.StringFlagIfSet("email"), command.ContextIfValidString(latest.Email), command.PromptString("email"))
		if err != nil {
			return err
		}

		token, err := command.FlagOrPromptString(cmd, "token")
		if err != nil {
			return err
		}

		body := idpClient.PostV1LoginVerifyTotpJSONRequestBody{
			Email: types.Email(email),
			Token: token,
		}

		res, err := client.PostV1LoginVerifyTotpWithResponse(context.TODO(), body)
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
