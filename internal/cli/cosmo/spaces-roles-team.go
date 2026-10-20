package cosmo

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/cosmoClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	spacesRolesTeamUpdateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesRolesTeamUpdateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesRolesTeamUpdateCommand.Flags().String("roleId", "", "Role ID")
	spacesRolesTeamUpdateCommand.Flags().String("token", "", "authorization token")
	spacesRolesTeamUpdateCommand.Flags().String("teamName", "", "Name of the team")
	spacesRolesTeamDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesRolesTeamDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesRolesTeamDeleteCommand.Flags().String("roleId", "", "Role ID")
	spacesRolesTeamDeleteCommand.Flags().String("teamName", "", "Team Name")
	spacesRolesTeamDeleteCommand.Flags().String("token", "", "authorization token")

	spacesRolesTeamCommand.AddCommand(spacesRolesTeamUpdateCommand)
	spacesRolesTeamCommand.AddCommand(spacesRolesTeamDeleteCommand)
}

var spacesRolesTeamCommand = &cobra.Command{
	Use:   "team",
	Short: "Commands for team",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesRolesTeamUpdateCommand = &cobra.Command{
	Use:   "update",
	Short: "Assign a team to a custom role",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewCosmoClient(cmd)
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

		roleId, err := command.FlagOrPromptString(cmd, "roleId")
		if err != nil {
			return err
		}
		roleIdParsed, err := uuid.Parse(roleId)
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		teamName, err := command.GetArg(cmd, command.StringFlagIfSet("teamName"), command.ContextIfValidString(latest.TeamName), command.PromptString("teamName"))
		if err != nil {
			return err
		}

		body := cosmoClient.PutV1OrgOrgNameSpacesSpaceNameRolesRoleIdTeamJSONRequestBody{
			TeamName: teamName,
		}

		res, err := client.PutV1OrgOrgNameSpacesSpaceNameRolesRoleIdTeamWithResponse(context.TODO(), orgName, spaceName, roleIdParsed, body, func(_ context.Context, req *http.Request) error {
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

var spacesRolesTeamDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Unassign a team from a custom role",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewCosmoClient(cmd)
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

		roleId, err := command.FlagOrPromptString(cmd, "roleId")
		if err != nil {
			return err
		}
		roleIdParsed, err := uuid.Parse(roleId)
		if err != nil {
			return err
		}

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

		res, err := client.DeleteV1OrgOrgNameSpacesSpaceNameRolesRoleIdTeamTeamNameWithResponse(context.TODO(), orgName, spaceName, roleIdParsed, teamName, func(_ context.Context, req *http.Request) error {
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
