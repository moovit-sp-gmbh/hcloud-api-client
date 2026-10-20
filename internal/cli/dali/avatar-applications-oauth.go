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
	avatarApplicationsOauthDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	avatarApplicationsOauthDeleteCommand.Flags().String("appId", "", "ID of the OAuth application (ObjectId identifier, consists of 24 hexadecimal digits)")
	avatarApplicationsOauthDeleteCommand.Flags().String("token", "", "authorization token")
	avatarApplicationsOauthCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	avatarApplicationsOauthCreateCommand.Flags().String("appId", "", "ID of the OAuth application (ObjectId identifier, consists of 24 hexadecimal digits)")
	avatarApplicationsOauthCreateCommand.Flags().String("token", "", "authorization token")

	avatarApplicationsOauthCommand.AddCommand(avatarApplicationsOauthDeleteCommand)
	avatarApplicationsOauthCommand.AddCommand(avatarApplicationsOauthCreateCommand)
	avatarApplicationsCommand.AddCommand(avatarApplicationsOauthCommand)
}

var avatarApplicationsCommand = &cobra.Command{
	Use:   "applications",
	Short: "Commands for applications",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var avatarApplicationsOauthCommand = &cobra.Command{
	Use:   "oauth",
	Short: "Commands for oauth",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var avatarApplicationsOauthDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete avatar for the OAuth application",
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

		appId, err := command.FlagOrPromptString(cmd, "appId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1AvatarOrgOrgNameApplicationsOauthAppIdWithResponse(context.TODO(), orgName, appId, &daliClient.DeleteV1AvatarOrgOrgNameApplicationsOauthAppIdParams{
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

var avatarApplicationsOauthCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new avatar for the OAuth application",
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

		appId, err := command.FlagOrPromptString(cmd, "appId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.PostV1AvatarOrgOrgNameApplicationsOauthAppIdWithResponse(context.TODO(), orgName, appId, &daliClient.PostV1AvatarOrgOrgNameApplicationsOauthAppIdParams{
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
