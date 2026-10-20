package high5

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/high5Client"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	spacesPermissionsSearchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesPermissionsSearchCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesPermissionsSearchCommand.Flags().String("token", "", "authorization token")
	spacesPermissionsSearchCommand.Flags().Float32("limit", 0, "Page size")
	spacesPermissionsSearchCommand.Flags().Float32("page", 0, "Page index")
	spacesPermissionsSearchCommand.Flags().String("body", "{}", "request body as JSON")
	spacesPermissionsTeamCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesPermissionsTeamCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesPermissionsTeamCommand.Flags().String("token", "", "authorization token")
	spacesPermissionsTeamCommand.Flags().String("permission", "", "The permission the team has within the space Use SpacePermission.NONE to remove permissions")
	spacesPermissionsTeamCommand.Flags().String("teamName", "", "Name of the team (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesPermissionsUserCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesPermissionsUserCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesPermissionsUserCommand.Flags().String("token", "", "authorization token")
	spacesPermissionsUserCommand.Flags().String("permission", "", "")
	spacesPermissionsUserCommand.Flags().String("userId", "", "The id of the user")

	spacesPermissionsCommand.AddCommand(spacesPermissionsSearchCommand)
	spacesPermissionsCommand.AddCommand(spacesPermissionsTeamCommand)
	spacesPermissionsCommand.AddCommand(spacesPermissionsUserCommand)
}

var spacesPermissionsCommand = &cobra.Command{
	Use:   "permissions",
	Short: "Commands for permissions",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesPermissionsSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "For the specified space, will retrieve all space permissions that match the provided search filter or all space permissions if no search filter is provided",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewHigh5Client(cmd)
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

		spaceName, err := command.FlagOrPromptString(cmd, "spaceName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &high5Client.PostV1OrgOrgNameSpacesSpaceNamePermissionsSearchParams{
			Authorization: jwt,
		}

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
		var body high5Client.PostV1OrgOrgNameSpacesSpaceNamePermissionsSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNamePermissionsSearchWithResponse(context.TODO(), orgName, spaceName, params, body)
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

var spacesPermissionsTeamCommand = &cobra.Command{
	Use:   "team",
	Short: "Add, update, delete a team's permission",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewHigh5Client(cmd)
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

		spaceName, err := command.FlagOrPromptString(cmd, "spaceName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		permission, err := command.FlagOrPromptString(cmd, "permission")
		if err != nil {
			return err
		}

		teamName, err := command.GetArg(cmd, command.StringFlagIfSet("teamName"), command.ContextIfValidString(latest.TeamName), command.PromptString("teamName"))
		if err != nil {
			return err
		}

		body := high5Client.PutV1OrgOrgNameSpacesSpaceNamePermissionsTeamJSONRequestBody{
			Permission: high5Client.SpacePermission(permission),
			TeamName:   teamName,
		}

		res, err := client.PutV1OrgOrgNameSpacesSpaceNamePermissionsTeamWithResponse(context.TODO(), orgName, spaceName, &high5Client.PutV1OrgOrgNameSpacesSpaceNamePermissionsTeamParams{
			Authorization: jwt,
		}, body)
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

var spacesPermissionsUserCommand = &cobra.Command{
	Use:   "user",
	Short: "Add, update, delete a user's permission",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewHigh5Client(cmd)
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

		spaceName, err := command.FlagOrPromptString(cmd, "spaceName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		permission, err := command.FlagOrPromptString(cmd, "permission")
		if err != nil {
			return err
		}

		userId, err := command.FlagOrPromptString(cmd, "userId")
		if err != nil {
			return err
		}

		body := high5Client.PutV1OrgOrgNameSpacesSpaceNamePermissionsUserJSONRequestBody{
			Permission: high5Client.SpaceUserPermissionDtoPermission(permission),
			UserId:     userId,
		}

		res, err := client.PutV1OrgOrgNameSpacesSpaceNamePermissionsUserWithResponse(context.TODO(), orgName, spaceName, &high5Client.PutV1OrgOrgNameSpacesSpaceNamePermissionsUserParams{
			Authorization: jwt,
		}, body)
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
