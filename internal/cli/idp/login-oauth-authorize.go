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
	loginOauthAuthorizeCreateCommand.Flags().String("token", "", "authorization token")
	loginOauthAuthorizeCreateCommand.Flags().String("client_id", "", "The ID of the application that is asking for authorization")
	loginOauthAuthorizeCreateCommand.Flags().String("redirect_uri", "", "A successful response will redirect to this URI")
	loginOauthAuthorizeCreateCommand.Flags().String("scope", "", "Space-delimited list of hcloud permissions that the application requires as well as openId if required")
	loginOauthAuthorizeCreateCommand.Flags().String("response_type", "", "Tells the authorization server which grant to execute. Currently only \"code\" is supported")
	loginOauthAuthorizeCreateCommand.Flags().String("state", "", "Optional value, used for security purposes (XSRF) or to store non-sensitive information between requests. If this request parameter is set in the request, then it is returned to the application as part of the redirect_uri with no changes")
	loginOauthAuthorizeCreateCommand.Flags().String("code_challenge", "", "Optional value, used for PKCE (https://www.rfc-editor.org/rfc/rfc7636#section-4)")
	loginOauthAuthorizeCreateCommand.Flags().String("code_challenge_method", "", "Optional value, used for PKCE (https://www.rfc-editor.org/rfc/rfc7636#section-4)")
}

var loginOauthAuthorizeCreateCommand = &cobra.Command{
	Use:   "create",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		clientId, err := command.FlagOrPromptString(cmd, "client_id")
		if err != nil {
			return err
		}

		redirectUri, err := command.FlagOrPromptString(cmd, "redirect_uri")
		if err != nil {
			return err
		}

		scope, err := command.FlagOrPromptString(cmd, "scope")
		if err != nil {
			return err
		}

		params := &idpClient.PostV1LoginOauthAuthorizeParams{
			ClientId:    clientId,
			RedirectUri: redirectUri,
			Scope:       scope,
		}

		if cmd.Flags().Changed("response_type") {
			responseType, err := cmd.Flags().GetString("response_type")
			if err != nil {
				return err
			}
			params.ResponseType = &responseType
		}

		if cmd.Flags().Changed("state") {
			state, err := cmd.Flags().GetString("state")
			if err != nil {
				return err
			}
			params.State = &state
		}

		if cmd.Flags().Changed("code_challenge") {
			codeChallenge, err := cmd.Flags().GetString("code_challenge")
			if err != nil {
				return err
			}
			params.CodeChallenge = &codeChallenge
		}

		if cmd.Flags().Changed("code_challenge_method") {
			codeChallengeMethod, err := cmd.Flags().GetString("code_challenge_method")
			if err != nil {
				return err
			}
			codeChallengeMethodValue := idpClient.PostV1LoginOauthAuthorizeParamsCodeChallengeMethod(codeChallengeMethod)
			params.CodeChallengeMethod = &codeChallengeMethodValue
		}

		res, err := client.PostV1LoginOauthAuthorizeWithResponse(context.TODO(), params, func(_ context.Context, req *http.Request) error {
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
