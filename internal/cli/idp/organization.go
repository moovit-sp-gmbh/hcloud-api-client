package idp

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/idpClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	orgCommand.Flags().StringP("name", "n", "", "Name of the organization")

	createOrgCommand.Flags().String("company", "", "Company name")

	orgCommand.AddCommand(createOrgCommand)
}

var orgCommand = &cobra.Command{
	Use:   "organization",
	Short: "Organization entrypoint",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

var createOrgCommand = &cobra.Command{
	Use:   "create <name>",
	Short: "Create an organization",
	Long:  "Create an organization with a given name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}
		orgName := args[0]
		company, err := cmd.Flags().GetString("company")
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

		body := map[string]string{
			"name": orgName,
		}

		if company != "" {
			body["company"] = company
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		var res *idpClient.PostV1OrgResponse
		if company == "" {
			res, err = client.PostV1OrgWithResponse(context.TODO(), nil, idpClient.OrganizationCreateDto{
				Name: orgName,
			}, func(_ context.Context, req *http.Request) error {
				req.Header.Set("authorization", ctx.Jwt.String)
				return nil
			})
		} else {
			res, err = client.PostV1OrgWithResponse(context.TODO(), nil, idpClient.OrganizationCreateDto{
				Name:    orgName,
				Company: &company,
			}, func(_ context.Context, req *http.Request) error {
				req.Header.Set("authorization", ctx.Jwt.String)
				return nil
			})
		}

		var text command.Text
		switch res.StatusCode() {
		case 200, 201:
			text, err = command.JSON(res.JSON201)
			if err != nil {
				return err
			}

			ctx.OrgName = sql.NullString{Valid: true, String: orgName}
			service.Push(ctx)
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
