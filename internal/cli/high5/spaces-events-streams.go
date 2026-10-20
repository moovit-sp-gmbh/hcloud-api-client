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
	spacesEventsStreamsCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesEventsStreamsCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesEventsStreamsCreateCommand.Flags().String("eventName", "", "Name of the event (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesEventsStreamsCreateCommand.Flags().String("token", "", "authorization token")
	spacesEventsStreamsCreateCommand.Flags().String("name", "", "Name of the stream, dots and spaces are allowed only inside a name (only alphanumeric characters, underscores, hyphens, dots and spaces are allowed; it should contain between 1 and 128 characters)")
	spacesEventsStreamsCreateCommand.Flags().String("streamId", "", "ID of the stream which will be duplicated")
	spacesEventsStreamsSearchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesEventsStreamsSearchCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesEventsStreamsSearchCommand.Flags().String("eventName", "", "Name of the event (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesEventsStreamsSearchCommand.Flags().String("token", "", "authorization token")
	spacesEventsStreamsSearchCommand.Flags().Float32("limit", 0, "Page size")
	spacesEventsStreamsSearchCommand.Flags().Float32("page", 0, "Page index")
	spacesEventsStreamsSearchCommand.Flags().String("body", "{}", "request body as JSON")
	spacesEventsStreamsDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesEventsStreamsDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesEventsStreamsDeleteCommand.Flags().String("eventName", "", "Name of the event (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesEventsStreamsDeleteCommand.Flags().String("streamId", "", "")
	spacesEventsStreamsDeleteCommand.Flags().String("token", "", "authorization token")
	spacesEventsStreamsGetCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesEventsStreamsGetCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesEventsStreamsGetCommand.Flags().String("eventName", "", "Name of the event (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesEventsStreamsGetCommand.Flags().String("streamId", "", "")
	spacesEventsStreamsGetCommand.Flags().String("token", "", "authorization token")
	spacesEventsStreamsActiveCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesEventsStreamsActiveCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesEventsStreamsActiveCommand.Flags().String("eventName", "", "Name of the event (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesEventsStreamsActiveCommand.Flags().String("streamId", "", "")
	spacesEventsStreamsActiveCommand.Flags().String("token", "", "authorization token")
	spacesEventsStreamsActiveCommand.Flags().Bool("active", false, "New state of the stream")
	spacesEventsStreamsEventCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesEventsStreamsEventCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesEventsStreamsEventCommand.Flags().String("eventName", "", "Name of the event (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesEventsStreamsEventCommand.Flags().String("streamId", "", "")
	spacesEventsStreamsEventCommand.Flags().String("token", "", "authorization token")
	spacesEventsStreamsEventCommand.Flags().String("bodyEventName", "", "Name of the event (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesEventsStreamsEventCommand.Flags().Float32("order", 0, "Optional order to set the stream in the new event. If omitted, the stream will be set last in the new event.")
	spacesEventsStreamsNameCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesEventsStreamsNameCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesEventsStreamsNameCommand.Flags().String("eventName", "", "Name of the event (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesEventsStreamsNameCommand.Flags().String("streamId", "", "")
	spacesEventsStreamsNameCommand.Flags().String("token", "", "authorization token")
	spacesEventsStreamsNameCommand.Flags().String("name", "", "Name of the stream, dots and spaces are allowed only inside a name (only alphanumeric characters, underscores, hyphens, dots and spaces are allowed; it should contain between 1 and 128 characters)")

	spacesEventsStreamsCommand.AddCommand(spacesEventsStreamsCreateCommand)
	spacesEventsStreamsCommand.AddCommand(spacesEventsStreamsOrderCommand)
	spacesEventsStreamsCommand.AddCommand(spacesEventsStreamsSearchCommand)
	spacesEventsStreamsCommand.AddCommand(spacesEventsStreamsDeleteCommand)
	spacesEventsStreamsCommand.AddCommand(spacesEventsStreamsGetCommand)
	spacesEventsStreamsCommand.AddCommand(spacesEventsStreamsActiveCommand)
	spacesEventsStreamsCommand.AddCommand(spacesEventsStreamsDesignCommand)
	spacesEventsStreamsCommand.AddCommand(spacesEventsStreamsEventCommand)
	spacesEventsStreamsCommand.AddCommand(spacesEventsStreamsNameCommand)
}

var spacesEventsStreamsCommand = &cobra.Command{
	Use:   "streams",
	Short: "Commands for streams",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesEventsStreamsCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Add a new stream for an event or duplicate an existing stream (if streamId is provided)",
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

		eventName, err := command.FlagOrPromptString(cmd, "eventName")
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

		body := high5Client.PostV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsJSONRequestBody{
			Name: name,
		}

		if cmd.Flags().Changed("streamId") {
			streamId, err := cmd.Flags().GetString("streamId")
			if err != nil {
				return err
			}
			body.StreamId = &streamId
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsWithResponse(context.TODO(), orgName, spaceName, eventName, &high5Client.PostV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsParams{
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

var spacesEventsStreamsSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Search for streams by event name and search filter",
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

		eventName, err := command.FlagOrPromptString(cmd, "eventName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &high5Client.PostV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsSearchParams{
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
		var body high5Client.PostV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsSearchWithResponse(context.TODO(), orgName, spaceName, eventName, params, body)
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

var spacesEventsStreamsDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete a stream by it's ID",
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

		eventName, err := command.FlagOrPromptString(cmd, "eventName")
		if err != nil {
			return err
		}

		streamId, err := command.FlagOrPromptString(cmd, "streamId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdWithResponse(context.TODO(), orgName, spaceName, eventName, streamId, &high5Client.DeleteV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdParams{
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

var spacesEventsStreamsGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Get stream by it's ID",
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

		eventName, err := command.FlagOrPromptString(cmd, "eventName")
		if err != nil {
			return err
		}

		streamId, err := command.FlagOrPromptString(cmd, "streamId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdWithResponse(context.TODO(), orgName, spaceName, eventName, streamId, &high5Client.GetV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdParams{
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

var spacesEventsStreamsActiveCommand = &cobra.Command{
	Use:   "active",
	Short: "Activate or deactivate a stream",
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

		eventName, err := command.FlagOrPromptString(cmd, "eventName")
		if err != nil {
			return err
		}

		streamId, err := command.FlagOrPromptString(cmd, "streamId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		active, err := command.FlagOrPromptBool(cmd, "active")
		if err != nil {
			return err
		}

		body := high5Client.PatchV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdActiveJSONRequestBody{
			Active: active,
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdActiveWithResponse(context.TODO(), orgName, spaceName, eventName, streamId, &high5Client.PatchV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdActiveParams{
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

var spacesEventsStreamsEventCommand = &cobra.Command{
	Use:   "event",
	Short: "Move a stream to another event",
	Long:  "Move a stream from one event to another",
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

		eventName, err := command.FlagOrPromptString(cmd, "eventName")
		if err != nil {
			return err
		}

		streamId, err := command.FlagOrPromptString(cmd, "streamId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		bodyEventName, err := command.FlagOrDifferentPromptString(cmd, "bodyEventName", "eventName")
		if err != nil {
			return err
		}

		body := high5Client.PatchV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdEventJSONRequestBody{
			EventName: bodyEventName,
		}

		if cmd.Flags().Changed("order") {
			order, err := cmd.Flags().GetFloat32("order")
			if err != nil {
				return err
			}
			body.Order = &order
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdEventWithResponse(context.TODO(), orgName, spaceName, eventName, streamId, &high5Client.PatchV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdEventParams{
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

var spacesEventsStreamsNameCommand = &cobra.Command{
	Use:   "name",
	Short: "Change a stream's name",
	Long:  "Update the name of a stream by providing a new one that is not already assigned to another stream within the space",
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

		eventName, err := command.FlagOrPromptString(cmd, "eventName")
		if err != nil {
			return err
		}

		streamId, err := command.FlagOrPromptString(cmd, "streamId")
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

		body := high5Client.PatchV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdNameJSONRequestBody{
			Name: name,
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdNameWithResponse(context.TODO(), orgName, spaceName, eventName, streamId, &high5Client.PatchV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdNameParams{
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
