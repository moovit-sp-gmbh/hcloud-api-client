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
	spacesPoolsCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesPoolsCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesPoolsCreateCommand.Flags().String("token", "", "authorization token")
	spacesPoolsCreateCommand.Flags().Bool("noAvatar", false, "")
	spacesPoolsCreateCommand.Flags().String("name", "", "Name of the pool")
	spacesPoolsCreateCommand.Flags().StringSlice("targets", nil, "IDs of the targets of the pool")
	spacesPoolsSearchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesPoolsSearchCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesPoolsSearchCommand.Flags().String("token", "", "authorization token")
	spacesPoolsSearchCommand.Flags().Float32("limit", 0, "Page size")
	spacesPoolsSearchCommand.Flags().Float32("page", 0, "Page index")
	spacesPoolsSearchCommand.Flags().Bool("totalTargetCount", false, "")
	spacesPoolsSearchCommand.Flags().Float32("targetsSample", 0, "Page size")
	spacesPoolsSearchCommand.Flags().String("body", "{}", "request body as JSON")
	spacesPoolsDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesPoolsDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesPoolsDeleteCommand.Flags().String("poolName", "", "Name of the pool (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesPoolsDeleteCommand.Flags().String("token", "", "authorization token")
	spacesPoolsGetCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesPoolsGetCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesPoolsGetCommand.Flags().String("poolName", "", "Name of the pool (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesPoolsGetCommand.Flags().String("token", "", "authorization token")
	spacesPoolsPatchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesPoolsPatchCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesPoolsPatchCommand.Flags().String("poolName", "", "Name of the pool (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesPoolsPatchCommand.Flags().String("token", "", "authorization token")
	spacesPoolsPatchCommand.Flags().String("name", "", "Name of the pool (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesPoolsUpdateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesPoolsUpdateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesPoolsUpdateCommand.Flags().String("poolName", "", "Name of the pool (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesPoolsUpdateCommand.Flags().String("token", "", "authorization token")
	spacesPoolsUpdateCommand.Flags().String("name", "", "Name of the pool")
	spacesPoolsUpdateCommand.Flags().StringSlice("targets", nil, "IDs of the targets of the pool")

	spacesPoolsCommand.AddCommand(spacesPoolsCreateCommand)
	spacesPoolsCommand.AddCommand(spacesPoolsSearchCommand)
	spacesPoolsCommand.AddCommand(spacesPoolsDeleteCommand)
	spacesPoolsCommand.AddCommand(spacesPoolsGetCommand)
	spacesPoolsCommand.AddCommand(spacesPoolsPatchCommand)
	spacesPoolsCommand.AddCommand(spacesPoolsUpdateCommand)
	spacesPoolsCommand.AddCommand(spacesPoolsTargetsCommand)
}

var spacesPoolsCommand = &cobra.Command{
	Use:   "pools",
	Short: "Commands for pools",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesPoolsCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new agent pool",
	Long:  "Create a new agent pool. An agent pool is a pool of execution targets",
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

		params := &high5Client.PostV1OrgOrgNameSpacesSpaceNamePoolsParams{
			Authorization: jwt,
		}

		if cmd.Flags().Changed("noAvatar") {
			noAvatar, err := cmd.Flags().GetBool("noAvatar")
			if err != nil {
				return err
			}
			params.NoAvatar = &noAvatar
		}

		name, err := command.FlagOrPromptString(cmd, "name")
		if err != nil {
			return err
		}

		targets, err := cmd.Flags().GetStringSlice("targets")
		if err != nil {
			return err
		}

		body := high5Client.PostV1OrgOrgNameSpacesSpaceNamePoolsJSONRequestBody{
			Name:    name,
			Targets: targets,
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNamePoolsWithResponse(context.TODO(), orgName, spaceName, params, body)
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

var spacesPoolsSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Search pools of a space",
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

		params := &high5Client.PostV1OrgOrgNameSpacesSpaceNamePoolsSearchParams{
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

		if cmd.Flags().Changed("totalTargetCount") {
			totalTargetCount, err := cmd.Flags().GetBool("totalTargetCount")
			if err != nil {
				return err
			}
			params.TotalTargetCount = &totalTargetCount
		}

		if cmd.Flags().Changed("targetsSample") {
			targetsSample, err := cmd.Flags().GetFloat32("targetsSample")
			if err != nil {
				return err
			}
			params.TargetsSample = &targetsSample
		}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body high5Client.PostV1OrgOrgNameSpacesSpaceNamePoolsSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNamePoolsSearchWithResponse(context.TODO(), orgName, spaceName, params, body)
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

var spacesPoolsDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete an agent pool",
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

		poolName, err := command.FlagOrPromptString(cmd, "poolName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1OrgOrgNameSpacesSpaceNamePoolsPoolNameWithResponse(context.TODO(), orgName, spaceName, poolName, &high5Client.DeleteV1OrgOrgNameSpacesSpaceNamePoolsPoolNameParams{
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

var spacesPoolsGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Fetch a pool by its name",
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

		poolName, err := command.FlagOrPromptString(cmd, "poolName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.GetV1OrgOrgNameSpacesSpaceNamePoolsPoolNameWithResponse(context.TODO(), orgName, spaceName, poolName, &high5Client.GetV1OrgOrgNameSpacesSpaceNamePoolsPoolNameParams{
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

var spacesPoolsPatchCommand = &cobra.Command{
	Use:   "patch",
	Short: "Patch an agent pool",
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

		poolName, err := command.FlagOrPromptString(cmd, "poolName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		body := high5Client.PatchV1OrgOrgNameSpacesSpaceNamePoolsPoolNameJSONRequestBody{}

		if cmd.Flags().Changed("name") {
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}
			body.Name = &name
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNamePoolsPoolNameWithResponse(context.TODO(), orgName, spaceName, poolName, &high5Client.PatchV1OrgOrgNameSpacesSpaceNamePoolsPoolNameParams{
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

var spacesPoolsUpdateCommand = &cobra.Command{
	Use:   "update",
	Short: "Replace an agent pool",
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

		poolName, err := command.FlagOrPromptString(cmd, "poolName")
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

		targets, err := cmd.Flags().GetStringSlice("targets")
		if err != nil {
			return err
		}

		body := high5Client.PutV1OrgOrgNameSpacesSpaceNamePoolsPoolNameJSONRequestBody{
			Name:    name,
			Targets: targets,
		}

		res, err := client.PutV1OrgOrgNameSpacesSpaceNamePoolsPoolNameWithResponse(context.TODO(), orgName, spaceName, poolName, &high5Client.PutV1OrgOrgNameSpacesSpaceNamePoolsPoolNameParams{
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
