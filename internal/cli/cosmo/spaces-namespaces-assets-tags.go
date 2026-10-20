package cosmo

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	spacesNamespacesAssetsTagsDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesAssetsTagsDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNamespacesAssetsTagsDeleteCommand.Flags().String("namespaceName", "", "Name of the namespace (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesAssetsTagsDeleteCommand.Flags().String("assetId", "", "Asset ID")
	spacesNamespacesAssetsTagsDeleteCommand.Flags().String("token", "", "authorization token")
	spacesNamespacesAssetsTagsPatchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesAssetsTagsPatchCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNamespacesAssetsTagsPatchCommand.Flags().String("namespaceName", "", "Name of the namespace (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesAssetsTagsPatchCommand.Flags().String("assetId", "", "Asset ID")
	spacesNamespacesAssetsTagsPatchCommand.Flags().String("tagName", "", "Name of the tag (it should contain between 1 and 20 characters)")
	spacesNamespacesAssetsTagsPatchCommand.Flags().String("token", "", "authorization token")

	spacesNamespacesAssetsTagsCommand.AddCommand(spacesNamespacesAssetsTagsDeleteCommand)
	spacesNamespacesAssetsTagsCommand.AddCommand(spacesNamespacesAssetsTagsPatchCommand)
	spacesNamespacesAssetsCommand.AddCommand(spacesNamespacesAssetsTagsCommand)
}

var spacesNamespacesAssetsCommand = &cobra.Command{
	Use:   "assets",
	Short: "Commands for assets",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesNamespacesAssetsTagsCommand = &cobra.Command{
	Use:   "tags",
	Short: "Commands for tags",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesNamespacesAssetsTagsDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Remove a tag from an asset",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewCosmoClient(cmd)
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

		namespaceName, err := command.FlagOrPromptString(cmd, "namespaceName")
		if err != nil {
			return err
		}

		assetId, err := command.FlagOrPromptString(cmd, "assetId")
		if err != nil {
			return err
		}
		assetIdParsed, err := uuid.Parse(assetId)
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameAssetsAssetIdTagsWithResponse(context.TODO(), orgName, spaceName, namespaceName, assetIdParsed, func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", jwt)
			return nil
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

var spacesNamespacesAssetsTagsPatchCommand = &cobra.Command{
	Use:   "patch",
	Short: "Add a tag to an asset",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewCosmoClient(cmd)
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

		namespaceName, err := command.FlagOrPromptString(cmd, "namespaceName")
		if err != nil {
			return err
		}

		assetId, err := command.FlagOrPromptString(cmd, "assetId")
		if err != nil {
			return err
		}
		assetIdParsed, err := uuid.Parse(assetId)
		if err != nil {
			return err
		}

		tagName, err := command.FlagOrPromptString(cmd, "tagName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameAssetsAssetIdTagsTagNameWithResponse(context.TODO(), orgName, spaceName, namespaceName, assetIdParsed, tagName, func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", jwt)
			return nil
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
