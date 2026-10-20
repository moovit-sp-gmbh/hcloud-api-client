package dali

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/daliClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	avatarSpacesAddressBookDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	avatarSpacesAddressBookDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	avatarSpacesAddressBookDeleteCommand.Flags().String("bookId", "", "")
	avatarSpacesAddressBookDeleteCommand.Flags().String("token", "", "authorization token")
	avatarSpacesAddressBookCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	avatarSpacesAddressBookCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	avatarSpacesAddressBookCreateCommand.Flags().String("bookId", "", "")
	avatarSpacesAddressBookCreateCommand.Flags().String("token", "", "authorization token")

	avatarSpacesAddressBookCommand.AddCommand(avatarSpacesAddressBookDeleteCommand)
	avatarSpacesAddressBookCommand.AddCommand(avatarSpacesAddressBookCreateCommand)
}

var avatarSpacesAddressBookCommand = &cobra.Command{
	Use:   "address-book",
	Short: "Commands for address-book",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var avatarSpacesAddressBookDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete avatar for the address book",
	Long:  "Removes the existing avatar of the specified address book",
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

		bookId, err := command.FlagOrPromptString(cmd, "bookId")
		if err != nil {
			return err
		}
		bookIdParsed, err := uuid.Parse(bookId)
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1AvatarOrgOrgNameSpacesSpaceNameAddressBookBookIdWithResponse(context.TODO(), orgName, spaceName, bookIdParsed, &daliClient.DeleteV1AvatarOrgOrgNameSpacesSpaceNameAddressBookBookIdParams{
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

var avatarSpacesAddressBookCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new avatar for the address book",
	Long:  "Generates and uploads a new default avatar for the specified address book and returns its public URL",
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

		bookId, err := command.FlagOrPromptString(cmd, "bookId")
		if err != nil {
			return err
		}
		bookIdParsed, err := uuid.Parse(bookId)
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.PostV1AvatarOrgOrgNameSpacesSpaceNameAddressBookBookIdWithResponse(context.TODO(), orgName, spaceName, bookIdParsed, &daliClient.PostV1AvatarOrgOrgNameSpacesSpaceNameAddressBookBookIdParams{
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
