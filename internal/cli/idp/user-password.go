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
	userPasswordResetCommand.Flags().String("email", "", "Email of the user")
	userPasswordResetCommand.Flags().String("newPassword", "", "New password")
	userPasswordResetCommand.Flags().String("nonce", "", "Nonce generated in the 1st phase of the password reset process.")
	userPasswordResetCommand.Flags().String("totpToken", "", "TOTP token. Required only if the user has 2FA enabled.")

	userPasswordCommand.AddCommand(userPasswordResetCommand)
}

var userPasswordCommand = &cobra.Command{
	Use:   "password",
	Short: "Commands for password",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var userPasswordResetCommand = &cobra.Command{
	Use:   "reset",
	Short: "Reset a user's password",
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

		newPassword, err := command.FlagOrPromptString(cmd, "newPassword")
		if err != nil {
			return err
		}

		nonce, err := command.FlagOrPromptString(cmd, "nonce")
		if err != nil {
			return err
		}

		body := idpClient.PostV1UserPasswordResetJSONRequestBody{
			Email:       types.Email(email),
			NewPassword: newPassword,
			Nonce:       nonce,
		}

		if cmd.Flags().Changed("totpToken") {
			totpToken, err := cmd.Flags().GetString("totpToken")
			if err != nil {
				return err
			}
			body.TotpToken = &totpToken
		}

		res, err := client.PostV1UserPasswordResetWithResponse(context.TODO(), body)
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
