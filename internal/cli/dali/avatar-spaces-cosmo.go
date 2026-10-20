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
	avatarSpacesCosmoDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	avatarSpacesCosmoDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	avatarSpacesCosmoDeleteCommand.Flags().String("token", "", "authorization token")
	avatarSpacesCosmoCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	avatarSpacesCosmoCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	avatarSpacesCosmoCreateCommand.Flags().String("token", "", "authorization token")

	avatarSpacesCosmoCommand.AddCommand(avatarSpacesCosmoDeleteCommand)
	avatarSpacesCosmoCommand.AddCommand(avatarSpacesCosmoCreateCommand)
	avatarSpacesCommand.AddCommand(avatarSpacesCosmoCommand)
	avatarSpacesCommand.AddCommand(avatarSpacesHigh5Command)
	avatarSpacesCommand.AddCommand(avatarSpacesAddressBookCommand)
}

var avatarSpacesCommand = &cobra.Command{
	Use:   "spaces",
	Short: "Commands for spaces",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var avatarSpacesCosmoCommand = &cobra.Command{
	Use:   "cosmo",
	Short: "Commands for cosmo",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var avatarSpacesCosmoDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete avatar for the Cosmo Space",
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

		spaceName, err := command.FlagOrPromptString(cmd, "spaceName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1AvatarOrgOrgNameSpacesCosmoSpaceNameWithResponse(context.TODO(), orgName, spaceName, &daliClient.DeleteV1AvatarOrgOrgNameSpacesCosmoSpaceNameParams{
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

var avatarSpacesCosmoCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new avatar for the Cosmo Space",
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

		spaceName, err := command.FlagOrPromptString(cmd, "spaceName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.PostV1AvatarOrgOrgNameSpacesCosmoSpaceNameWithResponse(context.TODO(), orgName, spaceName, &daliClient.PostV1AvatarOrgOrgNameSpacesCosmoSpaceNameParams{
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
