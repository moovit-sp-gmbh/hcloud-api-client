package idp

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/idpClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/oapi-codegen/runtime/types"
	"github.com/spf13/cobra"
)

func init() {
	userSettingsSecurity2faTotpDeleteCommand.Flags().String("token", "", "authorization token")
	userSettingsSecurity2faTotpDeleteCommand.Flags().String("bodyToken", "", "Token to match against a TOTP secret")
	userSettingsSecurity2faTotpCreateCommand.Flags().String("token", "", "authorization token")
	userSettingsSecurity2faTotpActivateCommand.Flags().String("token", "", "authorization token")
	userSettingsSecurity2faTotpActivateCommand.Flags().String("bodyToken", "", "Token to match against a TOTP secret")
	userSettingsSecurity2faTotpRecoverCommand.Flags().String("code", "", "Code to match against the TOTP recoverCodes")
	userSettingsSecurity2faTotpRecoverCommand.Flags().String("email", "", "Email of the user")
	userSettingsSecurity2faTotpRecoverCommand.Flags().String("password", "", "Password of the user (only alphanumeric, underscores, hyphens, dots and special characters are allowed; it should contain between 8 and 255 characters)")

	userSettingsSecurity2faTotpCommand.AddCommand(userSettingsSecurity2faTotpDeleteCommand)
	userSettingsSecurity2faTotpCommand.AddCommand(userSettingsSecurity2faTotpCreateCommand)
	userSettingsSecurity2faTotpCommand.AddCommand(userSettingsSecurity2faTotpActivateCommand)
	userSettingsSecurity2faTotpCommand.AddCommand(userSettingsSecurity2faTotpRecoverCommand)
	userSettingsSecurity2faCommand.AddCommand(userSettingsSecurity2faTotpCommand)
	userSettingsSecurityCommand.AddCommand(userSettingsSecurity2faCommand)
}

var userSettingsSecurityCommand = &cobra.Command{
	Use:   "security",
	Short: "Commands for security",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var userSettingsSecurity2faCommand = &cobra.Command{
	Use:   "2fa",
	Short: "Commands for 2fa",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var userSettingsSecurity2faTotpCommand = &cobra.Command{
	Use:   "totp",
	Short: "Commands for totp",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var userSettingsSecurity2faTotpDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete a user TOTP",
	Long:  "Required scopes: IDP_USER_DELETE",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		bodyToken, err := command.FlagOrDifferentPromptString(cmd, "bodyToken", "token")
		if err != nil {
			return err
		}

		body := idpClient.DeleteV1UserSettingsSecurity2faTotpJSONRequestBody{
			Token: bodyToken,
		}

		res, err := client.DeleteV1UserSettingsSecurity2faTotpWithResponse(context.TODO(), body, func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", jwt)
			return nil
		})
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

var userSettingsSecurity2faTotpCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new user TOTP",
	Long:  "Required scopes: IDP_USER_DELETE",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.PostV1UserSettingsSecurity2faTotpWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", jwt)
			return nil
		})
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

var userSettingsSecurity2faTotpActivateCommand = &cobra.Command{
	Use:   "activate",
	Short: "Activate a user TOTP",
	Long:  "Required scopes: IDP_USER_DELETE",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		bodyToken, err := command.FlagOrDifferentPromptString(cmd, "bodyToken", "token")
		if err != nil {
			return err
		}

		body := idpClient.PatchV1UserSettingsSecurity2faTotpActivateJSONRequestBody{
			Token: bodyToken,
		}

		res, err := client.PatchV1UserSettingsSecurity2faTotpActivateWithResponse(context.TODO(), body, func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", jwt)
			return nil
		})
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

var userSettingsSecurity2faTotpRecoverCommand = &cobra.Command{
	Use:   "recover",
	Short: "Delete a user TOTP by using a recover code",
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

		code, err := command.FlagOrPromptString(cmd, "code")
		if err != nil {
			return err
		}
		codeParsed, err := uuid.Parse(code)
		if err != nil {
			return err
		}

		email, err := command.GetArg(cmd, command.StringFlagIfSet("email"), command.ContextIfValidString(latest.Email), command.PromptString("email"))
		if err != nil {
			return err
		}

		password, err := command.FlagOrPromptString(cmd, "password")
		if err != nil {
			return err
		}

		body := idpClient.DeleteV1UserSettingsSecurity2faTotpRecoverJSONRequestBody{
			Code:     codeParsed,
			Email:    types.Email(email),
			Password: password,
		}

		res, err := client.DeleteV1UserSettingsSecurity2faTotpRecoverWithResponse(context.TODO(), body)
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
