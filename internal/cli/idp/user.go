package idp

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/idpClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	patchCommand.Flags().String("name", "", "New user name")

	whoamiCommand.AddCommand(patchCommand)
}

var whoamiCommand = &cobra.Command{
	Use:   "user",
	Short: "Get user info",
	Long:  "Get user info",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.Server(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server + "/api/account")
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		res, err := client.GetV1UserWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}

		text, err := command.JSON(res.JSON200)
		if err != nil {
			return err
		}
		command.SetResult(cmd, text)

		return nil
	},
}

var patchCommand = &cobra.Command{
	Use:   "patch",
	Short: "Patch user info",
	Long:  "Patch user info",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.Server(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server + "/api/account")
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		newName, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		res, err := client.PatchV1UserWithResponse(context.TODO(), idpClient.PatchV1UserJSONRequestBody{
			Name: &newName,
		}, func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}
		log.Println(res.JSON200.Name)

		ctx.NewName = sql.NullString{Valid: true, String: res.JSON200.Name}
		err = service.Push(ctx)
		if err != nil {
			log.Println(err.Error())
		}

		text, err := command.JSON(res.JSON200)
		if err != nil {
			return err
		}
		command.SetResult(cmd, text)

		return nil
	},
}
