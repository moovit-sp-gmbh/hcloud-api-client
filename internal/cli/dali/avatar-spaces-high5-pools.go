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
	avatarSpacesHigh5PoolsDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	avatarSpacesHigh5PoolsDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	avatarSpacesHigh5PoolsDeleteCommand.Flags().String("poolName", "", "Name of the pool (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	avatarSpacesHigh5PoolsDeleteCommand.Flags().String("token", "", "authorization token")
	avatarSpacesHigh5PoolsCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	avatarSpacesHigh5PoolsCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	avatarSpacesHigh5PoolsCreateCommand.Flags().String("poolName", "", "Name of the pool (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	avatarSpacesHigh5PoolsCreateCommand.Flags().String("token", "", "authorization token")

	avatarSpacesHigh5PoolsCommand.AddCommand(avatarSpacesHigh5PoolsDeleteCommand)
	avatarSpacesHigh5PoolsCommand.AddCommand(avatarSpacesHigh5PoolsCreateCommand)
}

var avatarSpacesHigh5PoolsCommand = &cobra.Command{
	Use:   "pools",
	Short: "Commands for pools",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var avatarSpacesHigh5PoolsDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete avatar for the High5 Pool",
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

		poolName, err := command.FlagOrPromptString(cmd, "poolName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1AvatarOrgOrgNameSpacesHigh5SpaceNamePoolsPoolNameWithResponse(context.TODO(), orgName, spaceName, poolName, &daliClient.DeleteV1AvatarOrgOrgNameSpacesHigh5SpaceNamePoolsPoolNameParams{
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

var avatarSpacesHigh5PoolsCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new avatar for the High5 Pool",
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

		poolName, err := command.FlagOrPromptString(cmd, "poolName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.PostV1AvatarOrgOrgNameSpacesHigh5SpaceNamePoolsPoolNameWithResponse(context.TODO(), orgName, spaceName, poolName, &daliClient.PostV1AvatarOrgOrgNameSpacesHigh5SpaceNamePoolsPoolNameParams{
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
