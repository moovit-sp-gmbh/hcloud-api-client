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
	spacesWatchfoldersScanListCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesWatchfoldersScanListCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesWatchfoldersScanListCommand.Flags().String("watchFolderId", "", "")
	spacesWatchfoldersScanListCommand.Flags().String("scanSecret", "", "")
	spacesWatchfoldersScanListCommand.Flags().String("token", "", "authorization token")
	spacesWatchfoldersScanCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesWatchfoldersScanCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesWatchfoldersScanCreateCommand.Flags().String("watchFolderId", "", "")
	spacesWatchfoldersScanCreateCommand.Flags().String("scanSecret", "", "")
	spacesWatchfoldersScanCreateCommand.Flags().String("token", "", "authorization token")
	spacesWatchfoldersScanCreateCommand.Flags().String("body", "{}", "request body as JSON")

	spacesWatchfoldersScanCommand.AddCommand(spacesWatchfoldersScanListCommand)
	spacesWatchfoldersScanCommand.AddCommand(spacesWatchfoldersScanCreateCommand)
}

var spacesWatchfoldersScanCommand = &cobra.Command{
	Use:   "scan",
	Short: "Commands for scan",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesWatchfoldersScanListCommand = &cobra.Command{
	Use:   "list",
	Short: "Get the scan configuration for a watch folder",
	Long:  "Returns the path and decrypted credentials to scan, for the agent that received this scan's NATS dispatch. Requires the one-time scan secret handed out with that dispatch, not general space read access.",
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

		scanSecret, err := command.FlagOrPromptString(cmd, "scanSecret")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameWatchfoldersWatchFolderIdScanScanSecretWithResponse(context.TODO(), orgName, spaceName, watchFolderId, scanSecret, &high5Client.GetV1OrgOrgNameSpacesSpaceNameWatchfoldersWatchFolderIdScanScanSecretParams{
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

var spacesWatchfoldersScanCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Submit a watch folder scan report",
	Long:  "Accepts a list of files discovered by the agent and updates file tracking state. Requires the one-time scan secret handed out with that scan's NATS dispatch, not general space write access.",
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

		scanSecret, err := command.FlagOrPromptString(cmd, "scanSecret")
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
		var body high5Client.PostV1OrgOrgNameSpacesSpaceNameWatchfoldersWatchFolderIdScanScanSecretJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameWatchfoldersWatchFolderIdScanScanSecretWithResponse(context.TODO(), orgName, spaceName, watchFolderId, scanSecret, &high5Client.PostV1OrgOrgNameSpacesSpaceNameWatchfoldersWatchFolderIdScanScanSecretParams{
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
