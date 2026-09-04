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
	loginPreCommand.Flags().String("email", "", "Email of the user")
	loginPreCommand.Flags().String("origin", "", "Origin of the request. Used when the response path is EXTERNAL")
	loginPreCommand.Flags().String("redirect", "false", "Whether or not redirect in the response. If true the response status code can be 302")
}

var loginPreCommand = &cobra.Command{
	Use:   "pre",
	Short: "Endpoint to get the way to proceed with the authentication",
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
		latest.Email = sql.NullString{Valid: true, String: email}

		origin, err := command.FlagOrPromptString(cmd, "origin")
		if err != nil {
			return err
		}

		params := &idpClient.GetV1LoginPreParams{
			Email:  types.Email(email),
			Origin: origin,
		}

		redirect, err := cmd.Flags().GetString("redirect")
		if err != nil {
			return err
		}
		redirectValue := idpClient.GetV1LoginPreParamsRedirect(redirect)
		params.Redirect = &redirectValue

		res, err := client.GetV1LoginPreWithResponse(context.TODO(), params)
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
