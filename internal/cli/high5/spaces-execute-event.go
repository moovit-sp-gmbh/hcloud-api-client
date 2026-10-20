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
	spacesExecuteEventNameCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesExecuteEventNameCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesExecuteEventNameCommand.Flags().String("eventName", "", "Name of the event (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 128 characters)")
	spacesExecuteEventNameCommand.Flags().String("token", "", "authorization token")
	spacesExecuteEventNameCommand.Flags().String("body", "{}", "request body as JSON")

	spacesExecuteEventCommand.AddCommand(spacesExecuteEventNameCommand)
	spacesExecuteCommand.AddCommand(spacesExecuteEventCommand)
	spacesExecuteCommand.AddCommand(spacesExecuteStreamCommand)
	spacesExecuteCommand.AddCommand(spacesExecuteWebhookCommand)
	spacesExecuteCommand.AddCommand(spacesExecuteDebugCommand)
}

var spacesExecuteCommand = &cobra.Command{
	Use:   "execute",
	Short: "Commands for execute",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesExecuteEventCommand = &cobra.Command{
	Use:   "event",
	Short: "Commands for event",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesExecuteEventNameCommand = &cobra.Command{
	Use:   "name",
	Short: "Execute all streams connected to an event by the event name",
	Long:  "Creates execution entries for all active streams of the event and dispatches the first stream to the target agent; each subsequent stream is triggered automatically once the agent reports successful completion of the previous one",
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
		var body high5Client.PostV1OrgOrgNameSpacesSpaceNameExecuteEventNameEventNameJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameExecuteEventNameEventNameWithResponse(context.TODO(), orgName, spaceName, eventName, &high5Client.PostV1OrgOrgNameSpacesSpaceNameExecuteEventNameEventNameParams{
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
