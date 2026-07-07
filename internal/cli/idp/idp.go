package idp

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/util"
	"github.com/spf13/cobra"
)

var server string

func init() {
	loginCommand.Flags().StringP("data", "d", "", "<{ email, password } JSON>")
}

var loginCommand = &cobra.Command{
	Use:   "login",
	Short: "Authenticate against helmut.cloud",
	Long:  "Obtain an authorization token to use in future requests. Saved in context as JWT",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := cmd.Flags().GetString("server")
		if err != nil {
			return err
		}

		data, err := cmd.Flags().GetString("data")
		var body io.Reader
		if (err != nil || data == "") && util.IsPipe() {
			body = cmd.InOrStdin()
		} else if len(data) > 0 {
			body = strings.NewReader(data)
		} else {
			return fmt.Errorf("data argument is missing")
		}

		url, err := util.BuildURL(server, "/api/account/v1/login")
		if err != nil {
			return err
		}
		req, err := http.NewRequest(http.MethodPost, url, body)
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
		fmt.Printf("Authorization: %s\n", authorization)

		text, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return err
		}
		fmt.Println(string(text))

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}
		ctx.Jwt = sql.NullString{Valid: true, String: authorization}
		service.Push(ctx)

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

	idpCommand.PersistentFlags().StringVar(&server, "server", "https://eu-central-1.app.helmut.cloud", "Set the target server/cluster")

	idpCommand.AddCommand(loginCommand)
	idpCommand.AddCommand(orgCommand)

	return idpCommand
}
