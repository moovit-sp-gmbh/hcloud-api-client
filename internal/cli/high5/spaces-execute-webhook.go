package high5

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/high5Client"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	spacesExecuteWebhookListCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesExecuteWebhookListCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesExecuteWebhookListCommand.Flags().String("webhookToken", "", "Last path segment of the webhook URL (128 char long hex string)")
	spacesExecuteWebhookListCommand.Flags().String("challenge", "", "Some third party services require to validate the target URL of a webhook first. To do so, they sent an initial GET request with a query-string called 'challenge' (only alphanumeric characters, underscores and hyphens are allowed; it should contain between 1 and 128 characters)")
	spacesExecuteWebhookCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesExecuteWebhookCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesExecuteWebhookCreateCommand.Flags().String("webhookToken", "", "Last path segment of the webhook URL (128 char long hex string)")
	spacesExecuteWebhookCallbackCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesExecuteWebhookCallbackCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesExecuteWebhookCallbackCommand.Flags().String("webhookToken", "", "Last path segment of the webhook URL (128 char long hex string)")
	spacesExecuteWebhookCallbackCommand.Flags().String("callbackUUID", "", "")
	spacesExecuteWebhookCallbackCommand.Flags().String("body", "{}", "request body as JSON")

	spacesExecuteWebhookCommand.AddCommand(spacesExecuteWebhookListCommand)
	spacesExecuteWebhookCommand.AddCommand(spacesExecuteWebhookCreateCommand)
	spacesExecuteWebhookCommand.AddCommand(spacesExecuteWebhookCallbackCommand)
}

var spacesExecuteWebhookCommand = &cobra.Command{
	Use:   "webhook",
	Short: "Commands for webhook",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesExecuteWebhookListCommand = &cobra.Command{
	Use:   "list",
	Short: "Validate webhook url and set headers",
	Long:  "Validate webhook url and set headers. This path is equivalent to the 'url' property of Webhook",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewHigh5Client(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("orgName"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		spaceName, err := command.FlagOrPromptString(cmd, "spaceName")
		if err != nil {
			return err
		}

		webhookToken, err := command.FlagOrPromptString(cmd, "webhookToken")
		if err != nil {
			return err
		}

		params := &high5Client.GetV1OrgOrgNameSpacesSpaceNameExecuteWebhookWebhookTokenParams{}

		if cmd.Flags().Changed("challenge") {
			challenge, err := cmd.Flags().GetString("challenge")
			if err != nil {
				return err
			}
			params.Challenge = &challenge
		}

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameExecuteWebhookWebhookTokenWithResponse(context.TODO(), orgName, spaceName, webhookToken, params)
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

var spacesExecuteWebhookCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Execute webhook",
	Long:  "Executes a webhook. This path is equivalent to the 'url' property of Webhook",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewHigh5Client(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("orgName"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		spaceName, err := command.FlagOrPromptString(cmd, "spaceName")
		if err != nil {
			return err
		}

		webhookToken, err := command.FlagOrPromptString(cmd, "webhookToken")
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameExecuteWebhookWebhookTokenWithResponse(context.TODO(), orgName, spaceName, webhookToken)
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

var spacesExecuteWebhookCallbackCommand = &cobra.Command{
	Use:   "callback",
	Short: "Webhook callback URL to send a non-generic response to webhook triggerer",
	Long:  "Forwards the incoming request defined the in the request body to the webhook triggerer via HTTP response",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewHigh5Client(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("orgName"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		spaceName, err := command.FlagOrPromptString(cmd, "spaceName")
		if err != nil {
			return err
		}

		webhookToken, err := command.FlagOrPromptString(cmd, "webhookToken")
		if err != nil {
			return err
		}

		callbackUUID, err := command.FlagOrPromptString(cmd, "callbackUUID")
		if err != nil {
			return err
		}
		callbackUUIDParsed, err := uuid.Parse(callbackUUID)
		if err != nil {
			return err
		}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body high5Client.PostV1OrgOrgNameSpacesSpaceNameExecuteWebhookWebhookTokenCallbackCallbackUUIDJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameExecuteWebhookWebhookTokenCallbackCallbackUUIDWithResponse(context.TODO(), orgName, spaceName, webhookToken, callbackUUIDParsed, body)
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
