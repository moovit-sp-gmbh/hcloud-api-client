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
	spacesCommand.AddCommand(spacesSearchCommand)
	spacesCommand.AddCommand(spacesDeleteCommand)
	spacesCommand.AddCommand(spacesDatabasesCommand)
	spacesCommand.AddCommand(spacesEventsCommand)
	spacesCommand.AddCommand(spacesExecuteCommand)
	spacesCommand.AddCommand(spacesExecutionCommand)
	spacesCommand.AddCommand(spacesJobsCommand)
	spacesCommand.AddCommand(spacesNameCommand)
	spacesCommand.AddCommand(spacesNodesCommand)
	spacesCommand.AddCommand(spacesPermissionsCommand)
	spacesCommand.AddCommand(spacesPoolsCommand)
	spacesCommand.AddCommand(spacesRequestCommand)
	spacesCommand.AddCommand(spacesSecretsCommand)
	spacesCommand.AddCommand(spacesStatsCommand)
	spacesCommand.AddCommand(spacesStreamsCommand)
	spacesCommand.AddCommand(spacesTargetCommand)
	spacesCommand.AddCommand(spacesWatchfoldersCommand)
	spacesCommand.AddCommand(spacesWaveCommand)
	spacesCommand.AddCommand(spacesWebhooksCommand)
	spacesCommand.AddCommand(spacesListCommand)
	spaceGetCommand.Flags().String("orgName", "", "")
	spaceGetCommand.Flags().String("spaceName", "", "")
	spaceGetCommand.Flags().String("token", "", "")
	spaceCreateCommand.Flags().String("orgName", "", "Name of the organization")
	spaceCreateCommand.Flags().String("token", "", "Authorization token to use in the API call")
	spaceCreateCommand.Flags().Bool("noAvatar", false, "Skip avatar for space")
	spaceCreateCommand.Flags().String("name", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spaceCreateCommand.Flags().Bool("enableDev", false, "Specifies to the front-end whether to display the versions under development of Wave Engines and Node Catalogs")
	spacesSearchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesSearchCommand.Flags().String("token", "", "authorization token")
	spacesSearchCommand.Flags().Float32("limit", 0, "Page size")
	spacesSearchCommand.Flags().Float32("page", 0, "Page index")
	spacesSearchCommand.Flags().String("body", "{}", "request body as JSON")
	spacesDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesDeleteCommand.Flags().String("token", "", "authorization token")
	spacesNameCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNameCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNameCommand.Flags().String("token", "", "authorization token")
	spacesNameCommand.Flags().String("newSpaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesStatsCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesStatsCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesStatsCommand.Flags().String("token", "", "authorization token")
	spacesTargetCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesTargetCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesTargetCommand.Flags().String("token", "", "authorization token")
	spacesTargetCommand.Flags().String("body", "{}", "request body as JSON")
	spacesListCommand.Flags().String("token", "", "authorization token")
	spaceCommand.AddCommand(spaceCreateCommand)
	spaceCommand.AddCommand(spaceEditCommand)
	spaceCommand.AddCommand(spaceGetCommand)
}

var spaceCommand = &cobra.Command{
	Use:   "space",
	Short: "Commands for space",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spaceCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new space in High5",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &high5Client.PostV1OrgOrgNameSpacesParams{
			Authorization: jwt,
		}

		if cmd.Flags().Changed("noAvatar") {
			noAvatar, err := cmd.Flags().GetBool("noAvatar")
			if err != nil {
				return err
			}
			params.NoAvatar = &noAvatar
		}

		name, err := command.GetArg(cmd, command.StringFlagIfSet("name"), command.ContextIfValidString(latest.High5SpaceName), command.PromptString("name"))
		if err != nil {
			return err
		}
		latest.High5SpaceName = sql.NullString{Valid: true, String: name}

		body := high5Client.PostV1OrgOrgNameSpacesJSONRequestBody{
			Name: name,
		}

		if cmd.Flags().Changed("enableDev") {
			enableDev, err := cmd.Flags().GetBool("enableDev")
			if err != nil {
				return err
			}
			body.EnableDev = &enableDev
		}

		res, err := client.PostV1OrgOrgNameSpacesWithResponse(context.TODO(), orgName, params, body)
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
			if res.JSON201 != nil {
				latest.High5SpaceID = sql.NullString{Valid: true, String: res.JSON201.UnderscoreId}
			}
			if res.JSON201 != nil {
				latest.High5SpaceName = sql.NullString{Valid: true, String: res.JSON201.Name}
			}
			service.Push(latest)
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var spaceGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Get a High5 space by its name",
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

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("organization name"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		spaceName, err := command.GetArg(cmd, command.StringFlagIfSet("spaceName"), command.ContextIfValidString(latest.High5SpaceName), command.PromptString("space name"))
		if err != nil {
			return err
		}
		latest.High5SpaceName = sql.NullString{Valid: true, String: spaceName}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameWithResponse(context.TODO(), orgName, spaceName, &high5Client.GetV1OrgOrgNameSpacesSpaceNameParams{
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

var spacesCommand = &cobra.Command{
	Use:   "spaces",
	Short: "Commands for spaces",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Search spaces in an organization",
	Long:  "Returns spaces in the organization matching the search filter, scoped to the requesting user's access, with total count in the 'total' header",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &high5Client.PostV1OrgOrgNameSpacesSearchParams{
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

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body high5Client.PostV1OrgOrgNameSpacesSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameSpacesSearchWithResponse(context.TODO(), orgName, params, body)
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

var spacesDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete a space and its corresponding data by the space's name",
	Long:  "Deletes the space and all of its associated data",
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

		res, err := client.DeleteV1OrgOrgNameSpacesSpaceNameWithResponse(context.TODO(), orgName, spaceName, &high5Client.DeleteV1OrgOrgNameSpacesSpaceNameParams{
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

var spacesNameCommand = &cobra.Command{
	Use:   "name",
	Short: "Rename a Space",
	Long:  "Renames the space after verifying the new name is not already in use within the organization",
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

		newSpaceName, err := command.FlagOrPromptString(cmd, "newSpaceName")
		if err != nil {
			return err
		}

		body := high5Client.PatchV1OrgOrgNameSpacesSpaceNameNameJSONRequestBody{
			NewSpaceName: newSpaceName,
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameNameWithResponse(context.TODO(), orgName, spaceName, &high5Client.PatchV1OrgOrgNameSpacesSpaceNameNameParams{
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

var spacesStatsCommand = &cobra.Command{
	Use:   "stats",
	Short: "Retrieve usage statistics for a specific space",
	Long:  "Provides detailed license usage metrics for a given space",
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

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameStatsWithResponse(context.TODO(), orgName, spaceName, &high5Client.GetV1OrgOrgNameSpacesSpaceNameStatsParams{
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

var spacesTargetCommand = &cobra.Command{
	Use:   "target",
	Short: "Set the execution target used when this Space fires its own internal events",
	Long:  "Validates that the given target resolves to an existing pool or an execution-target-flagged org member, then sets it as the default execution target for this space's internal hcloud.high5.* events",
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

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body high5Client.PatchV1OrgOrgNameSpacesSpaceNameTargetJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameTargetWithResponse(context.TODO(), orgName, spaceName, &high5Client.PatchV1OrgOrgNameSpacesSpaceNameTargetParams{
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

var spacesListCommand = &cobra.Command{
	Use: "list",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.GetV1SpacesWithResponse(context.TODO(), &high5Client.GetV1SpacesParams{
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
