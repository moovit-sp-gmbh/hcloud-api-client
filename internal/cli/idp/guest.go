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
	guestEmailCommand.Flags().String("email", "", "")
	guestCreateCommand.Flags().String("email", "", "Email of the guest")
	guestCreateCommand.Flags().String("hmac", "", "Url-safe base64 representation of the HMAC of the email address")
	guestCreateCommand.Flags().String("name", "", "Name of the guest account")
	guestUpdateCommand.Flags().String("email", "", "Email of the guest")
	guestUpdateCommand.Flags().String("hmac", "", "Url-safe base64 representation of the HMAC of the email address")
	guestCommand.AddCommand(guestEmailCommand)
	guestCommand.AddCommand(guestCreateCommand)
	guestCommand.AddCommand(guestUpdateCommand)
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

var guestCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a guest user",
	Long:  "Create a guest user by providing an email and HMAC",
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

		hmac, err := command.FlagOrPromptString(cmd, "hmac")
		if err != nil {
			return err
		}

		name, err := command.FlagOrPromptString(cmd, "name")
		if err != nil {
			return err
		}

		body := idpClient.PostV1GuestJSONRequestBody{
			Email: types.Email(email),
			Hmac:  hmac,
			Name:  name,
		}

		res, err := client.PostV1GuestWithResponse(context.TODO(), body)
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

var guestUpdateCommand = &cobra.Command{
	Use:   "update",
	Short: "Login a guest user",
	Long:  "Login as a guest user by providing an email and HMAC",
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

		hmac, err := command.FlagOrPromptString(cmd, "hmac")
		if err != nil {
			return err
		}

		body := idpClient.PutV1GuestJSONRequestBody{
			Email: types.Email(email),
			Hmac:  hmac,
		}

		res, err := client.PutV1GuestWithResponse(context.TODO(), body)
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
