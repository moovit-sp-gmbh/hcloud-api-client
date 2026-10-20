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
	executeQueueCommand.Flags().String("token", "", "authorization token")
	executeQueueCommand.Flags().Bool("blocking", false, "Whether this agent is ready to accept blocking executions")
	executeQueueCommand.Flags().StringSlice("organizationIds", nil, "IDs of the organizations to poll the queue for")
	executeQueueCommand.Flags().Int("limit", 0, "Maximum number of executions to return")

	executeCommand.AddCommand(executeQueueCommand)
}

var executeCommand = &cobra.Command{
	Use:   "execute",
	Short: "Commands for execute",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var executeQueueCommand = &cobra.Command{
	Use:   "queue",
	Short: "Poll the global execution queue",
	Long:  "Fetches pending queued executions across all provided organizations, dispatches each one to the calling agent via NATS and waits until all have been accepted or failed",
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

		blocking, err := command.FlagOrPromptBool(cmd, "blocking")
		if err != nil {
			return err
		}

		organizationIds, err := cmd.Flags().GetStringSlice("organizationIds")
		if err != nil {
			return err
		}

		body := high5Client.PostV1ExecuteQueueJSONRequestBody{
			Blocking:        blocking,
			OrganizationIds: organizationIds,
		}

		if cmd.Flags().Changed("limit") {
			limit, err := cmd.Flags().GetInt("limit")
			if err != nil {
				return err
			}
			body.Limit = &limit
		}

		res, err := client.PostV1ExecuteQueueWithResponse(context.TODO(), &high5Client.PostV1ExecuteQueueParams{
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
