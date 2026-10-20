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
	teamsCommand.AddCommand(teamsCreateCommand)
	teamsCommand.AddCommand(teamsSearchCommand)
	teamsCommand.AddCommand(teamsDeleteCommand)
	teamsCommand.AddCommand(teamsPatchCommand)
	teamsCommand.AddCommand(teamsMembersCommand)
	organizationTeamsGetCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	organizationTeamsGetCommand.Flags().String("teamName", "", "Name of the team (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	organizationTeamsGetCommand.Flags().String("token", "", "authorization token")
	organizationTeamsGetCommand.Flags().Float32("membersLimit", 0, "Max number ob members that shall be returned in the team object. If not defined, no members will be returned")
	teamsCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	teamsCreateCommand.Flags().String("token", "", "authorization token")
	teamsCreateCommand.Flags().Bool("noAvatar", false, "")
	teamsCreateCommand.Flags().String("name", "", "Name of the team (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	teamsCreateCommand.Flags().StringSlice("userIds", nil, "Array of user database IDs that should be part of the team")
	teamsSearchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	teamsSearchCommand.Flags().String("token", "", "authorization token")
	teamsSearchCommand.Flags().Float32("limit", 0, "Page size")
	teamsSearchCommand.Flags().Float32("page", 0, "Page index")
	teamsSearchCommand.Flags().Bool("totalMemberCount", false, "If true, an additional property holding the total member count of the organization will be returned")
	teamsSearchCommand.Flags().Float32("membersSample", 0, "If set, an additional property holding a sample of members will be returned. The provided number represents the max size of the sample")
	teamsSearchCommand.Flags().String("body", "{}", "request body as JSON")
	teamsDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	teamsDeleteCommand.Flags().String("teamName", "", "Name of the team (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	teamsDeleteCommand.Flags().String("token", "", "authorization token")
	teamsPatchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	teamsPatchCommand.Flags().String("teamName", "", "Name of the team (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	teamsPatchCommand.Flags().String("token", "", "authorization token")
	teamsPatchCommand.Flags().String("name", "", "Name of the team (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	teamsPatchCommand.Flags().StringSlice("userIds", nil, "Array of user database IDs that should be part of the team")
	teamsPatchCommand.Flags().String("usersOperation", "", "Define which operation should be performed with the provided users array on the existing users array")
	organizationTeamsCommand.AddCommand(organizationTeamsGetCommand)
}

var organizationTeamsCommand = &cobra.Command{
	Use:   "teams",
	Short: "Commands for teams",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var organizationTeamsGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Get a team by its name",
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

		teamName, err := command.GetArg(cmd, command.StringFlagIfSet("teamName"), command.ContextIfValidString(latest.TeamName), command.PromptString("team name"))
		if err != nil {
			return err
		}
		latest.TeamName = sql.NullString{Valid: true, String: teamName}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &idpClient.GetV1OrgOrgNameTeamsTeamNameParams{}

		if cmd.Flags().Changed("membersLimit") {
			membersLimit, err := cmd.Flags().GetFloat32("membersLimit")
			if err != nil {
				return err
			}
			params.MembersLimit = &membersLimit
		}

		res, err := client.GetV1OrgOrgNameTeamsTeamNameWithResponse(context.TODO(), orgName, teamName, params, func(_ context.Context, req *http.Request) error {
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

var teamsCommand = &cobra.Command{
	Use:   "teams",
	Short: "Commands for teams",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var teamsCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new team",
	Long:  "Required role: ADMIN\nRequired scopes: IDP_ORGANIZATION_WRITE",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &idpClient.PostV1OrgOrgNameTeamsParams{}

		if cmd.Flags().Changed("noAvatar") {
			noAvatar, err := cmd.Flags().GetBool("noAvatar")
			if err != nil {
				return err
			}
			params.NoAvatar = &noAvatar
		}

		name, err := command.FlagOrPromptString(cmd, "name")
		if err != nil {
			return err
		}

		userIds, err := cmd.Flags().GetStringSlice("userIds")
		if err != nil {
			return err
		}

		body := idpClient.PostV1OrgOrgNameTeamsJSONRequestBody{
			Name:    name,
			UserIds: userIds,
		}

		res, err := client.PostV1OrgOrgNameTeamsWithResponse(context.TODO(), orgName, params, body, func(_ context.Context, req *http.Request) error {
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

var teamsSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Search all teams of organization by organization name and the supplied filters",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &idpClient.PostV1OrgOrgNameTeamsSearchParams{}

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

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body idpClient.PostV1OrgOrgNameTeamsSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameTeamsSearchWithResponse(context.TODO(), orgName, params, body, func(_ context.Context, req *http.Request) error {
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

var teamsDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete Team",
	Long:  "Required role: ADMIN\nRequired scopes: IDP_ORGANIZATION_WRITE",
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

		res, err := client.DeleteV1OrgOrgNameTeamsTeamNameWithResponse(context.TODO(), orgName, teamName, func(_ context.Context, req *http.Request) error {
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

var teamsPatchCommand = &cobra.Command{
	Use:   "patch",
	Short: "Patch a team by its name",
	Long:  "Required role: MANAGER\nRequired scopes: IDP_ORGANIZATION_WRITE",
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

		body := idpClient.PatchV1OrgOrgNameTeamsTeamNameJSONRequestBody{}

		if cmd.Flags().Changed("name") {
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}
			body.Name = &name
		}

		if cmd.Flags().Changed("userIds") {
			userIds, err := cmd.Flags().GetStringSlice("userIds")
			if err != nil {
				return err
			}
			body.UserIds = &userIds
		}

		if cmd.Flags().Changed("usersOperation") {
			usersOperation, err := cmd.Flags().GetString("usersOperation")
			if err != nil {
				return err
			}
			usersOperationValue := idpClient.UsersOperation(usersOperation)
			body.UsersOperation = &usersOperationValue
		}

		res, err := client.PatchV1OrgOrgNameTeamsTeamNameWithResponse(context.TODO(), orgName, teamName, body, func(_ context.Context, req *http.Request) error {
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
