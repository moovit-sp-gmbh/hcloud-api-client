package idp

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/idpClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	userSettingsNotificationsListCommand.Flags().String("token", "", "authorization token")
	userSettingsNotificationsPatchCommand.Flags().String("token", "", "authorization token")
	userSettingsNotificationsPatchCommand.Flags().Bool("cosmoNotifications", false, "Enable/disable cosmo notifications")
	userSettingsNotificationsPatchCommand.Flags().Bool("organizationInvitationEmail", false, "Enable/disable email notification when invited to an organization")
	userSettingsNotificationsPatchCommand.Flags().Bool("systemNotifications", false, "Enable/disable system notifications")

	userSettingsNotificationsCommand.AddCommand(userSettingsNotificationsListCommand)
	userSettingsNotificationsCommand.AddCommand(userSettingsNotificationsPatchCommand)
}

var userSettingsNotificationsCommand = &cobra.Command{
	Use:   "notifications",
	Short: "Commands for notifications",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var userSettingsNotificationsListCommand = &cobra.Command{
	Use:   "list",
	Short: "Retrieves the notification settings of the requesting user",
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

		res, err := client.GetV1UserSettingsNotificationsWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
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

var userSettingsNotificationsPatchCommand = &cobra.Command{
	Use:   "patch",
	Short: "Updates the notification settings of the requesting user",
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

		body := idpClient.PatchV1UserSettingsNotificationsJSONRequestBody{}

		if cmd.Flags().Changed("cosmoNotifications") {
			cosmoNotifications, err := cmd.Flags().GetBool("cosmoNotifications")
			if err != nil {
				return err
			}
			body.CosmoNotifications = &cosmoNotifications
		}

		if cmd.Flags().Changed("organizationInvitationEmail") {
			organizationInvitationEmail, err := cmd.Flags().GetBool("organizationInvitationEmail")
			if err != nil {
				return err
			}
			body.OrganizationInvitationEmail = &organizationInvitationEmail
		}

		if cmd.Flags().Changed("systemNotifications") {
			systemNotifications, err := cmd.Flags().GetBool("systemNotifications")
			if err != nil {
				return err
			}
			body.SystemNotifications = &systemNotifications
		}

		res, err := client.PatchV1UserSettingsNotificationsWithResponse(context.TODO(), body, func(_ context.Context, req *http.Request) error {
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
