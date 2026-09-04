package idp

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/idpClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	loginOauthInfoCommand.Flags().String("clientId", "", "The ID of the application that is asking for authorization")
	loginOauthInfoCommand.Flags().String("redirectURI", "", "A successful response will redirect to this URI")
	loginOauthInfoCommand.Flags().String("scope", "", "Space-delimited list of hcloud permissions that the application requires as well as openId if required")
	loginOauthInfoCommand.Flags().String("responseType", "code", "Tells the authorization server which grant to execute. Currently only \"code\" is supported")
	loginOauthInfoCommand.Flags().String("state", "", "Optional value, used for security purposes (XSRF) or to store non-sensitive information between requests. If this request parameter is set in the request, then it is returned to the application as part of the redirect_uri with no changes")
	loginOauthInfoCommand.Flags().String("codeChallenge", "", "Optional value, used for PKCE (https://www.rfc-editor.org/rfc/rfc7636#section-4)")
	loginOauthInfoCommand.Flags().String("codeChallengeMethod", "", "Optional value, used for PKCE (https://www.rfc-editor.org/rfc/rfc7636#section-4)")
	loginOauthDeviceAuthorizeCommand.Flags().String("clientId", "", "Client ID of the OAuth application")
	loginOauthDeviceAuthorizeCommand.Flags().String("userCode", "", "User code to validate")
	loginOauthDeviceCommand.AddCommand(loginOauthDeviceAuthorizeCommand)
	loginOauthAuthorizeCommand.Flags().String("clientId", "", "The ID of the application that is asking for authorization")
	loginOauthAuthorizeCommand.Flags().String("redirectURI", "", "A successful response will redirect to this URI")
	loginOauthAuthorizeCommand.Flags().String("scope", "", "Space-delimited list of hcloud permissions that the application requires as well as openId if required")
	loginOauthAuthorizeCommand.Flags().String("response_type", "code", "Tells the authorization server which grant to execute. Currently only \"code\" is supported")
	loginOauthAuthorizeCommand.Flags().String("state", "", "Optional value, used for security purposes (XSRF) or to store non-sensitive information between requests. If this request parameter is set in the request, then it is returned to the application as part of the redirect_uri with no changes")
	loginOauthAuthorizeCommand.Flags().String("codeChallenge", "", "Optional value, used for PKCE (https://www.rfc-editor.org/rfc/rfc7636#section-4)")
	loginOauthAuthorizeCommand.Flags().String("codeChallengeMethod", "", "Optional value, used for PKCE (https://www.rfc-editor.org/rfc/rfc7636#section-4)")
	oauthCommand.AddCommand(oauthWellKnownCommand)
}

