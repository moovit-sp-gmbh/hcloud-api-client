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
	executionLogsSearchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	executionLogsSearchCommand.Flags().String("token", "", "authorization token")
	executionLogsSearchCommand.Flags().Float32("limit", 0, "Page size")
	executionLogsSearchCommand.Flags().Float32("page", 0, "Page index")
	executionLogsSearchCommand.Flags().Bool("payload", false, "Optional parameter to include the payload in the response")
	executionLogsSearchCommand.Flags().String("body", "{}", "request body as JSON")

	executionLogsCommand.AddCommand(executionLogsSearchCommand)
	executionCommand.AddCommand(executionLogsCommand)
	executionCommand.AddCommand(executionStatusCommand)
}

var executionCommand = &cobra.Command{
	Use:   "execution",
	Short: "Commands for execution",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var executionLogsCommand = &cobra.Command{
	Use:   "logs",
	Short: "Commands for logs",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var executionLogsSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Retrieve the Stream execution logs of an Organization",
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

		params := &high5Client.PostV1OrgOrgNameExecutionLogsSearchParams{
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

		if cmd.Flags().Changed("payload") {
			payload, err := cmd.Flags().GetBool("payload")
			if err != nil {
				return err
			}
			params.Payload = &payload
		}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body high5Client.PostV1OrgOrgNameExecutionLogsSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameExecutionLogsSearchWithResponse(context.TODO(), orgName, params, body)
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
