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
	spacesNodesListCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNodesListCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNodesListCommand.Flags().String("token", "", "authorization token")
	spacesNodesListCommand.Flags().Float32("limit", 0, "Page size")
	spacesNodesListCommand.Flags().Float32("page", 0, "Page index")
	spacesNodesCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNodesCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNodesCreateCommand.Flags().String("token", "", "authorization token")
	spacesNodesCreateCommand.Flags().String("category", "", "The category of the node")
	spacesNodesCreateCommand.Flags().String("specification", "", "The specification of the node as a base64 string")
	spacesNodesCreateCommand.Flags().String("typescript", "", "The typescript code content as a base64 string")
	spacesNodesContentCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNodesContentCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNodesContentCommand.Flags().String("secret", "", "")
	spacesNodesContentCommand.Flags().String("token", "", "authorization token")
	spacesNodesDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNodesDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNodesDeleteCommand.Flags().String("nodeId", "", "")
	spacesNodesDeleteCommand.Flags().String("token", "", "authorization token")
	spacesNodesGetCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNodesGetCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNodesGetCommand.Flags().String("nodeId", "", "")
	spacesNodesGetCommand.Flags().String("token", "", "authorization token")
	spacesNodesPatchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNodesPatchCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNodesPatchCommand.Flags().String("nodeId", "", "")
	spacesNodesPatchCommand.Flags().String("token", "", "authorization token")
	spacesNodesPatchCommand.Flags().Bool("regenerateSecret", false, "")
	spacesNodesPatchCommand.Flags().String("specification", "", "The specification of the node as a base64 string")
	spacesNodesPatchCommand.Flags().String("typescript", "", "The typescript code content as a base64 string")

	spacesNodesCommand.AddCommand(spacesNodesListCommand)
	spacesNodesCommand.AddCommand(spacesNodesCreateCommand)
	spacesNodesCommand.AddCommand(spacesNodesContentCommand)
	spacesNodesCommand.AddCommand(spacesNodesDeleteCommand)
	spacesNodesCommand.AddCommand(spacesNodesGetCommand)
	spacesNodesCommand.AddCommand(spacesNodesPatchCommand)
}

var spacesNodesCommand = &cobra.Command{
	Use:   "nodes",
	Short: "Commands for nodes",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesNodesListCommand = &cobra.Command{
	Use:   "list",
	Short: "Search for nodes by space name",
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

		params := &high5Client.GetV1OrgOrgNameSpacesSpaceNameNodesParams{
			Authorization: jwt,
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

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameNodesWithResponse(context.TODO(), orgName, spaceName, params)
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

var spacesNodesCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Add a new node for a space",
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

		category, err := command.FlagOrPromptString(cmd, "category")
		if err != nil {
			return err
		}

		specification, err := command.FlagOrPromptString(cmd, "specification")
		if err != nil {
			return err
		}

		typescript, err := command.FlagOrPromptString(cmd, "typescript")
		if err != nil {
			return err
		}

		body := high5Client.PostV1OrgOrgNameSpacesSpaceNameNodesJSONRequestBody{
			Category:      category,
			Specification: specification,
			Typescript:    typescript,
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameNodesWithResponse(context.TODO(), orgName, spaceName, &high5Client.PostV1OrgOrgNameSpacesSpaceNameNodesParams{
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

var spacesNodesContentCommand = &cobra.Command{
	Use:   "content",
	Short: "Get nodes content as plain text javascript by a secret",
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

		secret, err := command.FlagOrPromptString(cmd, "secret")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameNodesContentSecretWithResponse(context.TODO(), orgName, spaceName, secret, &high5Client.GetV1OrgOrgNameSpacesSpaceNameNodesContentSecretParams{
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

var spacesNodesDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete a node by it's ID",
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

		nodeId, err := command.FlagOrPromptString(cmd, "nodeId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1OrgOrgNameSpacesSpaceNameNodesNodeIdWithResponse(context.TODO(), orgName, spaceName, nodeId, &high5Client.DeleteV1OrgOrgNameSpacesSpaceNameNodesNodeIdParams{
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

var spacesNodesGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Get node by it's ID",
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

		nodeId, err := command.FlagOrPromptString(cmd, "nodeId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameNodesNodeIdWithResponse(context.TODO(), orgName, spaceName, nodeId, &high5Client.GetV1OrgOrgNameSpacesSpaceNameNodesNodeIdParams{
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

var spacesNodesPatchCommand = &cobra.Command{
	Use:   "patch",
	Short: "Update a node by it's ID",
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

		nodeId, err := command.FlagOrPromptString(cmd, "nodeId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &high5Client.PatchV1OrgOrgNameSpacesSpaceNameNodesNodeIdParams{
			Authorization: jwt,
		}

		if cmd.Flags().Changed("regenerateSecret") {
			regenerateSecret, err := cmd.Flags().GetBool("regenerateSecret")
			if err != nil {
				return err
			}
			params.RegenerateSecret = &regenerateSecret
		}

		body := high5Client.PatchV1OrgOrgNameSpacesSpaceNameNodesNodeIdJSONRequestBody{}

		if cmd.Flags().Changed("specification") {
			specification, err := cmd.Flags().GetString("specification")
			if err != nil {
				return err
			}
			body.Specification = &specification
		}

		if cmd.Flags().Changed("typescript") {
			typescript, err := cmd.Flags().GetString("typescript")
			if err != nil {
				return err
			}
			body.Typescript = &typescript
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameNodesNodeIdWithResponse(context.TODO(), orgName, spaceName, nodeId, params, body)
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
