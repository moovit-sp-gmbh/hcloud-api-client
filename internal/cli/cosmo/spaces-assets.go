package cosmo

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/cosmoClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	spacesAssetsDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesAssetsDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesAssetsDeleteCommand.Flags().String("token", "", "authorization token")
	spacesAssetsDeleteCommand.Flags().String("body", "{}", "request body as JSON")
	spacesAssetsCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesAssetsCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesAssetsCreateCommand.Flags().String("token", "", "authorization token")
	spacesAssetsCreateCommand.Flags().String("name", "", "Media_Asset name. It must include a file extension")
	spacesAssetsCreateCommand.Flags().String("parentId", "", "Parent ID. If omitted, the asset will be placed at the root level.")
	spacesAssetsBulkDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesAssetsBulkDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesAssetsBulkDeleteCommand.Flags().String("token", "", "authorization token")
	spacesAssetsBulkDeleteCommand.Flags().Int("limit", 0, "Maximum number of assets to delete. -1 deletes all matches.")
	spacesAssetsBulkDeleteCommand.Flags().Float32("page", 0, "Page index")
	spacesAssetsBulkDeleteCommand.Flags().String("body", "{}", "request body as JSON")
	spacesAssetsBulkMoveCopyCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesAssetsBulkMoveCopyCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesAssetsBulkMoveCopyCommand.Flags().String("token", "", "authorization token")
	spacesAssetsBulkMoveCopyCommand.Flags().String("body", "{}", "request body as JSON")
	spacesAssetsCloneCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesAssetsCloneCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesAssetsCloneCommand.Flags().String("token", "", "authorization token")
	spacesAssetsCloneCommand.Flags().String("body", "{}", "request body as JSON")
	spacesAssetsDownloadUrlsCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesAssetsDownloadUrlsCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesAssetsDownloadUrlsCommand.Flags().String("token", "", "authorization token")
	spacesAssetsDownloadUrlsCommand.Flags().String("body", "{}", "request body as JSON")
	spacesAssetsPermissionsCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesAssetsPermissionsCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesAssetsPermissionsCommand.Flags().String("token", "", "authorization token")
	spacesAssetsPermissionsCommand.Flags().String("body", "{}", "request body as JSON")
	spacesAssetsGetCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesAssetsGetCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesAssetsGetCommand.Flags().String("assetId", "", "")
	spacesAssetsGetCommand.Flags().String("token", "", "authorization token")
	spacesAssetsActivityCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesAssetsActivityCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesAssetsActivityCommand.Flags().String("assetId", "", "")
	spacesAssetsActivityCommand.Flags().String("token", "", "authorization token")
	spacesAssetsActivityCommand.Flags().Float32("limit", 0, "Page size")
	spacesAssetsActivityCommand.Flags().Float32("page", 0, "Page index")
	spacesAssetsLocationCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesAssetsLocationCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesAssetsLocationCommand.Flags().String("assetId", "", "")
	spacesAssetsLocationCommand.Flags().String("token", "", "authorization token")
	spacesAssetsLocationCommand.Flags().String("locationId", "", "Location ID")
	spacesAssetsNameCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesAssetsNameCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesAssetsNameCommand.Flags().String("assetId", "", "Asset ID")
	spacesAssetsNameCommand.Flags().String("token", "", "authorization token")
	spacesAssetsNameCommand.Flags().String("name", "", "Media_Asset name")
	spacesAssetsThumbnailCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesAssetsThumbnailCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesAssetsThumbnailCommand.Flags().String("assetId", "", "")
	spacesAssetsThumbnailCommand.Flags().String("token", "", "authorization token")
	spacesAssetsThumbnailCommand.Flags().Float64("thumbnailTime", 0, "Time in milliseconds to grab the thumbnail from")
	spacesAssetsUsersCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesAssetsUsersCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesAssetsUsersCommand.Flags().String("assetId", "", "")
	spacesAssetsUsersCommand.Flags().String("token", "", "authorization token")
	spacesAssetsUsersCommand.Flags().Float32("limit", 0, "Page size")
	spacesAssetsUsersCommand.Flags().Float32("page", 0, "Page index")
	spacesAssetsUsersCommand.Flags().String("search", "", "")
	spacesAssetsVttCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesAssetsVttCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesAssetsVttCommand.Flags().String("assetId", "", "")
	spacesAssetsVttCommand.Flags().String("token", "", "authorization token")
	spacesAssetsWaveformCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesAssetsWaveformCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesAssetsWaveformCommand.Flags().String("assetId", "", "")
	spacesAssetsWaveformCommand.Flags().String("token", "", "authorization token")

	spacesAssetsCommand.AddCommand(spacesAssetsDeleteCommand)
	spacesAssetsCommand.AddCommand(spacesAssetsCreateCommand)
	spacesAssetsCommand.AddCommand(spacesAssetsBulkDeleteCommand)
	spacesAssetsCommand.AddCommand(spacesAssetsBulkMoveCopyCommand)
	spacesAssetsCommand.AddCommand(spacesAssetsCloneCommand)
	spacesAssetsCommand.AddCommand(spacesAssetsDownloadUrlsCommand)
	spacesAssetsCommand.AddCommand(spacesAssetsPermissionsCommand)
	spacesAssetsCommand.AddCommand(spacesAssetsSubscribeCommand)
	spacesAssetsCommand.AddCommand(spacesAssetsGetCommand)
	spacesAssetsCommand.AddCommand(spacesAssetsActivityCommand)
	spacesAssetsCommand.AddCommand(spacesAssetsLocationCommand)
	spacesAssetsCommand.AddCommand(spacesAssetsNameCommand)
	spacesAssetsCommand.AddCommand(spacesAssetsThumbnailCommand)
	spacesAssetsCommand.AddCommand(spacesAssetsUsersCommand)
	spacesAssetsCommand.AddCommand(spacesAssetsVttCommand)
	spacesAssetsCommand.AddCommand(spacesAssetsWaveformCommand)
}

