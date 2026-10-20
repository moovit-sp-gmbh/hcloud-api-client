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
	cronResolveCommand.Flags().String("token", "", "authorization token")
	cronResolveCommand.Flags().String("expression", "", "Quartz cron expression")
	cronResolveCommand.Flags().String("timezone", "", "Timezone of the Quartz expression in IANA format, e.g. Europe/Berlin")
	cronResolveCommand.Flags().Int("amount", 0, "Amount of returned execution-dates")

	cronCommand.AddCommand(cronResolveCommand)
}

var cronCommand = &cobra.Command{
	Use:   "cron",
	Short: "Commands for cron",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var cronResolveCommand = &cobra.Command{
	Use:   "resolve",
	Short: "Resolve a Quartz cron expression along with a timezone to get the next execution(s)",
	Long:  "Resolves a Quartz cron expression against a given timezone and returns the UTC timestamps (in milliseconds) of the next n upcoming executions without requiring an existing cronjob",
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

		expression, err := command.FlagOrPromptString(cmd, "expression")
		if err != nil {
			return err
		}

		timezone, err := command.FlagOrPromptString(cmd, "timezone")
		if err != nil {
			return err
		}

		body := high5Client.PostV1CronResolveJSONRequestBody{
			Expression: expression,
			Timezone:   timezone,
		}

		if cmd.Flags().Changed("amount") {
			amount, err := cmd.Flags().GetInt("amount")
			if err != nil {
				return err
			}
			body.Amount = &amount
		}

		res, err := client.PostV1CronResolveWithResponse(context.TODO(), &high5Client.PostV1CronResolveParams{
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
