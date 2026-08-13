package idp

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/idpClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	patchCommand.Flags().String("name", "", "New user name")
	patchPwCommand.Flags().String("old", "", "Old password")
	patchPwCommand.Flags().String("new", "", "New password")

	patchGeneralSettingsCommand.Flags().String("last-url", "", "Last URL visited by the user")
	patchGeneralSettingsCommand.Flags().String("last-view", "", "Last view the user was in")
	patchGeneralSettingsCommand.Flags().String("timezone", "", "IANA timezone")
	patchGeneralSettingsCommand.Flags().String("language", "", "New language")
	patchGeneralSettingsCommand.Flags().String("entrypoint", "", "New entrypoint in helmut.cloud")
	patchGeneralSettingsCommand.Flags().String("date-format", "", "New format to use for dates")
	patchGeneralSettingsCommand.Flags().String("theme", "", "Theme to use")

	patchLastVisitCommand.Flags().String("last-url", "", "Last URL visited by the user")
	patchLastVisitCommand.Flags().String("last-view", "", "Last view the user was in")

	patchNotificationSettingsCommand.Flags().Bool("organization-invitation-email", false, "Enable/disable email notification when invited to an organization")
	patchNotificationSettingsCommand.Flags().Bool("system-notifications", false, "Enable/disable system notifications")
	patchNotificationSettingsCommand.Flags().Bool("cosmo-notifications", false, "Enable/disable cosmo notifications")

	activateTotpCommand.Flags().String("token", "", "Token to match against a TOTP secret")
	deleteTotpCommand.Flags().String("token", "", "Token to match against a TOTP secret")

	whoamiCommand.AddCommand(patchCommand)
	whoamiCommand.AddCommand(getUserEmailCommand)
	whoamiCommand.AddCommand(getUserPasskeysCommand)
	whoamiCommand.AddCommand(createPasskeyOptionsCommand)
	whoamiCommand.AddCommand(getUserNatsPermissionsCommand)
	whoamiCommand.AddCommand(getUserGeneralSettingsCommand)
	whoamiCommand.AddCommand(getUserNotificationSettingsCommand)
	whoamiCommand.AddCommand(getUserConnectionSettingsCommand)
	whoamiCommand.AddCommand(patchGeneralSettingsCommand)
	whoamiCommand.AddCommand(patchLastVisitCommand)
	whoamiCommand.AddCommand(patchNotificationSettingsCommand)
	whoamiCommand.AddCommand(patchPwCommand)
	whoamiCommand.AddCommand(revokeOAuthAppCommand)
	whoamiCommand.AddCommand(deleteAllPatsCommand)
	whoamiCommand.AddCommand(getPatCommand)
	whoamiCommand.AddCommand(deletePatCommand)
	whoamiCommand.AddCommand(createTotpCommand)
	whoamiCommand.AddCommand(activateTotpCommand)
	whoamiCommand.AddCommand(deleteTotpCommand)

	whoamiCommand.AddCommand(pingCommand)
	whoamiCommand.AddCommand(logoutCommand)

}

