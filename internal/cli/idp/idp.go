package idp

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/cli"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/util"
	"github.com/spf13/cobra"
)

func init() {
	loginCommand.Flags().StringP("email", "e", "", "email")
	loginCommand.Flags().StringP("password", "p", "", "password")
}

var loginCommand = &cobra.Command{
	Use:   "login",
	Short: "Authenticate against helmut.cloud",
	Long:  "Obtain an authorization token to use in future requests. Saved in context as JWT",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.Server(cmd)
		if err != nil {
			return err
		}

		email, err := cmd.Flags().GetString("email")
		if err != nil {
			return err
		} else if email == "" {
			email = cli.Prompt("Email")
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			return err
		} else if password == "" {
			password = cli.Prompt("Password")
		}

		body, err := json.Marshal(struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{email, password})
		if err != nil {
			return err
		}

		url, err := util.BuildURL(server, "/api/account/v1/login")
		if err != nil {
			return err
		}
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("content-type", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		if res.StatusCode != 200 {
			text, err := io.ReadAll(res.Body)
			defer res.Body.Close()
			if err != nil {
				return err
			}
			return errors.New(string(text))
		}

		authorization := res.Header.Get("authorization")

		text, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}
		ctx.Jwt = sql.NullString{Valid: true, String: authorization}
		service.Push(ctx)

		command.SetResult(cmd, command.Text(string(text)))
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

	return idpCommand
}
