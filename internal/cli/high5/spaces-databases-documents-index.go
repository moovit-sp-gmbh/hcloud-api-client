package high5

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	spacesDatabasesDocumentsIndexDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesDocumentsIndexDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesDocumentsIndexDeleteCommand.Flags().String("dbName", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesDatabasesDocumentsIndexDeleteCommand.Flags().String("key", "", "Key of the document (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 64 characters)")
	spacesDatabasesDocumentsIndexDeleteCommand.Flags().Float64("index", 0, "Numeric identifier representing the element’s position in the array")
	spacesDatabasesDocumentsIndexListCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesDocumentsIndexListCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesDocumentsIndexListCommand.Flags().String("dbName", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesDatabasesDocumentsIndexListCommand.Flags().String("key", "", "Key of the document (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 64 characters)")
	spacesDatabasesDocumentsIndexListCommand.Flags().Float64("index", 0, "Numeric identifier representing the element’s position in the array")

	spacesDatabasesDocumentsIndexCommand.AddCommand(spacesDatabasesDocumentsIndexDeleteCommand)
	spacesDatabasesDocumentsIndexCommand.AddCommand(spacesDatabasesDocumentsIndexListCommand)
}

var spacesDatabasesDocumentsIndexCommand = &cobra.Command{
	Use:   "index",
	Short: "Commands for index",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesDatabasesDocumentsIndexDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete element with specified index from the array",
	Long:  "Removes the element at the specified index from an array document",
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

		key, err := command.FlagOrPromptString(cmd, "key")
		if err != nil {
			return err
		}

		index, err := command.FlagOrPromptFloat(cmd, "index")
		if err != nil {
			return err
		}

		res, err := client.DeleteV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsKeyIndexIndexWithResponse(context.TODO(), orgName, spaceName, dbName, key, float32(index))
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

var spacesDatabasesDocumentsIndexListCommand = &cobra.Command{
	Use:   "list",
	Short: "Get element with specified index from the array",
	Long:  "Retrieves a single element from an array document by its zero-based index position",
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

		key, err := command.FlagOrPromptString(cmd, "key")
		if err != nil {
			return err
		}

		index, err := command.FlagOrPromptFloat(cmd, "index")
		if err != nil {
			return err
		}

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsKeyIndexIndexWithResponse(context.TODO(), orgName, spaceName, dbName, key, float32(index))
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
