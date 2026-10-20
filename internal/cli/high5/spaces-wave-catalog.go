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
	spacesWaveCatalogListCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesWaveCatalogListCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesWaveCatalogListCommand.Flags().String("token", "", "authorization token")
	spacesWaveCatalogCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesWaveCatalogCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesWaveCatalogCreateCommand.Flags().String("token", "", "authorization token")
	spacesWaveCatalogCreateCommand.Flags().String("url", "", "Url to the catalog")
	spacesWaveCatalogCreateCommand.Flags().String("version", "", "Version that should be used to serve nodes from for a space")
	spacesWaveCatalogCreateCommand.Flags().String("minimumEngineVersion", "", "Minimum engine version the catalog can run against")
	spacesWaveCatalogDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesWaveCatalogDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesWaveCatalogDeleteCommand.Flags().String("catalogId", "", "")
	spacesWaveCatalogDeleteCommand.Flags().String("token", "", "authorization token")

	spacesWaveCatalogCommand.AddCommand(spacesWaveCatalogListCommand)
	spacesWaveCatalogCommand.AddCommand(spacesWaveCatalogCreateCommand)
	spacesWaveCatalogCommand.AddCommand(spacesWaveCatalogDeleteCommand)
	spacesWaveCommand.AddCommand(spacesWaveCatalogCommand)
	spacesWaveCommand.AddCommand(spacesWaveEngineCommand)
}

var spacesWaveCommand = &cobra.Command{
	Use:   "wave",
	Short: "Commands for wave",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesWaveCatalogCommand = &cobra.Command{
	Use:   "catalog",
	Short: "Commands for catalog",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesWaveCatalogListCommand = &cobra.Command{
	Use:   "list",
	Short: "Get the currently active wave catalogs for a space",
	Long:  "Returns the wave catalogs currently active for this space",
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

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameWaveCatalogWithResponse(context.TODO(), orgName, spaceName, &high5Client.GetV1OrgOrgNameSpacesSpaceNameWaveCatalogParams{
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

var spacesWaveCatalogCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Add or update a new wave catalog to a space",
	Long:  "Adds a new wave catalog to the space or updates an existing one, checking the license quota unless the catalog is already registered",
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

		body := high5Client.PostV1OrgOrgNameSpacesSpaceNameWaveCatalogJSONRequestBody{
			Url:     url,
			Version: version,
		}

		if cmd.Flags().Changed("minimumEngineVersion") {
			minimumEngineVersion, err := cmd.Flags().GetString("minimumEngineVersion")
			if err != nil {
				return err
			}
			body.MinimumEngineVersion = &minimumEngineVersion
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameWaveCatalogWithResponse(context.TODO(), orgName, spaceName, &high5Client.PostV1OrgOrgNameSpacesSpaceNameWaveCatalogParams{
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

var spacesWaveCatalogDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete a catalog from a space",
	Long:  "Removes the specified wave catalog from the space",
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

		catalogId, err := command.FlagOrPromptString(cmd, "catalogId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1OrgOrgNameSpacesSpaceNameWaveCatalogCatalogIdWithResponse(context.TODO(), orgName, spaceName, catalogId, &high5Client.DeleteV1OrgOrgNameSpacesSpaceNameWaveCatalogCatalogIdParams{
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
