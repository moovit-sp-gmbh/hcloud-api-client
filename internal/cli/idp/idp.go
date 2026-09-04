package idp

import (
	"context"
	"database/sql"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/idpClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/oapi-codegen/runtime/types"
	"github.com/spf13/cobra"
)

func init() {
	loginCommand.AddCommand(loginOauthCommand)
	loginCommand.Flags().StringP("email", "e", "", "email")
	loginCommand.Flags().StringP("password", "p", "", "password")
}

var loginCommand = &cobra.Command{
	Use:   "login",
	Short: "Authenticate against helmut.cloud",
	Long:  "Obtain an authorization token to use in future requests. Saved in context as JWT",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		email, err := command.FlagOrPromptString(cmd, "email")
		if err != nil {
			return err
		}

		password, err := command.FlagOrPromptString(cmd, "password")
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		res, err := client.PostV1LoginWithResponse(context.TODO(), nil, idpClient.PostV1LoginJSONRequestBody{
			Email:    types.Email(email),
			Password: password,
		})
		if err != nil {
			return err
		}

		var text command.Text
		switch res.StatusCode() {
		case 200:
			authorization := res.HTTPResponse.Header.Get("authorization")

			service := internal.New()
			ctx, err := service.GetLatest()
			if err != nil {
				return err
			}
			ctx.Jwt = sql.NullString{Valid: true, String: authorization}
			err = service.Push(ctx)
			if err != nil {
				return err
			}

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

func Command() *cobra.Command {
	idpCommand := &cobra.Command{
		Use:   "idp",
		Short: "Entrypoint to IDP commands",
		Long:  `Entrypoint to helmut.cloud Identity Provider's commands`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	idpCommand.PersistentFlags().String("server", "https://eu-central-1.app.helmut.cloud", "Set the target server/cluster")

	idpCommand.AddCommand(loginCommand)
	idpCommand.AddCommand(orgCommand)
	idpCommand.AddCommand(whoamiCommand)
	idpCommand.AddCommand(oauthCommand)
	idpCommand.AddCommand(guestCommand)

	return idpCommand
}
