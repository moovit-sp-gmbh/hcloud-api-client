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
	spacesDatabasesApikeyCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesApikeyCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesApikeyCreateCommand.Flags().String("dbName", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesDatabasesApikeyCreateCommand.Flags().String("token", "", "authorization token")
	spacesDatabasesApikeyCreateCommand.Flags().String("name", "", "Human-readable name for the API key to help identify its purpose (only alphanumeric characters, underscores and hyphens are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesApikeyCreateCommand.Flags().String("permission", "", "* read - Allows read-only access to resources without making modifications * write - Allows full access to create, update, and delete resources")
	spacesDatabasesApikeyDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesApikeyDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesApikeyDeleteCommand.Flags().String("dbName", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesDatabasesApikeyDeleteCommand.Flags().String("keyName", "", "Human-readable name for the API key to help identify its purpose (only alphanumeric characters, underscores and hyphens are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesApikeyDeleteCommand.Flags().String("token", "", "authorization token")

	spacesDatabasesApikeyCommand.AddCommand(spacesDatabasesApikeyCreateCommand)
	spacesDatabasesApikeyCommand.AddCommand(spacesDatabasesApikeyDeleteCommand)
}

var spacesDatabasesApikeyCommand = &cobra.Command{
	Use:   "apikey",
	Short: "Commands for apikey",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesDatabasesApikeyCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a database API key",
	Long:  "Generates a new API key for a specific database within the space",
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

		dbName, err := command.FlagOrPromptString(cmd, "dbName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		name, err := command.FlagOrPromptString(cmd, "name")
		if err != nil {
			return err
		}

		permission, err := command.FlagOrPromptString(cmd, "permission")
		if err != nil {
			return err
		}

		body := high5Client.PostV1OrgOrgNameSpacesSpaceNameDatabasesDbNameApikeyJSONRequestBody{
			Name:       name,
			Permission: high5Client.ApiKeyPermission(permission),
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameDatabasesDbNameApikeyWithResponse(context.TODO(), orgName, spaceName, dbName, &high5Client.PostV1OrgOrgNameSpacesSpaceNameDatabasesDbNameApikeyParams{
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

var spacesDatabasesApikeyDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete a database API key",
	Long:  "Removes an existing API key by its name, revoking access to the database",
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

		dbName, err := command.FlagOrPromptString(cmd, "dbName")
		if err != nil {
			return err
		}

		keyName, err := command.FlagOrPromptString(cmd, "keyName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1OrgOrgNameSpacesSpaceNameDatabasesDbNameApikeyKeyNameWithResponse(context.TODO(), orgName, spaceName, dbName, keyName, &high5Client.DeleteV1OrgOrgNameSpacesSpaceNameDatabasesDbNameApikeyKeyNameParams{
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