var spacesAssetsCommand = &cobra.Command{
	Use:   "assets",
	Short: "Commands for assets",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesAssetsDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete assets",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &cosmoClient.DeleteV1OrgOrgNameSpacesSpaceNameAssetsParams{}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body cosmoClient.DeleteV1OrgOrgNameSpacesSpaceNameAssetsJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.DeleteV1OrgOrgNameSpacesSpaceNameAssetsWithResponse(context.TODO(), orgName, spaceName, params, body, func(_ context.Context, req *http.Request) error {
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

var spacesAssetsCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create an asset",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		name, err := command.FlagOrPromptString(cmd, "name")
		if err != nil {
			return err
		}

		body := cosmoClient.PostV1OrgOrgNameSpacesSpaceNameAssetsJSONRequestBody{
			Name: name,
		}

		if cmd.Flags().Changed("parentId") {
			parentId, err := cmd.Flags().GetString("parentId")
			if err != nil {
				return err
			}
			parentIdParsed, err := uuid.Parse(parentId)
			if err != nil {
				return err
			}
			body.ParentId = &parentIdParsed
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameAssetsWithResponse(context.TODO(), orgName, spaceName, body, func(_ context.Context, req *http.Request) error {
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

var spacesAssetsBulkDeleteCommand = &cobra.Command{
	Use:   "bulk-delete",
	Short: "Bulk delete assets by filter",
	Long:  "Deletes the assets of the space that match the given filters. All filters must match for an asset to be deleted (AND combination).\nTo match one of multiple values for the same field, provide an array as the filter value.\nDepending on the mode, the matching assets are either moved to the trash bin (TRASHBIN, default) or permanently deleted (PERMANENT).\nFolders are never deleted by this endpoint and assets that are already in the trash bin are ignored. Stacks are deleted as a whole, including all of their versions.\nAt most `limit` assets are deleted per call. Matching more than that is not an error - `total` reports how many assets match in total, so the caller can repeat the call until `total` is 0.\nWhen deleting in batches, keep requesting page 0: deleting a page shifts the remaining matches, so incrementing the page would skip assets.\nIf the user is not allowed to delete one of the matching assets, the call fails. Moving to the trash bin is atomic, a permanent deletion commits in batches, so assets removed by an earlier batch stay deleted.",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &cosmoClient.DeleteV1OrgOrgNameSpacesSpaceNameAssetsBulkDeleteParams{}

		if cmd.Flags().Changed("limit") {
			limit, err := cmd.Flags().GetInt("limit")
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
		var body cosmoClient.DeleteV1OrgOrgNameSpacesSpaceNameAssetsBulkDeleteJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.DeleteV1OrgOrgNameSpacesSpaceNameAssetsBulkDeleteWithResponse(context.TODO(), orgName, spaceName, params, body, func(_ context.Context, req *http.Request) error {
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

var spacesAssetsBulkMoveCopyCommand = &cobra.Command{
	Use:   "bulk-move-copy",
	Short: "Bulk move or copy multiple assets",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body cosmoClient.PutV1OrgOrgNameSpacesSpaceNameAssetsBulkMoveCopyJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PutV1OrgOrgNameSpacesSpaceNameAssetsBulkMoveCopyWithResponse(context.TODO(), orgName, spaceName, body, func(_ context.Context, req *http.Request) error {
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

var spacesAssetsCloneCommand = &cobra.Command{
	Use:   "clone",
	Short: "Clone assets",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body cosmoClient.PutV1OrgOrgNameSpacesSpaceNameAssetsCloneJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PutV1OrgOrgNameSpacesSpaceNameAssetsCloneWithResponse(context.TODO(), orgName, spaceName, body, func(_ context.Context, req *http.Request) error {
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

var spacesAssetsDownloadUrlsCommand = &cobra.Command{
	Use:   "download-urls",
	Short: "Get download URLs for multiple assets",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body cosmoClient.PostV1OrgOrgNameSpacesSpaceNameAssetsDownloadUrlsJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameAssetsDownloadUrlsWithResponse(context.TODO(), orgName, spaceName, body, func(_ context.Context, req *http.Request) error {
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

var spacesAssetsPermissionsCommand = &cobra.Command{
	Use:   "permissions",
	Short: "Fetch permissions for a group of assets.",
	Long:  "The permissions are related to the user making the request.",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body cosmoClient.PostV1OrgOrgNameSpacesSpaceNameAssetsPermissionsJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameAssetsPermissionsWithResponse(context.TODO(), orgName, spaceName, body, func(_ context.Context, req *http.Request) error {
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

var spacesAssetsGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Fetch an asset",
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

		params := &cosmoClient.GetV1OrgOrgNameSpacesSpaceNameAssetsAssetIdParams{}

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameAssetsAssetIdWithResponse(context.TODO(), orgName, spaceName, assetIdParsed, params, func(_ context.Context, req *http.Request) error {
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

var spacesAssetsActivityCommand = &cobra.Command{
	Use:   "activity",
	Short: "Fetch activity relative to an asset",
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

		params := &cosmoClient.GetV1OrgOrgNameSpacesSpaceNameAssetsAssetIdActivityParams{}

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

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameAssetsAssetIdActivityWithResponse(context.TODO(), orgName, spaceName, assetIdParsed, params, func(_ context.Context, req *http.Request) error {
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

var spacesAssetsLocationCommand = &cobra.Command{
	Use:   "location",
	Short: "Attach an asset to a location",
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

		locationId, err := command.FlagOrPromptString(cmd, "locationId")
		if err != nil {
			return err
		}
		locationIdParsed, err := uuid.Parse(locationId)
		if err != nil {
			return err
		}

		body := cosmoClient.PatchV1OrgOrgNameSpacesSpaceNameAssetsAssetIdLocationJSONRequestBody{
			LocationId: locationIdParsed,
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameAssetsAssetIdLocationWithResponse(context.TODO(), orgName, spaceName, assetIdParsed, body, func(_ context.Context, req *http.Request) error {
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

var spacesAssetsNameCommand = &cobra.Command{
	Use:   "name",
	Short: "Rename an asset",
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

		name, err := command.FlagOrPromptString(cmd, "name")
		if err != nil {
			return err
		}

		body := cosmoClient.PatchV1OrgOrgNameSpacesSpaceNameAssetsAssetIdNameJSONRequestBody{
			Name: name,
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameAssetsAssetIdNameWithResponse(context.TODO(), orgName, spaceName, assetIdParsed, body, func(_ context.Context, req *http.Request) error {
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

var spacesAssetsThumbnailCommand = &cobra.Command{
	Use:   "thumbnail",
	Short: "Get a thumbnail for a video asset at a specific timecode",
	Long:  "Returns a thumbnail image (binary) generated from the video asset at the given timecode.",
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

		thumbnailTime, err := command.FlagOrPromptFloat(cmd, "thumbnailTime")
		if err != nil {
			return err
		}

		body := cosmoClient.PostV1OrgOrgNameSpacesSpaceNameAssetsAssetIdThumbnailJSONRequestBody{
			ThumbnailTime: float32(thumbnailTime),
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameAssetsAssetIdThumbnailWithResponse(context.TODO(), orgName, spaceName, assetIdParsed, body, func(_ context.Context, req *http.Request) error {
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

var spacesAssetsUsersCommand = &cobra.Command{
	Use:   "users",
	Short: "List users and teams that are allowed to be mentioned for the asset",
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

		params := &cosmoClient.GetV1OrgOrgNameSpacesSpaceNameAssetsAssetIdUsersParams{}

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

		if cmd.Flags().Changed("search") {
			search, err := cmd.Flags().GetString("search")
			if err != nil {
				return err
			}
			params.Search = &search
		}

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameAssetsAssetIdUsersWithResponse(context.TODO(), orgName, spaceName, assetIdParsed, params, func(_ context.Context, req *http.Request) error {
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

var spacesAssetsVttCommand = &cobra.Command{
	Use:   "vtt",
	Short: "Get VTT for a specific asset",
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

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameAssetsAssetIdVttWithResponse(context.TODO(), orgName, spaceName, assetIdParsed, func(_ context.Context, req *http.Request) error {
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

var spacesAssetsWaveformCommand = &cobra.Command{
	Use:   "waveform",
	Short: "Get waveform data for one video asset",
	Long:  "Returns multi-channel audio waveform data for one asset, sourced from HCU.",
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

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameAssetsAssetIdWaveformWithResponse(context.TODO(), orgName, spaceName, assetIdParsed, func(_ context.Context, req *http.Request) error {
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
