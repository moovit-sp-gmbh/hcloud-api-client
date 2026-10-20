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
	spacesDatabasesDocumentsIncPatchAllCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesDocumentsIncPatchAllCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesDocumentsIncPatchAllCommand.Flags().String("dbName", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesDatabasesDocumentsIncPatchAllCommand.Flags().String("key", "", "Key of the document (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 64 characters)")
	spacesDatabasesDocumentsIncPatchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesDocumentsIncPatchCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesDocumentsIncPatchCommand.Flags().String("dbName", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesDatabasesDocumentsIncPatchCommand.Flags().String("key", "", "Key of the document (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 64 characters)")
	spacesDatabasesDocumentsIncPatchCommand.Flags().Float64("value", 0, "Numeric value extracted from the path, stored in the document and used in subsequent processing")

	spacesDatabasesDocumentsIncCommand.AddCommand(spacesDatabasesDocumentsIncPatchAllCommand)
	spacesDatabasesDocumentsIncCommand.AddCommand(spacesDatabasesDocumentsIncPatchCommand)
}

var spacesDatabasesDocumentsIncCommand = &cobra.Command{
	Use:   "inc",
	Short: "Commands for inc",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesDatabasesDocumentsIncPatchAllCommand = &cobra.Command{
	Use:   "patch-all",
	Short: "Increment the numeric value stored in the document by one",
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

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsKeyIncWithResponse(context.TODO(), orgName, spaceName, dbName, key)
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

var spacesDatabasesDocumentsIncPatchCommand = &cobra.Command{
	Use:   "patch",
	Short: "Increase the numeric value stored in the document by custom value",
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

		value, err := command.FlagOrPromptFloat(cmd, "value")
		if err != nil {
			return err
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsKeyIncValueWithResponse(context.TODO(), orgName, spaceName, dbName, key, float32(value))
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
