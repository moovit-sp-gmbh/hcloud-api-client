package idp

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/idpClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	organizationGetCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	organizationGetCommand.Flags().String("token", "", "authorization token")
	organizationGetCommand.Flags().Bool("teamsOfUser", false, "If true, an additional property holding all teams that the requesting user is part of and that belong to the organization will be returned")
	organizationGetCommand.Flags().Bool("totalMemberCount", false, "If true, an additional property holding the total member count of the organization will be returned")
	organizationGetCommand.Flags().Float32("membersSample", 0, "If set, an additional property holding a sample of members will be returned. The provided number represents the max. size of the sample")
	organizationGetCommand.Flags().Bool("license", false, "If set, an additional property holding the organization's license will be included")
	orgCommand.Flags().StringP("name", "n", "", "Name of the organization")

	createOrgCommand.Flags().String("company", "", "Company name")

	orgCommand.AddCommand(createOrgCommand)
	orgCommand.AddCommand(organizationGetCommand)
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

var organizationGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Retrieves an organization by it's name",
	Long:  "\nRequired role: MEMBER\nRequired scopes: IDP_ORGANIZATION_READ",
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

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("organization name"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &idpClient.GetV1OrgOrgNameParams{}

		if cmd.Flags().Changed("teamsOfUser") {
			teamsOfUser, err := cmd.Flags().GetBool("teamsOfUser")
			if err != nil {
				return err
			}
			params.TeamsOfUser = &teamsOfUser
		}

		if cmd.Flags().Changed("totalMemberCount") {
			totalMemberCount, err := cmd.Flags().GetBool("totalMemberCount")
			if err != nil {
				return err
			}
			params.TotalMemberCount = &totalMemberCount
		}

		if cmd.Flags().Changed("membersSample") {
			membersSample, err := cmd.Flags().GetFloat32("membersSample")
			if err != nil {
				return err
			}
			params.MembersSample = &membersSample
		}

		if cmd.Flags().Changed("license") {
			license, err := cmd.Flags().GetBool("license")
			if err != nil {
				return err
			}
			params.License = &license
		}

		res, err := client.GetV1OrgOrgNameWithResponse(context.TODO(), orgName, params, func(_ context.Context, req *http.Request) error {
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
