package idp

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/idpClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/oapi-codegen/runtime/types"
	"github.com/spf13/cobra"
)

func init() {
	registerCreateCommand.Flags().String("validationCode", "", "Optional validation code (required for complete registration that was not triggered before)")
	registerCreateCommand.Flags().String("captcha", "", "Valid google reCAPTCHAV2")
	registerCreateCommand.Flags().String("email", "", "Email of the user")
	registerCreateCommand.Flags().String("name", "", "Name of the user (only alphanumeric characters, underscores, hyphens, dots and spaces are allowed; it should contain between 3 and 255 characters)")
	registerCreateCommand.Flags().String("password", "", "Password of the user (only alphanumeric, underscores, hyphens, dots and special characters are allowed; it should contain between 8 and 255 characters)")
	registerCreateCommand.Flags().String("company", "", "Name of the company the organization belongs to (only alphanumeric characters, underscores, hyphens, dots and spaces are allowed; it should contain between 1 and 255 characters)")
	registerCreateCommand.Flags().String("language", "", "")
	registerCreateCommand.Flags().String("targetUrl", "", "URL to create the redirect link in the email. Defaults to https://app.helmut.cloud.")
	registerCreateCommand.Flags().String("timezone", "", "Define the user's timezone")
	registerResendRegistrationEmailCommand.Flags().String("email", "", "Email of the user")
	registerResendRegistrationEmailCommand.Flags().String("targetUrl", "", "URL to create the redirect link in the email. It will default to https://app.helmut.cloud.")
	registerVerifyCommand.Flags().String("email", "", "Email of the user")
	registerVerifyCommand.Flags().String("verificationCode", "", "Verification code")

	registerCommand.AddCommand(registerCreateCommand)
	registerCommand.AddCommand(registerResendRegistrationEmailCommand)
	registerCommand.AddCommand(registerVerifyCommand)
}

var registerCommand = &cobra.Command{
	Use:   "register",
	Short: "Commands for register",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var registerCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Register a user in hcloud system",
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

		params := &idpClient.PostV1RegisterParams{}

		if cmd.Flags().Changed("validationCode") {
			validationCode, err := cmd.Flags().GetString("validationCode")
			if err != nil {
				return err
			}
			validationCodeParsed, err := uuid.Parse(validationCode)
			if err != nil {
				return err
			}
			params.ValidationCode = &validationCodeParsed
		}

		captcha, err := command.FlagOrPromptString(cmd, "captcha")
		if err != nil {
			return err
		}

		email, err := command.GetArg(cmd, command.StringFlagIfSet("email"), command.ContextIfValidString(latest.Email), command.PromptString("email"))
		if err != nil {
			return err
		}

		name, err := command.FlagOrPromptString(cmd, "name")
		if err != nil {
			return err
		}

		password, err := command.FlagOrPromptString(cmd, "password")
		if err != nil {
			return err
		}

		body := idpClient.PostV1RegisterJSONRequestBody{
			Captcha:  captcha,
			Email:    types.Email(email),
			Name:     name,
			Password: password,
		}

		if cmd.Flags().Changed("company") {
			company, err := cmd.Flags().GetString("company")
			if err != nil {
				return err
			}
			body.Company = &company
		}

		if cmd.Flags().Changed("language") {
			language, err := cmd.Flags().GetString("language")
			if err != nil {
				return err
			}
			languageValue := idpClient.UserRegistrationDtoLanguage(language)
			body.Language = &languageValue
		}

		if cmd.Flags().Changed("targetUrl") {
			targetUrl, err := cmd.Flags().GetString("targetUrl")
			if err != nil {
				return err
			}
			body.TargetUrl = &targetUrl
		}

		if cmd.Flags().Changed("timezone") {
			timezone, err := cmd.Flags().GetString("timezone")
			if err != nil {
				return err
			}
			body.Timezone = &timezone
		}

		res, err := client.PostV1RegisterWithResponse(context.TODO(), params, body)
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

var registerResendRegistrationEmailCommand = &cobra.Command{
	Use:   "resend-registration-email",
	Short: "Generate a new email verification code and resend the verification email",
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

		body := idpClient.PatchV1RegisterResendRegistrationEmailJSONRequestBody{
			Email: types.Email(email),
		}

		if cmd.Flags().Changed("targetUrl") {
			targetUrl, err := cmd.Flags().GetString("targetUrl")
			if err != nil {
				return err
			}
			body.TargetUrl = &targetUrl
		}

		res, err := client.PatchV1RegisterResendRegistrationEmailWithResponse(context.TODO(), body)
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

var registerVerifyCommand = &cobra.Command{
	Use:   "verify",
	Short: "Verify a previous registration",
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

		verificationCode, err := command.FlagOrPromptString(cmd, "verificationCode")
		if err != nil {
			return err
		}
		verificationCodeParsed, err := uuid.Parse(verificationCode)
		if err != nil {
			return err
		}

		body := idpClient.PatchV1RegisterVerifyJSONRequestBody{
			Email:            types.Email(email),
			VerificationCode: verificationCodeParsed,
		}

		res, err := client.PatchV1RegisterVerifyWithResponse(context.TODO(), body)
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
