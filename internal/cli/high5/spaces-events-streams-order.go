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
	spacesEventsStreamsOrderPatchAllCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesEventsStreamsOrderPatchAllCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesEventsStreamsOrderPatchAllCommand.Flags().String("eventName", "", "Name of the event (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesEventsStreamsOrderPatchAllCommand.Flags().String("token", "", "authorization token")
	spacesEventsStreamsOrderPatchAllCommand.Flags().String("body", "{}", "request body as JSON")
	spacesEventsStreamsOrderPatchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesEventsStreamsOrderPatchCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesEventsStreamsOrderPatchCommand.Flags().String("eventName", "", "Name of the event (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesEventsStreamsOrderPatchCommand.Flags().String("streamId", "", "")
	spacesEventsStreamsOrderPatchCommand.Flags().String("token", "", "authorization token")
	spacesEventsStreamsOrderPatchCommand.Flags().String("direction", "", "Direction of the changing of stream order")

	spacesEventsStreamsOrderCommand.AddCommand(spacesEventsStreamsOrderPatchAllCommand)
	spacesEventsStreamsOrderCommand.AddCommand(spacesEventsStreamsOrderPatchCommand)
}

var spacesEventsStreamsOrderCommand = &cobra.Command{
	Use:   "order",
	Short: "Commands for order",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesEventsStreamsOrderPatchAllCommand = &cobra.Command{
	Use:   "patch-all",
	Short: "Update the stream order for an event",
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

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body high5Client.PatchV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsOrderJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsOrderWithResponse(context.TODO(), orgName, spaceName, eventName, &high5Client.PatchV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsOrderParams{
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

var spacesEventsStreamsOrderPatchCommand = &cobra.Command{
	Use:   "patch",
	Short: "Update the stream order for a single stream",
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

		direction, err := command.FlagOrPromptString(cmd, "direction")
		if err != nil {
			return err
		}

		body := high5Client.PatchV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdOrderJSONRequestBody{
			Direction: high5Client.MoveDirection(direction),
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdOrderWithResponse(context.TODO(), orgName, spaceName, eventName, streamId, &high5Client.PatchV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdOrderParams{
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
