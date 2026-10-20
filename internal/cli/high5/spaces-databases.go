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
	spacesDatabasesCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesCreateCommand.Flags().String("token", "", "authorization token")
	spacesDatabasesCreateCommand.Flags().String("name", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesDatabasesSearchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesSearchCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesSearchCommand.Flags().String("token", "", "authorization token")
	spacesDatabasesSearchCommand.Flags().Float32("limit", 0, "Page size")
	spacesDatabasesSearchCommand.Flags().Float32("page", 0, "Page index")
	spacesDatabasesSearchCommand.Flags().String("body", "{}", "request body as JSON")
	spacesDatabasesDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesDeleteCommand.Flags().String("dbName", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesDatabasesDeleteCommand.Flags().String("token", "", "authorization token")
	spacesDatabasesGetCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesGetCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesGetCommand.Flags().String("dbName", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesDatabasesGetCommand.Flags().String("token", "", "authorization token")
	spacesDatabasesNameCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesNameCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesNameCommand.Flags().String("dbName", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesDatabasesNameCommand.Flags().String("token", "", "authorization token")
	spacesDatabasesNameCommand.Flags().String("name", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")

	spacesDatabasesCommand.AddCommand(spacesDatabasesCreateCommand)
	spacesDatabasesCommand.AddCommand(spacesDatabasesSearchCommand)
	spacesDatabasesCommand.AddCommand(spacesDatabasesDeleteCommand)
	spacesDatabasesCommand.AddCommand(spacesDatabasesGetCommand)
	spacesDatabasesCommand.AddCommand(spacesDatabasesApikeyCommand)
	spacesDatabasesCommand.AddCommand(spacesDatabasesDocumentsCommand)
	spacesDatabasesCommand.AddCommand(spacesDatabasesNameCommand)
}

var spacesDatabasesCommand = &cobra.Command{
	Use:   "databases",
	Short: "Commands for databases",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesDatabasesCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new database for a space",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		name, err := command.FlagOrPromptString(cmd, "name")
		if err != nil {
			return err
		}

		body := high5Client.PostV1OrgOrgNameSpacesSpaceNameDatabasesJSONRequestBody{
			Name: name,
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameDatabasesWithResponse(context.TODO(), orgName, spaceName, &high5Client.PostV1OrgOrgNameSpacesSpaceNameDatabasesParams{
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

var spacesDatabasesSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Search databases of a space",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &high5Client.PostV1OrgOrgNameSpacesSpaceNameDatabasesSearchParams{
			Authorization: jwt,
		}

		if cmd.Flags().Changed("limit") {
			limit, err := cmd.Flags().GetFloat32("limit")
			if err != nil {
				return err
			}
			params.Limit = &limit
		}

		if cmd.Flags().Changed("page") {
			page, err := cmd.Flags().GetFloat32("page")
			if err != nil {
				return err
			}
			params.Page = &page
		}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body high5Client.PostV1OrgOrgNameSpacesSpaceNameDatabasesSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameDatabasesSearchWithResponse(context.TODO(), orgName, spaceName, params, body)
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

var spacesDatabasesDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete  a database by it's name",
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

		res, err := client.DeleteV1OrgOrgNameSpacesSpaceNameDatabasesDbNameWithResponse(context.TODO(), orgName, spaceName, dbName, &high5Client.DeleteV1OrgOrgNameSpacesSpaceNameDatabasesDbNameParams{
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

var spacesDatabasesGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Retrieve the database by it's name",
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

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameDatabasesDbNameWithResponse(context.TODO(), orgName, spaceName, dbName, &high5Client.GetV1OrgOrgNameSpacesSpaceNameDatabasesDbNameParams{
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

var spacesDatabasesNameCommand = &cobra.Command{
	Use:   "name",
	Short: "Change name of the database",
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

		body := high5Client.PatchV1OrgOrgNameSpacesSpaceNameDatabasesDbNameNameJSONRequestBody{
			Name: name,
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameDatabasesDbNameNameWithResponse(context.TODO(), orgName, spaceName, dbName, &high5Client.PatchV1OrgOrgNameSpacesSpaceNameDatabasesDbNameNameParams{
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