var oauthCommand = &cobra.Command{
	Use:   "oauth",
	Short: "Commands for oauth",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var oauthWellKnownCommand = &cobra.Command{
	Use: "well-known",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
		if err != nil {
			return err
		}

		res, err := client.GetV1WellKnownOauthAuthorizationServerWithResponse(context.TODO())
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

var loginOauthAuthorizeCommand = &cobra.Command{
	Use:   "authorize",
	Short: "OAuth Authorization Endpoint",
	Long:  "Request access code for the user\nRequired scopes: HCLOUD_FULL",
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
		latest.OauthAppClientID = sql.NullString{Valid: true, String: clientId}

		redirectUri, err := command.GetArg(cmd, command.StringFlagIfSet("redirectURI"), command.ContextIfValidString(latest.OauthAppRedirectUri), command.PromptString("OAuth App redirect URI"))
		if err != nil {
			return err
		}
		latest.OauthAppRedirectUri = sql.NullString{Valid: true, String: redirectUri}

		scope, err := command.FlagOrPromptString(cmd, "scope")
		if err != nil {
			return err
		}

		params := &idpClient.GetV1LoginOauthAuthorizeParams{
			ClientId:    clientId,
			RedirectUri: redirectUri,
			Scope:       scope,
		}

		responseType, err := cmd.Flags().GetString("response_type")
		if err != nil {
			return err
		}
		params.ResponseType = &responseType

		if cmd.Flags().Changed("state") {
			state, err := cmd.Flags().GetString("state")
			if err != nil {
				return err
			}
			params.State = &state
		}

		if cmd.Flags().Changed("codeChallenge") {
			codeChallenge, err := cmd.Flags().GetString("codeChallenge")
			if err != nil {
				return err
			}
			params.CodeChallenge = &codeChallenge
		}

		if cmd.Flags().Changed("codeChallengeMethod") {
			codeChallengeMethod, err := cmd.Flags().GetString("codeChallengeMethod")
			if err != nil {
				return err
			}
			codeChallengeMethodValue := idpClient.GetV1LoginOauthAuthorizeParamsCodeChallengeMethod(codeChallengeMethod)
			params.CodeChallengeMethod = &codeChallengeMethodValue
		}

		res, err := client.GetV1LoginOauthAuthorizeWithResponse(context.TODO(), params)
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

var loginOauthDeviceCommand = &cobra.Command{
	Use:   "device",
	Short: "Commands for device",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var loginOauthDeviceAuthorizeCommand = &cobra.Command{
	Use:   "authorize",
	Short: "Device flow - Validate the user code",
	Long:  "Authorization endpoint used in the OAuth 2.0 Device Authorization Flow to validate the user code",
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

		clientId, err := command.GetArg(cmd, command.StringFlagIfSet("clientId"), command.ContextIfValidString(latest.OauthAppClientID), command.PromptString("OAuth App Client ID"))
		if err != nil {
			return err
		}
		latest.OauthAppClientID = sql.NullString{Valid: true, String: clientId}

		userCode, err := command.FlagOrDifferentPromptString(cmd, "userCode", "user code")
		if err != nil {
			return err
		}

		res, err := client.GetV1LoginOauthDeviceAuthorizeWithResponse(context.TODO(), &idpClient.GetV1LoginOauthDeviceAuthorizeParams{
			ClientId: clientId,
			UserCode: userCode,
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

var loginOauthInfoCommand = &cobra.Command{
	Use:   "info",
	Short: "OAuth Preflight Info Endpoint",
	Long:  "This endpoint servers as a pseudo **authorize endpoint**. It will return an error if the user has not given permission to it yet or it will redirect like the **authorize endpoint** if the user has already consented to the scopes specified",
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

		clientId, err := command.GetArg(cmd, command.StringFlagIfSet("clientId"), command.ContextIfValidString(latest.OauthAppClientID), command.PromptString("OAuth client ID"))
		if err != nil {
			return err
		}
		latest.OauthAppClientID = sql.NullString{Valid: true, String: clientId}

		redirectUri, err := command.GetArg(cmd, command.StringFlagIfSet("redirectURI"), command.ContextIfValidString(latest.OauthAppRedirectUri), command.PromptString("OAuth app redirect URI"))
		if err != nil {
			return err
		}
		latest.OauthAppRedirectUri = sql.NullString{Valid: true, String: redirectUri}

		scope, err := command.FlagOrPromptString(cmd, "scope")
		if err != nil {
			return err
		}

		params := &idpClient.GetV1LoginOauthInfoParams{
			ClientId:    clientId,
			RedirectUri: redirectUri,
			Scope:       scope,
		}

		responseType, err := cmd.Flags().GetString("responseType")
		if err != nil {
			return err
		}
		params.ResponseType = &responseType

		if cmd.Flags().Changed("state") {
			state, err := cmd.Flags().GetString("state")
			if err != nil {
				return err
			}
			params.State = &state
		}

		if cmd.Flags().Changed("codeChallenge") {
			codeChallenge, err := cmd.Flags().GetString("codeChallenge")
			if err != nil {
				return err
			}
			params.CodeChallenge = &codeChallenge
		}

		if cmd.Flags().Changed("codeChallengeMethod") {
			codeChallengeMethod, err := cmd.Flags().GetString("codeChallengeMethod")
			if err != nil {
				return err
			}
			codeChallengeMethodValue := idpClient.GetV1LoginOauthInfoParamsCodeChallengeMethod(codeChallengeMethod)
			params.CodeChallengeMethod = &codeChallengeMethodValue
		}

		res, err := client.GetV1LoginOauthInfoWithResponse(context.TODO(), params)
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
