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
	spacesExecutionStreamsExecutionTreeCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesExecutionStreamsExecutionTreeCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesExecutionStreamsExecutionTreeCommand.Flags().String("streamStatusId", "", "")
	spacesExecutionStreamsExecutionTreeCommand.Flags().String("token", "", "authorization token")
	spacesExecutionStreamsLogsCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesExecutionStreamsLogsCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesExecutionStreamsLogsCommand.Flags().String("streamLogId", "", "")
	spacesExecutionStreamsLogsCommand.Flags().String("token", "", "authorization token")
	spacesExecutionStreamsLogsCommand.Flags().Bool("payload", false, "Optional parameter to include the payload in the response")
	spacesExecutionStreamsStatusCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesExecutionStreamsStatusCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesExecutionStreamsStatusCommand.Flags().String("streamStatusId", "", "")
	spacesExecutionStreamsStatusCommand.Flags().String("token", "", "authorization token")
	spacesExecutionStreamsPatchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesExecutionStreamsPatchCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesExecutionStreamsPatchCommand.Flags().String("streamSecret", "", "Secret that got generated on execution package creation at one of the POST execute endpoints (it should contain between 3 and 1024 characters)")
	spacesExecutionStreamsPatchCommand.Flags().String("token", "", "authorization token")
	spacesExecutionStreamsPatchCommand.Flags().String("body", "{}", "request body as JSON")
	spacesExecutionStreamsCrashCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesExecutionStreamsCrashCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesExecutionStreamsCrashCommand.Flags().String("streamSecret", "", "Secret that got generated on execution package creation at one of the POST execute endpoints (it should contain between 3 and 1024 characters)")
	spacesExecutionStreamsCrashCommand.Flags().String("token", "", "authorization token")
	spacesExecutionStreamsCrashCommand.Flags().String("message", "", "Crash log")

	spacesExecutionStreamsCommand.AddCommand(spacesExecutionStreamsExecutionTreeCommand)
	spacesExecutionStreamsCommand.AddCommand(spacesExecutionStreamsLogsCommand)
	spacesExecutionStreamsCommand.AddCommand(spacesExecutionStreamsStatusCommand)
	spacesExecutionStreamsCommand.AddCommand(spacesExecutionStreamsPatchCommand)
	spacesExecutionStreamsCommand.AddCommand(spacesExecutionStreamsCrashCommand)
}

var spacesExecutionStreamsCommand = &cobra.Command{
	Use:   "streams",
	Short: "Commands for streams",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesExecutionStreamsExecutionTreeCommand = &cobra.Command{
	Use:   "execution-tree",
	Short: "Retrieve the complete executionTree the stream execution is part of. An executionTree represents the hierarchical structure of stream executions initiated by a 'Trigger Stream' node",
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

		streamStatusId, err := command.FlagOrPromptString(cmd, "streamStatusId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameExecutionStreamsExecutionTreeStreamStatusIdWithResponse(context.TODO(), orgName, spaceName, streamStatusId, &high5Client.GetV1OrgOrgNameSpacesSpaceNameExecutionStreamsExecutionTreeStreamStatusIdParams{
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

var spacesExecutionStreamsLogsCommand = &cobra.Command{
	Use:   "logs",
	Short: "Retrieve a Stream execution log by its ID",
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

		streamLogId, err := command.FlagOrPromptString(cmd, "streamLogId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &high5Client.GetV1OrgOrgNameSpacesSpaceNameExecutionStreamsLogsStreamLogIdParams{
			Authorization: jwt,
		}

		if cmd.Flags().Changed("payload") {
			payload, err := cmd.Flags().GetBool("payload")
			if err != nil {
				return err
			}
			params.Payload = &payload
		}

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameExecutionStreamsLogsStreamLogIdWithResponse(context.TODO(), orgName, spaceName, streamLogId, params)
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

var spacesExecutionStreamsStatusCommand = &cobra.Command{
	Use:   "status",
	Short: "Retrieve a Stream execution status by its ID",
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

		streamStatusId, err := command.FlagOrPromptString(cmd, "streamStatusId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameExecutionStreamsStatusStreamStatusIdWithResponse(context.TODO(), orgName, spaceName, streamStatusId, &high5Client.GetV1OrgOrgNameSpacesSpaceNameExecutionStreamsStatusStreamStatusIdParams{
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

var spacesExecutionStreamsPatchCommand = &cobra.Command{
	Use:   "patch",
	Short: "Update stream log and status during execution or after execution is finished",
	Long:  "Live update stream log and status during execution or after execution is finished",
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

		streamSecret, err := command.FlagOrPromptString(cmd, "streamSecret")
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
		var body high5Client.PatchV1OrgOrgNameSpacesSpaceNameExecutionStreamsStreamSecretJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameExecutionStreamsStreamSecretWithResponse(context.TODO(), orgName, spaceName, streamSecret, &high5Client.PatchV1OrgOrgNameSpacesSpaceNameExecutionStreamsStreamSecretParams{
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

var spacesExecutionStreamsCrashCommand = &cobra.Command{
	Use:   "crash",
	Short: "Mark a stream has failed with an optional message",
	Long:  "Mark a stream has failed with a message. To be used on for some reason regular logging is not possible",
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

		streamSecret, err := command.FlagOrPromptString(cmd, "streamSecret")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		body := high5Client.PatchV1OrgOrgNameSpacesSpaceNameExecutionStreamsStreamSecretCrashJSONRequestBody{}

		if cmd.Flags().Changed("message") {
			message, err := cmd.Flags().GetString("message")
			if err != nil {
				return err
			}
			body.Message = &message
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameExecutionStreamsStreamSecretCrashWithResponse(context.TODO(), orgName, spaceName, streamSecret, &high5Client.PatchV1OrgOrgNameSpacesSpaceNameExecutionStreamsStreamSecretCrashParams{
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
