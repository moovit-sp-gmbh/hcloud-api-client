package idp

import (
	"context"
	"encoding/json"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/idpClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	loginOauthAccessTokenCommand.Flags().String("client_id", "", "HcloudClientId of the OAuth app")
	loginOauthAccessTokenCommand.Flags().String("grant_type", "", "Type of authorization grant.")
	loginOauthAccessTokenCommand.Flags().String("client_secret", "", "Client secret of the OAuth app")
	loginOauthAccessTokenCommand.Flags().String("code", "", "Code from the OAuth code grant. Will be ignored if the grant type is not authorization_code")
	loginOauthAccessTokenCommand.Flags().String("code_verifier", "", "Optional value, used for PKCE (https://www.rfc-editor.org/rfc/rfc7636#section-4)")
	loginOauthAccessTokenCommand.Flags().String("redirect_uri", "", "Callback URI of the OAuth app. If a redirect_uri was sent in the authorization phase, then it must be sent again. Optional for the authorization_code grant and ignored in the refresh_token grant")
	loginOauthAccessTokenCommand.Flags().String("refresh_token", "", "RefreshToken from the OAuth token grant. Required in the refresh_token grant and ignored in the authorization_code grant")
	loginOauthAccessTokenCommand.Flags().String("scope", "", "Scope that the access token should be granted. Optional for the refresh_token grant and ignored in the authorization_code grant. The requested scope MUST NOT include any scope not originally granted, and if omitted is treated as equal to the scope originally granted")
}

var loginOauthAccessTokenCommand = &cobra.Command{
	Use:   "access-token",
	Short: "OAuth Token Endpoint",
	Long:  "Use this endpoint to exchange a code for an access token. The code can be obtained via the authorization endpoint",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
		if err != nil {
			return err
		}

		clientId, err := command.FlagOrPromptString(cmd, "client_id")
		if err != nil {
			return err
		}

		grantType, err := command.FlagOrPromptString(cmd, "grant_type")
		if err != nil {
			return err
		}

		body := idpClient.PostV1LoginOauthAccessTokenJSONRequestBody{
			ClientId:  clientId,
			GrantType: idpClient.GrantType(grantType),
		}

		if cmd.Flags().Changed("client_secret") {
			clientSecret, err := cmd.Flags().GetString("client_secret")
			if err != nil {
				return err
			}
			body.ClientSecret = &clientSecret
		}

		if cmd.Flags().Changed("code") {
			code, err := cmd.Flags().GetString("code")
			if err != nil {
				return err
			}
			body.Code = &code
		}

		if cmd.Flags().Changed("code_verifier") {
			codeVerifier, err := cmd.Flags().GetString("code_verifier")
			if err != nil {
				return err
			}
			body.CodeVerifier = &codeVerifier
		}

		if cmd.Flags().Changed("redirect_uri") {
			redirectUri, err := cmd.Flags().GetString("redirect_uri")
			if err != nil {
				return err
			}
			body.RedirectUri = &redirectUri
		}

		if cmd.Flags().Changed("refresh_token") {
			refreshToken, err := cmd.Flags().GetString("refresh_token")
			if err != nil {
				return err
			}
			body.RefreshToken = &refreshToken
		}

		if cmd.Flags().Changed("scope") {
			scope, err := cmd.Flags().GetString("scope")
			if err != nil {
				return err
			}
			body.Scope = &scope
		}

		res, err := client.PostV1LoginOauthAccessTokenWithResponse(context.TODO(), body)
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
