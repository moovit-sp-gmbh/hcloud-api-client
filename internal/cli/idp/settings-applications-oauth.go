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
	settingsApplicationsOauthCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	settingsApplicationsOauthCreateCommand.Flags().String("token", "", "authorization token")
	settingsApplicationsOauthCreateCommand.Flags().StringSlice("callback", nil, "Callback URI for the OAuthApp")
	settingsApplicationsOauthCreateCommand.Flags().String("name", "", "Name of the OAuthApp")
	settingsApplicationsOauthCreateCommand.Flags().Bool("allowAuthorizationCodeFlow", false, "Set to true to allow authorization code flow through this application")
	settingsApplicationsOauthCreateCommand.Flags().Bool("allowDeviceFlow", false, "Set to true to allow device flow through this application")
	settingsApplicationsOauthCreateCommand.Flags().String("description", "", "Optional description of the OAuthApp")
	settingsApplicationsOauthCreateCommand.Flags().String("homepage", "", "Optional homepage URL for the OAuthApp")
	settingsApplicationsOauthSearchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	settingsApplicationsOauthSearchCommand.Flags().String("token", "", "authorization token")
	settingsApplicationsOauthSearchCommand.Flags().Float32("limit", 0, "Page size")
	settingsApplicationsOauthSearchCommand.Flags().Float32("page", 0, "Page index")
	settingsApplicationsOauthSearchCommand.Flags().String("body", "{}", "request body as JSON")
	settingsApplicationsOauthDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	settingsApplicationsOauthDeleteCommand.Flags().String("oAuthAppId", "", "")
	settingsApplicationsOauthDeleteCommand.Flags().String("token", "", "authorization token")
	settingsApplicationsOauthUpdateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	settingsApplicationsOauthUpdateCommand.Flags().String("oAuthAppId", "", "")
	settingsApplicationsOauthUpdateCommand.Flags().String("token", "", "authorization token")
	settingsApplicationsOauthUpdateCommand.Flags().StringSlice("callback", nil, "Callback URI for the OAuthApp")
	settingsApplicationsOauthUpdateCommand.Flags().String("name", "", "Name of the OAuthApp")
	settingsApplicationsOauthUpdateCommand.Flags().Bool("allowAuthorizationCodeFlow", false, "Set to true to allow authorization code flow through this application")
	settingsApplicationsOauthUpdateCommand.Flags().Bool("allowDeviceFlow", false, "Set to true to allow device flow through this application")
	settingsApplicationsOauthUpdateCommand.Flags().String("description", "", "Optional description of the OAuthApp")
	settingsApplicationsOauthUpdateCommand.Flags().String("homepage", "", "Optional homepage URL for the OAuthApp")
	settingsApplicationsOauthSecretCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	settingsApplicationsOauthSecretCommand.Flags().String("oAuthAppId", "", "")
	settingsApplicationsOauthSecretCommand.Flags().String("token", "", "authorization token")
	settingsApplicationsOauthSecretCommand.Flags().String("secretName", "", "Name of the secret (only alphanumeric, underscores, hyphens and dots are allowed; it should contain between 1 and 255 characters)")

	settingsApplicationsOauthCommand.AddCommand(settingsApplicationsOauthCreateCommand)
	settingsApplicationsOauthCommand.AddCommand(settingsApplicationsOauthSearchCommand)
	settingsApplicationsOauthCommand.AddCommand(settingsApplicationsOauthDeleteCommand)
	settingsApplicationsOauthCommand.AddCommand(settingsApplicationsOauthUpdateCommand)
	settingsApplicationsOauthCommand.AddCommand(settingsApplicationsOauthSecretCommand)
	settingsApplicationsCommand.AddCommand(settingsApplicationsOauthCommand)
	settingsCommand.AddCommand(settingsApplicationsCommand)
	settingsCommand.AddCommand(settingsDomainsCommand)
}

