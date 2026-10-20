package idp

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/idpClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	storageCommand.AddCommand(storageCreateCommand)
	storageCommand.AddCommand(storageSearchCommand)
	storageCommand.AddCommand(storageDeleteCommand)
	storageCommand.AddCommand(storagePatchCommand)
	organizationStorageTestCommand.Flags().String("storageId", "", "")
	organizationStorageTestCommand.Flags().String("token", "", "authorization token")
	storageCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	storageCreateCommand.Flags().String("token", "", "authorization token")
	storageCreateCommand.Flags().Bool("noAvatar", false, "")
	storageCreateCommand.Flags().String("accessKeyId", "", "")
	storageCreateCommand.Flags().String("bucket", "", "")
	storageCreateCommand.Flags().String("endpoint", "", "")
	storageCreateCommand.Flags().String("name", "", "Name of the storage (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	storageCreateCommand.Flags().String("region", "", "")
	storageCreateCommand.Flags().String("secretAccessKey", "", "")
	storageSearchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	storageSearchCommand.Flags().String("token", "", "authorization token")
	storageSearchCommand.Flags().Float32("limit", 0, "Page size")
	storageSearchCommand.Flags().Float32("page", 0, "Page index")
	storageSearchCommand.Flags().String("body", "{}", "request body as JSON")
	storageDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	storageDeleteCommand.Flags().String("storageId", "", "")
	storageDeleteCommand.Flags().String("token", "", "authorization token")
	storagePatchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	storagePatchCommand.Flags().String("storageId", "", "")
	storagePatchCommand.Flags().String("token", "", "authorization token")
	storagePatchCommand.Flags().String("accessKeyId", "", "")
	storagePatchCommand.Flags().String("bucket", "", "")
	storagePatchCommand.Flags().String("endpoint", "", "")
	storagePatchCommand.Flags().String("name", "", "Name of the storage (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	storagePatchCommand.Flags().String("region", "", "")
	storagePatchCommand.Flags().String("secretAccessKey", "", "")
	organizationStorageCommand.AddCommand(organizationStorageTestCommand)
}

var organizationStorageCommand = &cobra.Command{
	Use:   "storage",
	Short: "Commands for storage",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var organizationStorageTestCommand = &cobra.Command{
	Use:   "test",
	Short: "Test storage configuration",
	Long:  "\nRequired scopes: IDP_STORAGE_READ",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		storageId, err := command.GetArg(cmd, command.StringFlagIfSet("storageId"), command.ContextIfValidString(latest.StorageID), command.PromptString("storage ID"))
		if err != nil {
			return err
		}
		latest.StorageID = sql.NullString{Valid: true, String: storageId}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.GetV1OrgOrgNameStorageStorageIdTestWithResponse(context.TODO(), storageId, func(_ context.Context, req *http.Request) error {
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

var storageCommand = &cobra.Command{
	Use:   "storage",
	Short: "Commands for storage",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var storageCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new storage",
	Long:  "Required scopes: IDP_ORGANIZATION_WRITE",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
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

		params := &idpClient.PostV1OrgOrgNameStorageParams{}

		if cmd.Flags().Changed("noAvatar") {
			noAvatar, err := cmd.Flags().GetBool("noAvatar")
			if err != nil {
				return err
			}
			params.NoAvatar = &noAvatar
		}

		accessKeyId, err := command.FlagOrPromptString(cmd, "accessKeyId")
		if err != nil {
			return err
		}

		bucket, err := command.FlagOrPromptString(cmd, "bucket")
		if err != nil {
			return err
		}

		endpoint, err := command.FlagOrPromptString(cmd, "endpoint")
		if err != nil {
			return err
		}

		name, err := command.FlagOrPromptString(cmd, "name")
		if err != nil {
			return err
		}

		region, err := command.FlagOrPromptString(cmd, "region")
		if err != nil {
			return err
		}

		secretAccessKey, err := command.FlagOrPromptString(cmd, "secretAccessKey")
		if err != nil {
			return err
		}

		body := idpClient.PostV1OrgOrgNameStorageJSONRequestBody{
			AccessKeyId:     accessKeyId,
			Bucket:          bucket,
			Endpoint:        endpoint,
			Name:            name,
			Region:          region,
			SecretAccessKey: secretAccessKey,
		}

		res, err := client.PostV1OrgOrgNameStorageWithResponse(context.TODO(), orgName, params, body, func(_ context.Context, req *http.Request) error {
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

var storageSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Search storages",
	Long:  "Required scopes: IDP_ORGANIZATION_READ",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
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

		params := &idpClient.PostV1OrgOrgNameStorageSearchParams{}

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

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body idpClient.PostV1OrgOrgNameStorageSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameStorageSearchWithResponse(context.TODO(), orgName, params, body, func(_ context.Context, req *http.Request) error {
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

var storageDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete storage",
	Long:  "Required scopes: IDP_ORGANIZATION_WRITE",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
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

		storageId, err := command.GetArg(cmd, command.StringFlagIfSet("storageId"), command.ContextIfValidString(latest.StorageID), command.PromptString("storageId"))
		if err != nil {
			return err
		}
		latest.StorageID = sql.NullString{Valid: true, String: string(storageId)}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1OrgOrgNameStorageStorageIdWithResponse(context.TODO(), orgName, idpClient.TheIdOfTheStorageToDelete(storageId), func(_ context.Context, req *http.Request) error {
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

var storagePatchCommand = &cobra.Command{
	Use:   "patch",
	Short: "Update a storage",
	Long:  "Required scopes: IDP_ORGANIZATION_WRITE",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
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

		storageId, err := command.GetArg(cmd, command.StringFlagIfSet("storageId"), command.ContextIfValidString(latest.StorageID), command.PromptString("storageId"))
		if err != nil {
			return err
		}
		latest.StorageID = sql.NullString{Valid: true, String: storageId}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		body := idpClient.PatchV1OrgOrgNameStorageStorageIdJSONRequestBody{}

		if cmd.Flags().Changed("accessKeyId") {
			accessKeyId, err := cmd.Flags().GetString("accessKeyId")
			if err != nil {
				return err
			}
			body.AccessKeyId = &accessKeyId
		}

		if cmd.Flags().Changed("bucket") {
			bucket, err := cmd.Flags().GetString("bucket")
			if err != nil {
				return err
			}
			body.Bucket = &bucket
		}

		if cmd.Flags().Changed("endpoint") {
			endpoint, err := cmd.Flags().GetString("endpoint")
			if err != nil {
				return err
			}
			body.Endpoint = &endpoint
		}

		if cmd.Flags().Changed("name") {
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}
			body.Name = &name
		}

		if cmd.Flags().Changed("region") {
			region, err := cmd.Flags().GetString("region")
			if err != nil {
				return err
			}
			body.Region = &region
		}

		if cmd.Flags().Changed("secretAccessKey") {
			secretAccessKey, err := cmd.Flags().GetString("secretAccessKey")
			if err != nil {
				return err
			}
			body.SecretAccessKey = &secretAccessKey
		}

		res, err := client.PatchV1OrgOrgNameStorageStorageIdWithResponse(context.TODO(), orgName, storageId, body, func(_ context.Context, req *http.Request) error {
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
