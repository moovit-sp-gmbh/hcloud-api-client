package idp

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/idpClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	teamsMembersSearchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	teamsMembersSearchCommand.Flags().String("teamName", "", "Name of the team (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	teamsMembersSearchCommand.Flags().String("token", "", "authorization token")
	teamsMembersSearchCommand.Flags().Float32("limit", 0, "Page size")
	teamsMembersSearchCommand.Flags().Float32("page", 0, "Page index")
	teamsMembersSearchCommand.Flags().String("body", "{}", "request body as JSON")

	teamsMembersCommand.AddCommand(teamsMembersSearchCommand)
}

var teamsMembersCommand = &cobra.Command{
	Use:   "members",
	Short: "Commands for members",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var teamsMembersSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Search for members of a team",
	Long:  "Required role: MEMBER\nRequired scopes: IDP_ORGANIZATION_READ",
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

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("orgName"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		teamName, err := command.GetArg(cmd, command.StringFlagIfSet("teamName"), command.ContextIfValidString(latest.TeamName), command.PromptString("teamName"))
		if err != nil {
			return err
		}
		latest.TeamName = sql.NullString{Valid: true, String: teamName}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &idpClient.PostV1OrgOrgNameTeamsTeamNameMembersSearchParams{}

		if cmd.Flags().Changed("limit") {
			limit, err := cmd.Flags().GetFloat32("limit")
			if err != nil {
				return err
			}
			params.Limit = &limit
		}

		if cmd.Flags().Changed("page") {
			page, err := cmd.Flags().GetFloat32("page")
			if err != nil {
				return err
			}
			params.Page = &page
		}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body idpClient.PostV1OrgOrgNameTeamsTeamNameMembersSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameTeamsTeamNameMembersSearchWithResponse(context.TODO(), orgName, teamName, params, body, func(_ context.Context, req *http.Request) error {
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
