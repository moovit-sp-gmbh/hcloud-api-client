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
	spacesWatchfoldersCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesWatchfoldersCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesWatchfoldersCreateCommand.Flags().String("token", "", "authorization token")
	spacesWatchfoldersCreateCommand.Flags().String("body", "{}", "request body as JSON")
	spacesWatchfoldersDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesWatchfoldersDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesWatchfoldersDeleteCommand.Flags().String("watchFolderId", "", "")
	spacesWatchfoldersDeleteCommand.Flags().String("token", "", "authorization token")
	spacesWatchfoldersGetCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesWatchfoldersGetCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesWatchfoldersGetCommand.Flags().String("watchFolderId", "", "")
	spacesWatchfoldersGetCommand.Flags().String("token", "", "authorization token")
	spacesWatchfoldersPatchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesWatchfoldersPatchCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesWatchfoldersPatchCommand.Flags().String("watchFolderId", "", "")
	spacesWatchfoldersPatchCommand.Flags().String("token", "", "authorization token")
	spacesWatchfoldersPatchCommand.Flags().String("body", "{}", "request body as JSON")
	spacesWatchfoldersFilesCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesWatchfoldersFilesCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesWatchfoldersFilesCommand.Flags().String("watchFolderId", "", "")
	spacesWatchfoldersFilesCommand.Flags().String("token", "", "authorization token")

	spacesWatchfoldersCommand.AddCommand(spacesWatchfoldersCreateCommand)
	spacesWatchfoldersCommand.AddCommand(spacesWatchfoldersDeleteCommand)
	spacesWatchfoldersCommand.AddCommand(spacesWatchfoldersGetCommand)
	spacesWatchfoldersCommand.AddCommand(spacesWatchfoldersPatchCommand)
	spacesWatchfoldersCommand.AddCommand(spacesWatchfoldersFilesCommand)
	spacesWatchfoldersCommand.AddCommand(spacesWatchfoldersScanCommand)
}

var spacesWatchfoldersCommand = &cobra.Command{
	Use:   "watchfolders",
	Short: "Commands for watchfolders",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesWatchfoldersCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new watch folder",
	Long:  "Creates a new watch folder in the specified space and organization",
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

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body high5Client.PostV1OrgOrgNameSpacesSpaceNameWatchfoldersJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameWatchfoldersWithResponse(context.TODO(), orgName, spaceName, &high5Client.PostV1OrgOrgNameSpacesSpaceNameWatchfoldersParams{
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

var spacesWatchfoldersDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete a watch folder by it's ID",
	Long:  "Permanently deletes a watch folder from the system",
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

		watchFolderId, err := command.FlagOrPromptString(cmd, "watchFolderId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1OrgOrgNameSpacesSpaceNameWatchfoldersWatchFolderIdWithResponse(context.TODO(), orgName, spaceName, watchFolderId, &high5Client.DeleteV1OrgOrgNameSpacesSpaceNameWatchfoldersWatchFolderIdParams{
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

var spacesWatchfoldersGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Get a watch folder by it's ID",
	Long:  "Retrieves a specific watch folder by its unique database identifier",
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

		watchFolderId, err := command.FlagOrPromptString(cmd, "watchFolderId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameWatchfoldersWatchFolderIdWithResponse(context.TODO(), orgName, spaceName, watchFolderId, &high5Client.GetV1OrgOrgNameSpacesSpaceNameWatchfoldersWatchFolderIdParams{
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

var spacesWatchfoldersPatchCommand = &cobra.Command{
	Use:   "patch",
	Short: "Update existing watch folder",
	Long:  "Updates an existing watch folder's properties",
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

		watchFolderId, err := command.FlagOrPromptString(cmd, "watchFolderId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body high5Client.PatchV1OrgOrgNameSpacesSpaceNameWatchfoldersWatchFolderIdJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameWatchfoldersWatchFolderIdWithResponse(context.TODO(), orgName, spaceName, watchFolderId, &high5Client.PatchV1OrgOrgNameSpacesSpaceNameWatchfoldersWatchFolderIdParams{
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

var spacesWatchfoldersFilesCommand = &cobra.Command{
	Use:   "files",
	Short: "Get scan results for a watch folder",
	Long:  "Returns the list of files metadata currently tracked by the watch folder",
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

		watchFolderId, err := command.FlagOrPromptString(cmd, "watchFolderId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameWatchfoldersWatchFolderIdFilesWithResponse(context.TODO(), orgName, spaceName, watchFolderId, &high5Client.GetV1OrgOrgNameSpacesSpaceNameWatchfoldersWatchFolderIdFilesParams{
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
