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
	spacesEventsStreamsDesignListCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesEventsStreamsDesignListCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesEventsStreamsDesignListCommand.Flags().String("eventName", "", "Name of the event (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesEventsStreamsDesignListCommand.Flags().String("streamId", "", "")
	spacesEventsStreamsDesignListCommand.Flags().String("token", "", "authorization token")
	spacesEventsStreamsDesignUpdateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesEventsStreamsDesignUpdateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesEventsStreamsDesignUpdateCommand.Flags().String("eventName", "", "Name of the event (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesEventsStreamsDesignUpdateCommand.Flags().String("streamId", "", "")
	spacesEventsStreamsDesignUpdateCommand.Flags().String("token", "", "authorization token")
	spacesEventsStreamsDesignUpdateCommand.Flags().String("body", "{}", "request body as JSON")
	spacesEventsStreamsDesignOperationsCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesEventsStreamsDesignOperationsCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesEventsStreamsDesignOperationsCommand.Flags().String("eventName", "", "Name of the event (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesEventsStreamsDesignOperationsCommand.Flags().String("streamId", "", "")
	spacesEventsStreamsDesignOperationsCommand.Flags().String("token", "", "authorization token")
	spacesEventsStreamsDesignOperationsCommand.Flags().String("body", "{}", "request body as JSON")
	spacesEventsStreamsDesignPublishCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesEventsStreamsDesignPublishCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesEventsStreamsDesignPublishCommand.Flags().String("eventName", "", "Name of the event (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesEventsStreamsDesignPublishCommand.Flags().String("streamId", "", "")
	spacesEventsStreamsDesignPublishCommand.Flags().String("token", "", "authorization token")

	spacesEventsStreamsDesignCommand.AddCommand(spacesEventsStreamsDesignListCommand)
	spacesEventsStreamsDesignCommand.AddCommand(spacesEventsStreamsDesignUpdateCommand)
	spacesEventsStreamsDesignCommand.AddCommand(spacesEventsStreamsDesignOperationsCommand)
	spacesEventsStreamsDesignCommand.AddCommand(spacesEventsStreamsDesignPublishCommand)
	spacesEventsStreamsDesignCommand.AddCommand(spacesEventsStreamsDesignSnapshotsCommand)
}

var spacesEventsStreamsDesignCommand = &cobra.Command{
	Use:   "design",
	Short: "Commands for design",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesEventsStreamsDesignListCommand = &cobra.Command{
	Use:   "list",
	Short: "Get a stream's design",
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

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignWithResponse(context.TODO(), orgName, spaceName, eventName, streamId, &high5Client.GetV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignParams{
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

var spacesEventsStreamsDesignUpdateCommand = &cobra.Command{
	Use:   "update",
	Short: "Create a new design for a stream",
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

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body high5Client.PutV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PutV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignWithResponse(context.TODO(), orgName, spaceName, eventName, streamId, &high5Client.PutV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignParams{
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

var spacesEventsStreamsDesignOperationsCommand = &cobra.Command{
	Use:   "operations",
	Short: "Apply a set of operations to a Design",
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

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body high5Client.PatchV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignOperationsJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignOperationsWithResponse(context.TODO(), orgName, spaceName, eventName, streamId, &high5Client.PatchV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignOperationsParams{
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

var spacesEventsStreamsDesignPublishCommand = &cobra.Command{
	Use:   "publish",
	Short: "Publish the design of a stream",
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

		res, err := client.PutV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignPublishWithResponse(context.TODO(), orgName, spaceName, eventName, streamId, &high5Client.PutV1OrgOrgNameSpacesSpaceNameEventsEventNameStreamsStreamIdDesignPublishParams{
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
