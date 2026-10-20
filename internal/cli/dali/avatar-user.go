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
	avatarUserDeleteCommand.Flags().String("token", "", "authorization token")
	avatarUserCreateCommand.Flags().String("token", "", "authorization token")

	avatarUserCommand.AddCommand(avatarUserDeleteCommand)
	avatarUserCommand.AddCommand(avatarUserCreateCommand)
}

var avatarUserCommand = &cobra.Command{
	Use:   "user",
	Short: "Commands for user",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var avatarUserDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete an Avatar for user: recursively deletes all images associated with the current user and the relevant user organization",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1AvatarUserWithResponse(context.TODO(), &daliClient.DeleteV1AvatarUserParams{
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

var avatarUserCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new Avatar for user",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.PostV1AvatarUserWithResponse(context.TODO(), &daliClient.PostV1AvatarUserParams{
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
