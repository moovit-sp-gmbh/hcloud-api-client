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
	spacesDatabasesDocumentsSearchCreateAllByDbNameCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesDocumentsSearchCreateAllByDbNameCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesDocumentsSearchCreateAllByDbNameCommand.Flags().String("dbName", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesDatabasesDocumentsSearchCreateAllByDbNameCommand.Flags().Float32("limit", 0, "Page size")
	spacesDatabasesDocumentsSearchCreateAllByDbNameCommand.Flags().Float32("page", 0, "Page index")
	spacesDatabasesDocumentsSearchCreateAllByDbNameCommand.Flags().String("body", "{}", "request body as JSON")
	spacesDatabasesDocumentsSearchCreateAllByDbName2Command.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesDocumentsSearchCreateAllByDbName2Command.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesDocumentsSearchCreateAllByDbName2Command.Flags().String("dbName", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesDatabasesDocumentsSearchCreateAllByDbName2Command.Flags().Float32("limit", 0, "Page size")
	spacesDatabasesDocumentsSearchCreateAllByDbName2Command.Flags().Float32("page", 0, "Page index")
	spacesDatabasesDocumentsSearchCreateAllByDbName2Command.Flags().String("body", "{}", "request body as JSON")

	spacesDatabasesDocumentsSearchCreateAllCommand.AddCommand(spacesDatabasesDocumentsSearchCreateAllByDbNameCommand)
	spacesDatabasesDocumentsSearchCreateAllCommand.AddCommand(spacesDatabasesDocumentsSearchCreateAllByDbName2Command)
	spacesDatabasesDocumentsSearchCommand.AddCommand(spacesDatabasesDocumentsSearchCreateAllCommand)
}

var spacesDatabasesDocumentsSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Commands for search",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesDatabasesDocumentsSearchCreateAllCommand = &cobra.Command{
	Use:   "create-all",
	Short: "Commands for create-all",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesDatabasesDocumentsSearchCreateAllByDbNameCommand = &cobra.Command{
	Use:   "by-db-name",
	Short: "Search for items within a specified array",
	Long:  "Searches for items within an array document using the specified filter criteria and returns a paginated list of matching results",
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

		params := &high5Client.PostV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsKeySearchParams{}

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
		var body high5Client.PostV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsKeySearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsKeySearchWithResponse(context.TODO(), orgName, spaceName, dbName, params, body)
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

var spacesDatabasesDocumentsSearchCreateAllByDbName2Command = &cobra.Command{
	Use:   "by-db-name-2",
	Short: "Search for documents by database name and search filter",
	Long:  "Returns a paginated list of documents in the specified database that match the given search filter. Use `limit` and `page` query parameters for pagination. The `total` response header contains the total number of matching documents regardless of pagination",
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

		params := &high5Client.PostV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsSearchParams{}

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
		var body high5Client.PostV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsSearchWithResponse(context.TODO(), orgName, spaceName, dbName, params, body)
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
