package idp

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/idpClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/oapi-codegen/runtime/types"
	"github.com/spf13/cobra"
)

func init() {
	loginSamlCommand.AddCommand(loginSamlMetadataCommand)
	loginSamlCommand.Flags().String("origin", "", "A URL marking the starting point of the login process. At the end the user will be redirected back to this URL with token=<JWT> added to the query string")
	loginSamlCommand.Flags().String("email", "", "The user's email. The email domain will be used to determine the appropriate SAML 2.0 provider to use going forward")
}

var loginSamlCommand = &cobra.Command{
	Use:   "saml",
	Short: "Start the login process using the SAML2.0 protocol",
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

		origin, err := command.FlagOrPromptString(cmd, "origin")
		if err != nil {
			return err
		}

		email, err := command.GetArg(cmd, command.StringFlagIfSet("email"), command.ContextIfValidString(latest.Email), command.PromptString("email"))
		if err != nil {
			return err
		}
		latest.Email = sql.NullString{Valid: true, String: email}

		res, err := client.GetV1LoginSamlWithResponse(context.TODO(), &idpClient.GetV1LoginSamlParams{
			Origin: origin,
			Email:  types.Email(email),
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

var loginSamlMetadataCommand = &cobra.Command{
	Use:   "metadata",
	Short: "Obtain the metadata XML file with information about HCloud as a SAML 2.0 Service Provider",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
		if err != nil {
			return err
		}

		res, err := client.GetV1LoginSamlMetadataXmlWithResponse(context.TODO())
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
