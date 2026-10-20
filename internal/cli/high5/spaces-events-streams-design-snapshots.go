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
	spacesEventsStreamsDesignSnapshotsCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesEventsStreamsDesignSnapshotsCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesEventsStreamsDesignSnapshotsCreateCommand.Flags().String("eventName", "", "Name of the event (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesEventsStreamsDesignSnapshotsCreateCommand.Flags().String("streamId", "", "")
	spacesEventsStreamsDesignSnapshotsCreateCommand.Flags().String("token", "", "authorization token")
	spacesEventsStreamsDesignSnapshotsCreateCommand.Flags().String("name", "", "Name of the snapshot (only alphanumeric characters, underscores, hyphens, dots and spaces are allowed; it should contain between 1 and 64 characters)")
	spacesEventsStreamsDesignSnapshotsSearchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesEventsStreamsDesignSnapshotsSearchCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesEventsStreamsDesignSnapshotsSearchCommand.Flags().String("eventName", "", "Name of the event (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesEventsStreamsDesignSnapshotsSearchCommand.Flags().String("streamId", "", "")
	spacesEventsStreamsDesignSnapshotsSearchCommand.Flags().String("token", "", "authorization token")
	spacesEventsStreamsDesignSnapshotsSearchCommand.Flags().Bool("content", false, "")
	spacesEventsStreamsDesignSnapshotsSearchCommand.Flags().Float32("limit", 0, "Page size")
	spacesEventsStreamsDesignSnapshotsSearchCommand.Flags().Float32("page", 0, "Page index")
	spacesEventsStreamsDesignSnapshotsSearchCommand.Flags().String("body", "{}", "request body as JSON")
	spacesEventsStreamsDesignSnapshotsDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesEventsStreamsDesignSnapshotsDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesEventsStreamsDesignSnapshotsDeleteCommand.Flags().String("eventName", "", "Name of the event (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesEventsStreamsDesignSnapshotsDeleteCommand.Flags().String("streamId", "", "")
	spacesEventsStreamsDesignSnapshotsDeleteCommand.Flags().String("snapshotId", "", "")
	spacesEventsStreamsDesignSnapshotsDeleteCommand.Flags().String("token", "", "authorization token")
	spacesEventsStreamsDesignSnapshotsPatchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesEventsStreamsDesignSnapshotsPatchCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesEventsStreamsDesignSnapshotsPatchCommand.Flags().String("eventName", "", "Name of the event (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesEventsStreamsDesignSnapshotsPatchCommand.Flags().String("streamId", "", "")
	spacesEventsStreamsDesignSnapshotsPatchCommand.Flags().String("snapshotId", "", "")
	spacesEventsStreamsDesignSnapshotsPatchCommand.Flags().String("token", "", "authorization token")
	spacesEventsStreamsDesignSnapshotsPatchCommand.Flags().String("name", "", "Name of the snapshot (only alphanumeric characters, underscores, hyphens, dots and spaces are allowed; it should contain between 1 and 64 characters)")
	spacesEventsStreamsDesignSnapshotsRestoreCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesEventsStreamsDesignSnapshotsRestoreCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesEventsStreamsDesignSnapshotsRestoreCommand.Flags().String("eventName", "", "Name of the event (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesEventsStreamsDesignSnapshotsRestoreCommand.Flags().String("streamId", "", "")
	spacesEventsStreamsDesignSnapshotsRestoreCommand.Flags().String("snapshotId", "", "")
	spacesEventsStreamsDesignSnapshotsRestoreCommand.Flags().String("token", "", "authorization token")

	spacesEventsStreamsDesignSnapshotsCommand.AddCommand(spacesEventsStreamsDesignSnapshotsCreateCommand)
	spacesEventsStreamsDesignSnapshotsCommand.AddCommand(spacesEventsStreamsDesignSnapshotsSearchCommand)
	spacesEventsStreamsDesignSnapshotsCommand.AddCommand(spacesEventsStreamsDesignSnapshotsDeleteCommand)
	spacesEventsStreamsDesignSnapshotsCommand.AddCommand(spacesEventsStreamsDesignSnapshotsPatchCommand)
	spacesEventsStreamsDesignSnapshotsCommand.AddCommand(spacesEventsStreamsDesignSnapshotsRestoreCommand)
}

var spacesEventsStreamsDesignSnapshotsCommand = &cobra.Command{
	Use:   "snapshots",
	Short: "Commands for snapshots",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesEventsStreamsDesignSnapshotsCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new snapshot of a design",
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

		body := high5Client.PostV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignSnapshotsJSONRequestBody{
			Name: name,
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignSnapshotsWithResponse(context.TODO(), orgName, spaceName, eventName, streamId, &high5Client.PostV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignSnapshotsParams{
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

var spacesEventsStreamsDesignSnapshotsSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Search the snapshots of a design",
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

		params := &high5Client.PostV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignSnapshotsSearchParams{
			Authorization: jwt,
		}

		if cmd.Flags().Changed("content") {
			content, err := cmd.Flags().GetBool("content")
			if err != nil {
				return err
			}
			params.Content = &content
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
		var body high5Client.PostV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignSnapshotsSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignSnapshotsSearchWithResponse(context.TODO(), orgName, spaceName, eventName, streamId, params, body)
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

var spacesEventsStreamsDesignSnapshotsDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete a snapshot of a design",
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

		snapshotId, err := command.FlagOrPromptString(cmd, "snapshotId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignSnapshotsSnapshotIdWithResponse(context.TODO(), orgName, spaceName, eventName, streamId, snapshotId, &high5Client.DeleteV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignSnapshotsSnapshotIdParams{
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

var spacesEventsStreamsDesignSnapshotsPatchCommand = &cobra.Command{
	Use:   "patch",
	Short: "Update a snapshot",
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

		snapshotId, err := command.FlagOrPromptString(cmd, "snapshotId")
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

		body := high5Client.PatchV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignSnapshotsSnapshotIdJSONRequestBody{
			Name: name,
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignSnapshotsSnapshotIdWithResponse(context.TODO(), orgName, spaceName, eventName, streamId, snapshotId, &high5Client.PatchV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignSnapshotsSnapshotIdParams{
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

var spacesEventsStreamsDesignSnapshotsRestoreCommand = &cobra.Command{
	Use:   "restore",
	Short: "Restore the design to this snapshot",
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

		snapshotId, err := command.FlagOrPromptString(cmd, "snapshotId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.PutV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignSnapshotsSnapshotIdRestoreWithResponse(context.TODO(), orgName, spaceName, eventName, streamId, snapshotId, &high5Client.PutV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignSnapshotsSnapshotIdRestoreParams{
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