var whoamiCommand = &cobra.Command{
	Use:   "user",
	Short: "Get user info",
	Long:  "Get user info",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		res, err := client.GetV1UserWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}

		var text command.Text
		switch res.StatusCode() {
		case 200:

			text, err = command.JSON(res.JSON200)
			if err != nil {
				return err
			}
		default:
			text, err = command.JSON(string(res.Body))
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var getUserPasskeysCommand = &cobra.Command{
	Use:   "passkeys",
	Short: "List user passkeys",
	Long:  "List the passkeys registered for the authenticated user",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		res, err := client.GetV1UserPasskeyWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}
		var text command.Text

		switch res.StatusCode() {
		case 200:

			text, err = command.JSON(res.JSON200)
			if err != nil {
				return err
			}
		default:
			text, err = command.JSON(string(res.Body))
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var createPasskeyOptionsCommand = &cobra.Command{
	Use:   "passkey-options",
	Short: "Create passkey registration options",
	Long:  "Create WebAuthn registration options to set up a new passkey for the authenticated user",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		res, err := client.PostV1UserPasskeyOptionsWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}
		var text command.Text

		switch res.StatusCode() {
		case 200:

			text, err = command.JSON(res.JSON200)
			if err != nil {
				return err
			}
		default:
			text, err = command.JSON(string(res.Body))
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var getUserEmailCommand = &cobra.Command{
	Use:   "email",
	Short: "Get user email",
	Long:  "Get user email",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)

		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		res, err := client.GetV1UserEmailWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}

		var text command.Text
		switch res.StatusCode() {
		case 200:

			text, err = command.JSON(res.JSON200)
			if err != nil {
				return err
			}
		default:
			text, err = command.JSON(string(res.Body))
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var getUserNatsPermissionsCommand = &cobra.Command{
	Use:   "nats",
	Short: "Get user Nats permissions",
	Long:  "Get user Nats permissions",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		res, err := client.GetV1UserNatsWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}
		var text command.Text

		switch res.StatusCode() {
		case 200:

			text, err = command.JSON(res.JSON200)
			if err != nil {
				return err
			}
		default:
			text, err = command.JSON(string(res.Body))
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var getUserConnectionSettingsCommand = &cobra.Command{
	Use:   "connection-settings",
	Short: "Get user connection settings",
	Long:  "Retrieves the connection settings of the requesting user",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		res, err := client.GetV1UserSettingsConnectionsWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}
		var text command.Text

		switch res.StatusCode() {
		case 200:

			text, err = command.JSON(res.JSON200)
			if err != nil {
				return err
			}
		default:
			text, err = command.JSON(string(res.Body))
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var getUserNotificationSettingsCommand = &cobra.Command{
	Use:   "notification-settings",
	Short: "Get user notification settings",
	Long:  "Retrieves the notification settings of the requesting user",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		res, err := client.GetV1UserSettingsNotificationsWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}
		var text command.Text

		switch res.StatusCode() {
		case 200:

			text, err = command.JSON(res.JSON200)
			if err != nil {
				return err
			}
		default:
			text, err = command.JSON(string(res.Body))
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var getUserGeneralSettingsCommand = &cobra.Command{
	Use:   "settings",
	Short: "Get user general settings",
	Long:  "Get the current general settings of the logged in user",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		res, err := client.GetV1UserSettingsGeneralWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}
		var text command.Text

		switch res.StatusCode() {
		case 200:

			text, err = command.JSON(res.JSON200)
			if err != nil {
				return err
			}
		default:
			text, err = command.JSON(string(res.Body))
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}


var pingCommand = &cobra.Command{
	Use:   "ping",
	Short: "Ping to keep the user session alive",
	Long:  "Notify the server that the user is online, keeping the session alive",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.Server(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server + "/api/account")
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		res, err := client.GetV1UserPing(context.TODO(), func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})
		if err != nil {
			return err
		}
		defer res.Body.Close()

		var text command.Text

		if res.StatusCode != http.StatusNoContent && (res.StatusCode < 200 || res.StatusCode >= 300) {
			data, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}
			text, err := command.JSON(string(data))
			if err != nil {
				return err
			}
			return errors.New(string(text))
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var logoutCommand = &cobra.Command{
    Use:   "logout",
    Short: "Log out of hcloud on all devices",
    Long:  "Resets all user sessions and therewith logs the user out of hcloud on all devices",
    RunE: func(cmd *cobra.Command, args []string) error {
        server, err := command.Server(cmd)
        if err != nil {
            return err
        }

        client, err := idpClient.NewClientWithResponses(server + "/api/account")
        if err != nil {
            return err
        }

        service := internal.New()
        ctx, err := service.GetLatest()
        if err != nil {
            return err
        }

        if !ctx.Jwt.Valid {
            return errors.New("JWT not set. Login first")
        }

        res, err := client.DeleteV1UserSessionsWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
            req.Header.Set("authorization", ctx.Jwt.String)
            return nil
        })
        if err != nil {
            return err
        }

		log.Println(res.StatusCode(), string(res.Body))

		var text command.Text

		if res.StatusCode() != http.StatusNoContent && (res.StatusCode() < 200 || res.StatusCode() >= 300) {
			text, err := command.JSON(string(res.Body))
			if err != nil {
				return err
			}
			return errors.New(string(text))
		}

		command.SetResult(cmd, text)

        return nil
    },
}

var patchGeneralSettingsCommand = &cobra.Command{
	Use:   "patch-settings",
	Short: "Patch user general settings",
	Long:  "Patch the current general settings of the logged in user",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		body := idpClient.PatchV1UserSettingsGeneralJSONRequestBody{}

		if cmd.Flags().Changed("last-url") {
			lastUrl, err := cmd.Flags().GetString("last-url")
			if err != nil {
				return err
			}
			body.LastUrl = &lastUrl
		}

		if cmd.Flags().Changed("last-view") {
			lastView, err := cmd.Flags().GetString("last-view")
			if err != nil {
				return err
			}
			views := idpClient.Views(lastView)
			body.LastView = &views
		}

		if cmd.Flags().Changed("timezone") {
			timezone, err := cmd.Flags().GetString("timezone")
			if err != nil {
				return err
			}
			body.Timezone = &timezone
		}

		if cmd.Flags().Changed("language") {
			language, err := cmd.Flags().GetString("language")
			if err != nil {
				return err
			}
			lang := idpClient.Language(language)
			body.Language = &lang
		}

		if cmd.Flags().Changed("entrypoint") {
			entrypoint, err := cmd.Flags().GetString("entrypoint")
			if err != nil {
				return err
			}
			ep := idpClient.Entrypoint(entrypoint)
			body.Entrypoint = &ep
		}

		if cmd.Flags().Changed("date-format") {
			dateFormat, err := cmd.Flags().GetString("date-format")
			if err != nil {
				return err
			}
			df := idpClient.DateFormat(dateFormat)
			body.DateFormat = &df
		}

		if cmd.Flags().Changed("theme") {
			theme, err := cmd.Flags().GetString("theme")
			if err != nil {
				return err
			}
			th := idpClient.Theme(theme)
			body.Theme = &th
		}

		res, err := client.PatchV1UserSettingsGeneralWithResponse(context.TODO(), body, func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}

		var text command.Text
		switch res.StatusCode() {
		case 200:
			text, err = command.JSON(res.JSON200)
			if err != nil {
				return err
			}
		default:
			text, err = command.JSON(string(res.Body))
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var patchNotificationSettingsCommand = &cobra.Command{
	Use:   "patch-notification-settings",
	Short: "Patch user notification settings",
	Long:  "Updates the notification settings of the requesting user",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		body := idpClient.PatchV1UserSettingsNotificationsJSONRequestBody{}

		if cmd.Flags().Changed("organization-invitation-email") {
			organizationInvitationEmail, err := cmd.Flags().GetBool("organization-invitation-email")
			if err != nil {
				return err
			}
			body.OrganizationInvitationEmail = &organizationInvitationEmail
		}

		if cmd.Flags().Changed("system-notifications") {
			systemNotifications, err := cmd.Flags().GetBool("system-notifications")
			if err != nil {
				return err
			}
			body.SystemNotifications = &systemNotifications
		}

		if cmd.Flags().Changed("cosmo-notifications") {
			cosmoNotifications, err := cmd.Flags().GetBool("cosmo-notifications")
			if err != nil {
				return err
			}
			body.CosmoNotifications = &cosmoNotifications
		}

		res, err := client.PatchV1UserSettingsNotificationsWithResponse(context.TODO(), body, func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}

		var text command.Text
		switch res.StatusCode() {
		case 200:
			text, err = command.JSON(res.JSON200)
			if err != nil {
				return err
			}
		default:
			text, err = command.JSON(string(res.Body))
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var patchLastVisitCommand = &cobra.Command{
	Use:   "patch-last-visit",
	Short: "Patch user last visit",
	Long:  "Patch the last view and/or last URL of the logged in user",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		body := idpClient.PatchV1UserSettingsGeneralLastVisitJSONRequestBody{}

		if cmd.Flags().Changed("last-url") {
			lastUrl, err := cmd.Flags().GetString("last-url")
			if err != nil {
				return err
			}
			body.LastUrl = &lastUrl
		}

		if cmd.Flags().Changed("last-view") {
			lastView, err := cmd.Flags().GetString("last-view")
			if err != nil {
				return err
			}
			views := idpClient.Views(lastView)
			body.LastView = &views
		}

		res, err := client.PatchV1UserSettingsGeneralLastVisitWithResponse(context.TODO(), body, func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}

		var text command.Text
		switch res.StatusCode() {
		case 200:
			text, err = command.JSON(res.JSON200)
			if err != nil {
				return err
			}
		default:
			text, err = command.JSON(string(res.Body))
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var getPatCommand = &cobra.Command{
	Use:   "pat <patId>",
	Short: "Get a PAT",
	Long:  "Get a single pat by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		patId := args[0]

		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		res, err := client.GetV1UserSettingsPatsPatIdWithResponse(context.TODO(), patId, func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}
		var text command.Text

		switch res.StatusCode() {
		case 200:

			text, err = command.JSON(res.JSON200)
			if err != nil {
				return err
			}
		default:
			text, err = command.JSON(string(res.Body))
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var deletePatCommand = &cobra.Command{
	Use:   "delete-pat <patId>",
	Short: "Delete a PAT",
	Long:  "Deletes a users pat by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		patId := args[0]

		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		res, err := client.DeleteV1UserSettingsPatsPatIdWithResponse(context.TODO(), patId, func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}

		var text command.Text
		switch res.StatusCode() {
		case 204:
			text = command.Text("Success")
		default:
			text, err = command.JSON(string(res.Body))
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var createTotpCommand = &cobra.Command{
	Use:   "create-totp",
	Short: "Create user TOTP",
	Long:  "Create a new user TOTP",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		res, err := client.PostV1UserSettingsSecurity2faTotpWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}
		var text command.Text

		switch res.StatusCode() {
		case 201:

			text, err = command.JSON(res.JSON201)
			if err != nil {
				return err
			}
		default:
			text, err = command.JSON(string(res.Body))
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var activateTotpCommand = &cobra.Command{
	Use:   "activate-totp",
	Short: "Activate user TOTP",
	Long:  "Activate a user TOTP",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		token, err := cmd.Flags().GetString("token")
		if err != nil {
			return err
		}

		res, err := client.PatchV1UserSettingsSecurity2faTotpActivateWithResponse(context.TODO(), idpClient.PatchV1UserSettingsSecurity2faTotpActivateJSONRequestBody{
			Token: token,
		}, func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}
		var text command.Text

		switch res.StatusCode() {
		case 200:

			text, err = command.JSON(res.JSON200)
			if err != nil {
				return err
			}
		default:
			text, err = command.JSON(string(res.Body))
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var deleteTotpCommand = &cobra.Command{
	Use:   "delete-totp",
	Short: "Delete user TOTP",
	Long:  "Delete a user TOTP",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		token, err := cmd.Flags().GetString("token")
		if err != nil {
			return err
		}

		res, err := client.DeleteV1UserSettingsSecurity2faTotpWithResponse(context.TODO(), idpClient.DeleteV1UserSettingsSecurity2faTotpJSONRequestBody{
			Token: token,
		}, func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}

		var text command.Text
		switch res.StatusCode() {
		case 204:
			text = command.Text("Success")
		default:
			text, err = command.JSON(string(res.Body))
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var deleteAllPatsCommand = &cobra.Command{
	Use:   "delete-pats",
	Short: "Delete all PATs",
	Long:  "Deletes all users pats",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		res, err := client.DeleteV1UserSettingsPatsWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}

		var text command.Text
		switch res.StatusCode() {
		case 204:
			text = command.Text("Success")
		default:
			text, err = command.JSON(string(res.Body))
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var patchCommand = &cobra.Command{
	Use:   "patch",
	Short: "Patch user info",
	Long:  "Patch user info",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		newName, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		res, err := client.PatchV1UserWithResponse(context.TODO(), idpClient.PatchV1UserJSONRequestBody{
			Name: &newName,
		}, func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}
		log.Println(res.JSON200.Name)

		ctx.NewName = sql.NullString{Valid: true, String: res.JSON200.Name}
		err = service.Push(ctx)
		if err != nil {
			log.Println(err.Error())
		}

		text, err := command.JSON(res.JSON200)
		if err != nil {
			return err
		}
		command.SetResult(cmd, text)

		return nil
	},
}

var revokeOAuthAppCommand = &cobra.Command{
	Use:   "revoke-oauth-app <oAuthAppId>",
	Short: "Revoke OAuth app access",
	Long:  "Revoke OAuth App Access",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		oAuthAppId := args[0]

		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		res, err := client.DeleteV1UserSettingsOauthOAuthAppIdRevokeWithResponse(context.TODO(), oAuthAppId, func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}

		var text command.Text
		switch res.StatusCode() {
		case 204:
			text = command.Text("Success")
		default:
			text, err = command.JSON(string(res.Body))
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var patchPwCommand = &cobra.Command{
	Use:   "patch-password",
	Short: "Patch user password",
	Long:  "Patch user password",
	RunE: func(cmd *cobra.Command, args []string) error {

		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		oldPw, err := cmd.Flags().GetString("old")
		if err != nil {
			return err
		}

		newPw, err := cmd.Flags().GetString("new")
		if err != nil {
			return err
		}

		res, err := client.PatchV1UserPasswordWithResponse(context.TODO(), idpClient.PatchV1UserPasswordJSONRequestBody{
			Old: oldPw,
			New: newPw,
		}, func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}

		var text command.Text
		switch res.StatusCode() {
		case 200:
			text, err = command.JSON(res.JSON200)
			if err != nil {
				return err
			}
		default:
			text, err = command.JSON(string(res.Body))
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}
