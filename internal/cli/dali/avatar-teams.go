package dali

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/daliClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	avatarTeamsDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	avatarTeamsDeleteCommand.Flags().String("teamName", "", "Name of the team (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	avatarTeamsDeleteCommand.Flags().String("token", "", "authorization token")
	avatarTeamsCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	avatarTeamsCreateCommand.Flags().String("teamName", "", "Name of the team (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	avatarTeamsCreateCommand.Flags().String("token", "", "authorization token")

	avatarTeamsCommand.AddCommand(avatarTeamsDeleteCommand)
	avatarTeamsCommand.AddCommand(avatarTeamsCreateCommand)
}

var avatarTeamsCommand = &cobra.Command{
	Use:   "teams",
	Short: "Commands for teams",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var avatarTeamsDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete avatar for the team",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewDaliClient(cmd)
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

		res, err := client.DeleteV1AvatarOrgOrgNameTeamsTeamNameWithResponse(context.TODO(), orgName, teamName, &daliClient.DeleteV1AvatarOrgOrgNameTeamsTeamNameParams{
			Authorization: jwt,
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

var avatarTeamsCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new avatar for the team",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewDaliClient(cmd)
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

		res, err := client.PostV1AvatarOrgOrgNameTeamsTeamNameWithResponse(context.TODO(), orgName, teamName, &daliClient.PostV1AvatarOrgOrgNameTeamsTeamNameParams{
			Authorization: jwt,
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
