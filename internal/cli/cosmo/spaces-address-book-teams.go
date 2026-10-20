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
	spacesAddressBookTeamsDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesAddressBookTeamsDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesAddressBookTeamsDeleteCommand.Flags().String("bookId", "", "")
	spacesAddressBookTeamsDeleteCommand.Flags().String("token", "", "authorization token")
	spacesAddressBookTeamsDeleteCommand.Flags().StringSlice("teamIds", nil, "IDs of IDP internal teams referenced by the address book")
	spacesAddressBookTeamsCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesAddressBookTeamsCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesAddressBookTeamsCreateCommand.Flags().String("bookId", "", "")
	spacesAddressBookTeamsCreateCommand.Flags().String("token", "", "authorization token")
	spacesAddressBookTeamsCreateCommand.Flags().StringSlice("teamIds", nil, "IDs of IDP internal teams referenced by the address book")

	spacesAddressBookTeamsCommand.AddCommand(spacesAddressBookTeamsDeleteCommand)
	spacesAddressBookTeamsCommand.AddCommand(spacesAddressBookTeamsCreateCommand)
}

var spacesAddressBookTeamsCommand = &cobra.Command{
	Use:   "teams",
	Short: "Commands for teams",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesAddressBookTeamsDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Remove team references from an address book",
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

		teamIds, err := cmd.Flags().GetStringSlice("teamIds")
		if err != nil {
			return err
		}

		body := cosmoClient.DeleteV1OrgOrgNameSpacesSpaceNameAddressBookBookIdTeamsJSONRequestBody{
			TeamIds: teamIds,
		}

		res, err := client.DeleteV1OrgOrgNameSpacesSpaceNameAddressBookBookIdTeamsWithResponse(context.TODO(), orgName, spaceName, bookIdParsed, body, func(_ context.Context, req *http.Request) error {
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

var spacesAddressBookTeamsCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Add team references to an address book",
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

		teamIds, err := cmd.Flags().GetStringSlice("teamIds")
		if err != nil {
			return err
		}

		body := cosmoClient.PostV1OrgOrgNameSpacesSpaceNameAddressBookBookIdTeamsJSONRequestBody{
			TeamIds: teamIds,
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameAddressBookBookIdTeamsWithResponse(context.TODO(), orgName, spaceName, bookIdParsed, body, func(_ context.Context, req *http.Request) error {
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
