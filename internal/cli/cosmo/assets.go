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
	assetsUploadCommand.Flags().String("token", "", "")
	assetsBrowseCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	assetsBrowseCommand.Flags().String("token", "", "authorization token")
	assetsBrowseCommand.Flags().String("context", "", "")
	assetsBrowseCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	assetsBrowseCommand.Flags().Float32("limit", 0, "Page size")
	assetsBrowseCommand.Flags().Float32("page", 0, "Page index")
	assetsBrowseCommand.Flags().String("parentId", "", "")
	assetsBrowseCommand.Flags().String("shareId", "", "")
	assetsSearchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	assetsSearchCommand.Flags().String("token", "", "authorization token")
	assetsSearchCommand.Flags().String("context", "", "")
	assetsSearchCommand.Flags().Float32("limit", 0, "Page size")
	assetsSearchCommand.Flags().Float32("page", 0, "Page index")
	assetsSearchCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	assetsSearchCommand.Flags().String("parentId", "", "")
	assetsSearchCommand.Flags().String("shareId", "", "")
	assetsSearchCommand.Flags().String("body", "{}", "request body as JSON")

	assetsCommand.AddCommand(assetsUploadCommand)
	assetsCommand.AddCommand(assetsBrowseCommand)
	assetsCommand.AddCommand(assetsSearchCommand)
}

var assetsCommand = &cobra.Command{
	Use:   "assets",
	Short: "Commands for assets",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var assetsUploadCommand = &cobra.Command{
	Use:   "upload",
	Short: "Upload an asset",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewCosmoClient(cmd)
		if err != nil {
			return err
		}

		token, err := command.FlagOrPromptString(cmd, "token")
		if err != nil {
			return err
		}

		res, err := client.PutV1AssetsUploadWithResponse(context.TODO(), &cosmoClient.PutV1AssetsUploadParams{
			Token: token,
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

		command.SetResult(cmd, text)

		return nil
	},
}

var assetsBrowseCommand = &cobra.Command{
	Use:   "browse",
	Short: "Browse the children of an asset",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &cosmoClient.PostV1OrgOrgNameAssetsBrowseParams{}

		if cmd.Flags().Changed("context") {
			contextValue, err := cmd.Flags().GetString("context")
			if err != nil {
				return err
			}
			contextValueValue := cosmoClient.PostV1OrgOrgNameAssetsBrowseParamsContext(contextValue)
			params.Context = &contextValueValue
		}

		if cmd.Flags().Changed("spaceName") {
			spaceName, err := cmd.Flags().GetString("spaceName")
			if err != nil {
				return err
			}
			params.SpaceName = &spaceName
		}

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

		if cmd.Flags().Changed("parentId") {
			parentId, err := cmd.Flags().GetString("parentId")
			if err != nil {
				return err
			}
			parentIdParsed, err := uuid.Parse(parentId)
			if err != nil {
				return err
			}
			params.ParentId = &parentIdParsed
		}

		if cmd.Flags().Changed("shareId") {
			shareId, err := cmd.Flags().GetString("shareId")
			if err != nil {
				return err
			}
			shareIdParsed, err := uuid.Parse(shareId)
			if err != nil {
				return err
			}
			params.ShareId = &shareIdParsed
		}

		body := cosmoClient.PostV1OrgOrgNameAssetsBrowseJSONRequestBody{}

		res, err := client.PostV1OrgOrgNameAssetsBrowseWithResponse(context.TODO(), orgName, params, body, func(_ context.Context, req *http.Request) error {
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

var assetsSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Search assets of an organization across all spaces",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &cosmoClient.PostV1OrgOrgNameAssetsSearchParams{}

		if cmd.Flags().Changed("context") {
			contextValue, err := cmd.Flags().GetString("context")
			if err != nil {
				return err
			}
			contextValueValue := cosmoClient.PostV1OrgOrgNameAssetsSearchParamsContext(contextValue)
			params.Context = &contextValueValue
		}

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

		if cmd.Flags().Changed("spaceName") {
			spaceName, err := cmd.Flags().GetString("spaceName")
			if err != nil {
				return err
			}
			params.SpaceName = &spaceName
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
			params.ParentId = &parentIdParsed
		}

		if cmd.Flags().Changed("shareId") {
			shareId, err := cmd.Flags().GetString("shareId")
			if err != nil {
				return err
			}
			shareIdParsed, err := uuid.Parse(shareId)
			if err != nil {
				return err
			}
			params.ShareId = &shareIdParsed
		}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body cosmoClient.PostV1OrgOrgNameAssetsSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameAssetsSearchWithResponse(context.TODO(), orgName, params, body, func(_ context.Context, req *http.Request) error {
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
