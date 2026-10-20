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
	spacesWaveEngineListCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesWaveEngineListCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesWaveEngineListCommand.Flags().String("token", "", "authorization token")
	spacesWaveEnginePatchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesWaveEnginePatchCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesWaveEnginePatchCommand.Flags().String("token", "", "authorization token")
	spacesWaveEnginePatchCommand.Flags().String("url", "", "Url to the engine")
	spacesWaveEnginePatchCommand.Flags().String("version", "", "Version of the release (it should contain between 1 and 128 characters)")
	spacesWaveEnginePatchCommand.Flags().Bool("latest", false, "Flag to use latest engine. Will be auto updated")

	spacesWaveEngineCommand.AddCommand(spacesWaveEngineListCommand)
	spacesWaveEngineCommand.AddCommand(spacesWaveEnginePatchCommand)
}

var spacesWaveEngineCommand = &cobra.Command{
	Use:   "engine",
	Short: "Commands for engine",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesWaveEngineListCommand = &cobra.Command{
	Use:   "list",
	Short: "Get the currently active wave engine for a space",
	Long:  "Returns the Wave Engine currently active for this space",
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

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameWaveEngineWithResponse(context.TODO(), orgName, spaceName, &high5Client.GetV1OrgOrgNameSpacesSpaceNameWaveEngineParams{
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

var spacesWaveEnginePatchCommand = &cobra.Command{
	Use:   "patch",
	Short: "Update the default wave engine that should be used on stream execution for that space",
	Long:  "Sets the default Wave Engine version used when executing streams in this space",
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

		url, err := command.FlagOrPromptString(cmd, "url")
		if err != nil {
			return err
		}

		version, err := command.FlagOrPromptString(cmd, "version")
		if err != nil {
			return err
		}

		body := high5Client.PatchV1OrgOrgNameSpacesSpaceNameWaveEngineJSONRequestBody{
			Url:     url,
			Version: version,
		}

		if cmd.Flags().Changed("latest") {
			latestValue, err := cmd.Flags().GetBool("latest")
			if err != nil {
				return err
			}
			body.Latest = &latestValue
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameWaveEngineWithResponse(context.TODO(), orgName, spaceName, &high5Client.PatchV1OrgOrgNameSpacesSpaceNameWaveEngineParams{
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
