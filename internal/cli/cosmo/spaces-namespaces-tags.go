package cosmo

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/cosmoClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	spacesNamespacesTagsListCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesTagsListCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNamespacesTagsListCommand.Flags().String("namespaceName", "", "Name of the namespace (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesTagsListCommand.Flags().String("token", "", "authorization token")
	spacesNamespacesTagsCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesTagsCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNamespacesTagsCreateCommand.Flags().String("namespaceName", "", "Name of the namespace (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesTagsCreateCommand.Flags().String("token", "", "authorization token")
	spacesNamespacesTagsCreateCommand.Flags().String("color", "", "Color for the tag")
	spacesNamespacesTagsCreateCommand.Flags().String("name", "", "Name for the tag")
	spacesNamespacesTagsReorderCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesTagsReorderCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNamespacesTagsReorderCommand.Flags().String("namespaceName", "", "Name of the namespace (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesTagsReorderCommand.Flags().String("token", "", "authorization token")
	spacesNamespacesTagsReorderCommand.Flags().String("body", "{}", "request body as JSON")
	spacesNamespacesTagsDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesTagsDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNamespacesTagsDeleteCommand.Flags().String("namespaceName", "", "Name of the namespace (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesTagsDeleteCommand.Flags().String("tagName", "", "Name of the tag (it should contain between 1 and 20 characters)")
	spacesNamespacesTagsDeleteCommand.Flags().String("token", "", "authorization token")
	spacesNamespacesTagsPatchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesTagsPatchCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNamespacesTagsPatchCommand.Flags().String("namespaceName", "", "Name of the namespace (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesTagsPatchCommand.Flags().String("tagName", "", "Name of the tag (it should contain between 1 and 20 characters)")
	spacesNamespacesTagsPatchCommand.Flags().String("token", "", "authorization token")
	spacesNamespacesTagsPatchCommand.Flags().String("color", "", "Color for the tag")
	spacesNamespacesTagsPatchCommand.Flags().String("name", "", "Name for the tag")

	spacesNamespacesTagsCommand.AddCommand(spacesNamespacesTagsListCommand)
	spacesNamespacesTagsCommand.AddCommand(spacesNamespacesTagsCreateCommand)
	spacesNamespacesTagsCommand.AddCommand(spacesNamespacesTagsReorderCommand)
	spacesNamespacesTagsCommand.AddCommand(spacesNamespacesTagsDeleteCommand)
	spacesNamespacesTagsCommand.AddCommand(spacesNamespacesTagsPatchCommand)
}

var spacesNamespacesTagsCommand = &cobra.Command{
	Use:   "tags",
	Short: "Commands for tags",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesNamespacesTagsListCommand = &cobra.Command{
	Use:   "list",
	Short: "List all tags of a namespace",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameTagsWithResponse(context.TODO(), orgName, spaceName, namespaceName, func(_ context.Context, req *http.Request) error {
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

var spacesNamespacesTagsCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a tag",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		color, err := command.FlagOrPromptString(cmd, "color")
		if err != nil {
			return err
		}

		name, err := command.FlagOrPromptString(cmd, "name")
		if err != nil {
			return err
		}

		body := cosmoClient.PostV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameTagsJSONRequestBody{
			Color: color,
			Name:  name,
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameTagsWithResponse(context.TODO(), orgName, spaceName, namespaceName, body, func(_ context.Context, req *http.Request) error {
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

var spacesNamespacesTagsReorderCommand = &cobra.Command{
	Use:   "reorder",
	Short: "Reorder tags of a namespace",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body cosmoClient.PutV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameTagsReorderJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PutV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameTagsReorderWithResponse(context.TODO(), orgName, spaceName, namespaceName, body, func(_ context.Context, req *http.Request) error {
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

var spacesNamespacesTagsDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete a tag",
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

		tagName, err := command.FlagOrPromptString(cmd, "tagName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameTagsTagNameWithResponse(context.TODO(), orgName, spaceName, namespaceName, tagName, func(_ context.Context, req *http.Request) error {
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

var spacesNamespacesTagsPatchCommand = &cobra.Command{
	Use:   "patch",
	Short: "Update a tag",
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

		tagName, err := command.FlagOrPromptString(cmd, "tagName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		body := cosmoClient.PatchV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameTagsTagNameJSONRequestBody{}

		if cmd.Flags().Changed("color") {
			color, err := cmd.Flags().GetString("color")
			if err != nil {
				return err
			}
			body.Color = &color
		}

		if cmd.Flags().Changed("name") {
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}
			body.Name = &name
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameTagsTagNameWithResponse(context.TODO(), orgName, spaceName, namespaceName, tagName, body, func(_ context.Context, req *http.Request) error {
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
