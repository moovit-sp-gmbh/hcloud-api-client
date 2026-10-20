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
	avatarDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	avatarDeleteCommand.Flags().String("token", "", "authorization token")
	avatarCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	avatarCreateCommand.Flags().String("token", "", "authorization token")

	avatarCommand.AddCommand(avatarDeleteCommand)
	avatarCommand.AddCommand(avatarCreateCommand)
	avatarCommand.AddCommand(avatarApplicationsCommand)
	avatarCommand.AddCommand(avatarServiceAccountsCommand)
	avatarCommand.AddCommand(avatarSpacesCommand)
	avatarCommand.AddCommand(avatarStoragesCommand)
	avatarCommand.AddCommand(avatarTeamsCommand)
	avatarCommand.AddCommand(avatarUserCommand)
}

var avatarCommand = &cobra.Command{
	Use:   "avatar",
	Short: "Commands for avatar",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var avatarDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete an Avatar for organization",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1AvatarOrgOrgNameWithResponse(context.TODO(), orgName, &daliClient.DeleteV1AvatarOrgOrgNameParams{
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

var avatarCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new Avatar for organization",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.PostV1AvatarOrgOrgNameWithResponse(context.TODO(), orgName, &daliClient.PostV1AvatarOrgOrgNameParams{
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
