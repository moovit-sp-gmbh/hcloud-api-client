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
	spacesDatabasesDocumentsCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesDocumentsCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesDocumentsCreateCommand.Flags().String("dbName", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesDatabasesDocumentsCreateCommand.Flags().String("key", "", "Key of the document (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 64 characters)")
	spacesDatabasesDocumentsCreateCommand.Flags().String("type", "", "Type of the document value")
	spacesDatabasesDocumentsCreateCommand.Flags().Int("ttl", 0, "Time-to-live in seconds (undefined if missing)")
	spacesDatabasesDocumentsDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesDocumentsDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesDocumentsDeleteCommand.Flags().String("dbName", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesDatabasesDocumentsDeleteCommand.Flags().String("key", "", "Key of the document (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 64 characters)")
	spacesDatabasesDocumentsGetCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesDocumentsGetCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesDocumentsGetCommand.Flags().String("dbName", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesDatabasesDocumentsGetCommand.Flags().String("key", "", "Key of the document (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 64 characters)")
	spacesDatabasesDocumentsPatchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesDocumentsPatchCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesDocumentsPatchCommand.Flags().String("dbName", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesDatabasesDocumentsPatchCommand.Flags().String("key", "", "Key of the document (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 64 characters)")
	spacesDatabasesDocumentsPatchCommand.Flags().Int("ttl", 0, "Time-to-live in seconds (undefined if missing)")
	spacesDatabasesDocumentsUpdateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesDocumentsUpdateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesDocumentsUpdateCommand.Flags().String("dbName", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesDatabasesDocumentsUpdateCommand.Flags().String("key", "", "Key of the document (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 64 characters)")
	spacesDatabasesDocumentsUpdateCommand.Flags().String("bodyKey", "", "Key of the document (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 64 characters)")
	spacesDatabasesDocumentsUpdateCommand.Flags().Int("ttl", 0, "Time-to-live in seconds (undefined if missing)")
	spacesDatabasesDocumentsUpdateCommand.Flags().String("type", "", "")
	spacesDatabasesDocumentsKeysCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesDocumentsKeysCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesDocumentsKeysCommand.Flags().String("dbName", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesDatabasesDocumentsKeysCommand.Flags().String("key", "", "Key of the document (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 64 characters)")
	spacesDatabasesDocumentsKeysCommand.Flags().StringSlice("keys", nil, "Keys to remove from the object document")
	spacesDatabasesDocumentsPartialCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesDocumentsPartialCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesDocumentsPartialCommand.Flags().String("dbName", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesDatabasesDocumentsPartialCommand.Flags().String("key", "", "Key of the document (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 64 characters)")
	spacesDatabasesDocumentsPartialCommand.Flags().String("body", "{}", "request body as JSON")
	spacesDatabasesDocumentsPushCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesDocumentsPushCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesDocumentsPushCommand.Flags().String("dbName", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesDatabasesDocumentsPushCommand.Flags().String("key", "", "Key of the document (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 64 characters)")
	spacesDatabasesDocumentsPushCommand.Flags().Int("ttl", 0, "Time-to-live in seconds (undefined if missing)")
	spacesDatabasesDocumentsUnshiftCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDatabasesDocumentsUnshiftCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDatabasesDocumentsUnshiftCommand.Flags().String("dbName", "", "Name of the database (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesDatabasesDocumentsUnshiftCommand.Flags().String("key", "", "Key of the document (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 64 characters)")
	spacesDatabasesDocumentsUnshiftCommand.Flags().Int("ttl", 0, "Time-to-live in seconds (undefined if missing)")

	spacesDatabasesDocumentsCommand.AddCommand(spacesDatabasesDocumentsCreateCommand)
	spacesDatabasesDocumentsCommand.AddCommand(spacesDatabasesDocumentsSearchCommand)
	spacesDatabasesDocumentsCommand.AddCommand(spacesDatabasesDocumentsDeleteCommand)
	spacesDatabasesDocumentsCommand.AddCommand(spacesDatabasesDocumentsGetCommand)
	spacesDatabasesDocumentsCommand.AddCommand(spacesDatabasesDocumentsPatchCommand)
	spacesDatabasesDocumentsCommand.AddCommand(spacesDatabasesDocumentsUpdateCommand)
	spacesDatabasesDocumentsCommand.AddCommand(spacesDatabasesDocumentsDecCommand)
	spacesDatabasesDocumentsCommand.AddCommand(spacesDatabasesDocumentsIncCommand)
	spacesDatabasesDocumentsCommand.AddCommand(spacesDatabasesDocumentsIndexCommand)
	spacesDatabasesDocumentsCommand.AddCommand(spacesDatabasesDocumentsKeysCommand)
	spacesDatabasesDocumentsCommand.AddCommand(spacesDatabasesDocumentsPartialCommand)
	spacesDatabasesDocumentsCommand.AddCommand(spacesDatabasesDocumentsPushCommand)
	spacesDatabasesDocumentsCommand.AddCommand(spacesDatabasesDocumentsUnshiftCommand)
}

var spacesDatabasesDocumentsCommand = &cobra.Command{
	Use:   "documents",
	Short: "Commands for documents",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesDatabasesDocumentsCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Add a new document to the database",
	Long:  "Creates a new document in the specified database. The `type` field determines the expected data type of `value`. Returns 409 if a document with the same key already exists in the space",
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

		typeValue, err := command.FlagOrPromptString(cmd, "type")
		if err != nil {
			return err
		}

		body := high5Client.PostV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsJSONRequestBody{
			Key:  key,
			Type: high5Client.DocumentType(typeValue),
		}

		if cmd.Flags().Changed("ttl") {
			ttl, err := cmd.Flags().GetInt("ttl")
			if err != nil {
				return err
			}
			body.Ttl = &ttl
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsWithResponse(context.TODO(), orgName, spaceName, dbName, body)
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

var spacesDatabasesDocumentsDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete a document from the database",
	Long:  "Permanently deletes the document identified by its key. For ARRAY-type documents, all associated array items stored in the sub-collection are also deleted",
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

		res, err := client.DeleteV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsKeyWithResponse(context.TODO(), orgName, spaceName, dbName, key)
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

var spacesDatabasesDocumentsGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Get document by its key",
	Long:  "Returns a single document identified by its key within the specified database. For ARRAY-type documents, the full array items are collected and included in the response value",
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

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsKeyWithResponse(context.TODO(), orgName, spaceName, dbName, key)
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

var spacesDatabasesDocumentsPatchCommand = &cobra.Command{
	Use:   "patch",
	Short: "Patch a document",
	Long:  "Updates the value stored in an existing document. The type of the new value must match the document's current type; a type mismatch returns an error. For ARRAY-type documents, existing items are replaced with the new ones",
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

		body := high5Client.PatchV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsKeyJSONRequestBody{}

		if cmd.Flags().Changed("ttl") {
			ttl, err := cmd.Flags().GetInt("ttl")
			if err != nil {
				return err
			}
			body.Ttl = &ttl
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsKeyWithResponse(context.TODO(), orgName, spaceName, dbName, key, body)
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

var spacesDatabasesDocumentsUpdateCommand = &cobra.Command{
	Use:   "update",
	Short: "Upsert a document",
	Long:  "Creates the document if it does not exist, or updates it if it does. The document is looked up by the `key` URL parameter. An optional `key` field in the body overrides the document key. If `type` is not provided, it is auto-detected from the value. When updating, the resolved type must match the existing document's type; a mismatch returns an error. An optional `ttl` sets the document's expiry time in seconds",
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

		body := high5Client.PutV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsKeyJSONRequestBody{}

		if cmd.Flags().Changed("bodyKey") {
			bodyKey, err := cmd.Flags().GetString("bodyKey")
			if err != nil {
				return err
			}
			body.Key = &bodyKey
		}

		if cmd.Flags().Changed("ttl") {
			ttl, err := cmd.Flags().GetInt("ttl")
			if err != nil {
				return err
			}
			body.Ttl = &ttl
		}

		if cmd.Flags().Changed("type") {
			typeValue, err := cmd.Flags().GetString("type")
			if err != nil {
				return err
			}
			typeValueValue := high5Client.DocumentUpsertDtoType(typeValue)
			body.Type = &typeValueValue
		}

		res, err := client.PutV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsKeyWithResponse(context.TODO(), orgName, spaceName, dbName, key, body)
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

var spacesDatabasesDocumentsKeysCommand = &cobra.Command{
	Use:   "keys",
	Short: "Remove one or more keys from an object document",
	Long:  "Removes one or more specified keys from an object document without affecting the remaining fields",
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

		keys, err := cmd.Flags().GetStringSlice("keys")
		if err != nil {
			return err
		}

		body := high5Client.PatchV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsKeyKeysJSONRequestBody{
			Keys: keys,
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsKeyKeysWithResponse(context.TODO(), orgName, spaceName, dbName, key, body)
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

var spacesDatabasesDocumentsPartialCommand = &cobra.Command{
	Use:   "partial",
	Short: "Add or update one or more key/value pairs inside an object document without replacing the entire value",
	Long:  "Creates or updates specific key/value pairs within an object document without overwriting other existing fields",
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

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body high5Client.PatchV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsKeyPartialJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsKeyPartialWithResponse(context.TODO(), orgName, spaceName, dbName, key, body)
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

var spacesDatabasesDocumentsPushCommand = &cobra.Command{
	Use:   "push",
	Short: "Add one or more elements to the end of an array",
	Long:  "Appends one or more elements to the end of an array document and returns the updated document",
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

		body := high5Client.PatchV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsKeyPushJSONRequestBody{}

		if cmd.Flags().Changed("ttl") {
			ttl, err := cmd.Flags().GetInt("ttl")
			if err != nil {
				return err
			}
			body.Ttl = &ttl
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsKeyPushWithResponse(context.TODO(), orgName, spaceName, dbName, key, body)
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

var spacesDatabasesDocumentsUnshiftCommand = &cobra.Command{
	Use:   "unshift",
	Short: "Add one or more elements to the beginning of an array",
	Long:  "Prepends one or more elements to the beginning of an array document and returns the updated document",
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

		body := high5Client.PatchV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsKeyUnshiftJSONRequestBody{}

		if cmd.Flags().Changed("ttl") {
			ttl, err := cmd.Flags().GetInt("ttl")
			if err != nil {
				return err
			}
			body.Ttl = &ttl
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameDatabasesDbNameDocumentsKeyUnshiftWithResponse(context.TODO(), orgName, spaceName, dbName, key, body)
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