var settingsCommand = &cobra.Command{
	Use:   "settings",
	Short: "Commands for settings",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var settingsApplicationsCommand = &cobra.Command{
	Use:   "applications",
	Short: "Commands for applications",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var settingsApplicationsOauthCommand = &cobra.Command{
	Use:   "oauth",
	Short: "Commands for oauth",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var settingsApplicationsOauthCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new oauth app within an organization",
	Long:  "Required role: ADMIN\nRequired scopes: IDP_ORGANIZATION_ADMIN",
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

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("orgName"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		callback, err := cmd.Flags().GetStringSlice("callback")
		if err != nil {
			return err
		}

		name, err := command.FlagOrPromptString(cmd, "name")
		if err != nil {
			return err
		}

		body := idpClient.PostV1OrgOrgNameSettingsApplicationsOauthJSONRequestBody{
			Callback: callback,
			Name:     name,
		}

		if cmd.Flags().Changed("allowAuthorizationCodeFlow") {
			allowAuthorizationCodeFlow, err := cmd.Flags().GetBool("allowAuthorizationCodeFlow")
			if err != nil {
				return err
			}
			body.AllowAuthorizationCodeFlow = &allowAuthorizationCodeFlow
		}

		if cmd.Flags().Changed("allowDeviceFlow") {
			allowDeviceFlow, err := cmd.Flags().GetBool("allowDeviceFlow")
			if err != nil {
				return err
			}
			body.AllowDeviceFlow = &allowDeviceFlow
		}

		if cmd.Flags().Changed("description") {
			description, err := cmd.Flags().GetString("description")
			if err != nil {
				return err
			}
			body.Description = &description
		}

		if cmd.Flags().Changed("homepage") {
			homepage, err := cmd.Flags().GetString("homepage")
			if err != nil {
				return err
			}
			body.Homepage = &homepage
		}

		res, err := client.PostV1OrgOrgNameSettingsApplicationsOauthWithResponse(context.TODO(), orgName, body, func(_ context.Context, req *http.Request) error {
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

var settingsApplicationsOauthSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Search oAuthApps of the active organization",
	Long:  "Required role: ADMIN\nRequired scopes: IDP_ORGANIZATION_ADMIN",
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

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("orgName"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &idpClient.PostV1OrgOrgNameSettingsApplicationsOauthSearchParams{}

		if cmd.Flags().Changed("limit") {
			limit, err := cmd.Flags().GetFloat32("limit")
			if err != nil {
				return err
			}
			params.Limit = &limit
		}

		if cmd.Flags().Changed("page") {
			page, err := cmd.Flags().GetFloat32("page")
			if err != nil {
				return err
			}
			params.Page = &page
		}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body idpClient.PostV1OrgOrgNameSettingsApplicationsOauthSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameSettingsApplicationsOauthSearchWithResponse(context.TODO(), orgName, params, body, func(_ context.Context, req *http.Request) error {
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

var settingsApplicationsOauthDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete OAuthApp",
	Long:  "Required role: ADMIN\nRequired scopes: IDP_ORGANIZATION_ADMIN",
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

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("orgName"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		oAuthAppId, err := command.GetArg(cmd, command.StringFlagIfSet("oAuthAppId"), command.ContextIfValidString(latest.OauthAppID), command.PromptString("oAuthAppId"))
		if err != nil {
			return err
		}
		latest.OauthAppID = sql.NullString{Valid: true, String: oAuthAppId}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1OrgOrgNameSettingsApplicationsOauthOAuthAppIdWithResponse(context.TODO(), orgName, oAuthAppId, func(_ context.Context, req *http.Request) error {
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

var settingsApplicationsOauthUpdateCommand = &cobra.Command{
	Use:   "update",
	Short: "Update an oAuthApp by its ID",
	Long:  "Required role: ADMIN\nRequired scopes: IDP_ORGANIZATION_ADMIN",
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

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("orgName"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		oAuthAppId, err := command.GetArg(cmd, command.StringFlagIfSet("oAuthAppId"), command.ContextIfValidString(latest.OauthAppID), command.PromptString("oAuthAppId"))
		if err != nil {
			return err
		}
		latest.OauthAppID = sql.NullString{Valid: true, String: oAuthAppId}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		callback, err := cmd.Flags().GetStringSlice("callback")
		if err != nil {
			return err
		}

		name, err := command.FlagOrPromptString(cmd, "name")
		if err != nil {
			return err
		}

		body := idpClient.PutV1OrgOrgNameSettingsApplicationsOauthOAuthAppIdJSONRequestBody{
			Callback: callback,
			Name:     name,
		}

		if cmd.Flags().Changed("allowAuthorizationCodeFlow") {
			allowAuthorizationCodeFlow, err := cmd.Flags().GetBool("allowAuthorizationCodeFlow")
			if err != nil {
				return err
			}
			body.AllowAuthorizationCodeFlow = &allowAuthorizationCodeFlow
		}

		if cmd.Flags().Changed("allowDeviceFlow") {
			allowDeviceFlow, err := cmd.Flags().GetBool("allowDeviceFlow")
			if err != nil {
				return err
			}
			body.AllowDeviceFlow = &allowDeviceFlow
		}

		if cmd.Flags().Changed("description") {
			description, err := cmd.Flags().GetString("description")
			if err != nil {
				return err
			}
			body.Description = &description
		}

		if cmd.Flags().Changed("homepage") {
			homepage, err := cmd.Flags().GetString("homepage")
			if err != nil {
				return err
			}
			body.Homepage = &homepage
		}

		res, err := client.PutV1OrgOrgNameSettingsApplicationsOauthOAuthAppIdWithResponse(context.TODO(), orgName, oAuthAppId, body, func(_ context.Context, req *http.Request) error {
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

var settingsApplicationsOauthSecretCommand = &cobra.Command{
	Use:   "secret",
	Short: "Creates a new secret inside an existing oAuthApp",
	Long:  "Required role: ADMIN\nRequired scopes: IDP_ORGANIZATION_ADMIN",
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

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("orgName"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		oAuthAppId, err := command.GetArg(cmd, command.StringFlagIfSet("oAuthAppId"), command.ContextIfValidString(latest.OauthAppID), command.PromptString("oAuthAppId"))
		if err != nil {
			return err
		}
		latest.OauthAppID = sql.NullString{Valid: true, String: oAuthAppId}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		secretName, err := command.FlagOrPromptString(cmd, "secretName")
		if err != nil {
			return err
		}

		body := idpClient.PostV1OrgOrgNameSettingsApplicationsOauthOAuthAppIdSecretJSONRequestBody{
			SecretName: secretName,
		}

		res, err := client.PostV1OrgOrgNameSettingsApplicationsOauthOAuthAppIdSecretWithResponse(context.TODO(), orgName, oAuthAppId, body, func(_ context.Context, req *http.Request) error {
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
