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
	loginSamlAssertCommand.Flags().String("RelayState", "", "State value forwarded by the SAML IdP")
	loginSamlAssertCommand.Flags().String("SAMLResponse", "", "Base64 URL encoded XML string containing the response as per the <a href=\"https://docs.oasis-open.org/security/saml/v2.0/saml-bindings-2.0-os.pdf\">SAML 2.0 HTTP POST binding (Section 3.5)</a>")
}

var loginSamlAssertCommand = &cobra.Command{
	Use:   "assert",
	Short: "Callback URL for the login process using the SAML2.0 protocol",
	Long:  "It is required that the HTTP Post SAML binding be used in this operation. <a href=\"https://docs.oasis-open.org/security/saml/v2.0/saml-bindings-2.0-os.pdf\">SAML 2.0 Bindings (Section 3.5)</a>",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
		if err != nil {
			return err
		}

		relayState, err := command.FlagOrPromptString(cmd, "RelayState")
		if err != nil {
			return err
		}

		sAMLResponse, err := command.FlagOrPromptString(cmd, "SAMLResponse")
		if err != nil {
			return err
		}

		body := idpClient.PostV1LoginSamlAssertJSONRequestBody{
			RelayState:   relayState,
			SAMLResponse: sAMLResponse,
		}

		res, err := client.PostV1LoginSamlAssertWithResponse(context.TODO(), body)
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
