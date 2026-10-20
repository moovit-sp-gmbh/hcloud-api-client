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
	spacesJobsCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesJobsCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesJobsCreateCommand.Flags().String("token", "", "authorization token")
	spacesJobsCreateCommand.Flags().Bool("nextExecution", false, "Alongside the cronjob also retrieve the timestamps for three next executions (for the specified timezone, in milliseconds)")
	spacesJobsCreateCommand.Flags().String("body", "{}", "request body as JSON")
	spacesJobsSearchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesJobsSearchCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesJobsSearchCommand.Flags().String("token", "", "authorization token")
	spacesJobsSearchCommand.Flags().Float32("limit", 0, "Page size")
	spacesJobsSearchCommand.Flags().Float32("page", 0, "Page index")
	spacesJobsSearchCommand.Flags().Bool("nextExecution", false, "Alongside the cronjob also retrieve the timestamps for three next executions (for the specified timezone, in milliseconds)")
	spacesJobsSearchCommand.Flags().String("body", "{}", "request body as JSON")
	spacesJobsDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesJobsDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesJobsDeleteCommand.Flags().String("jobId", "", "")
	spacesJobsDeleteCommand.Flags().String("token", "", "authorization token")
	spacesJobsGetCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesJobsGetCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesJobsGetCommand.Flags().String("jobId", "", "")
	spacesJobsGetCommand.Flags().String("token", "", "authorization token")
	spacesJobsGetCommand.Flags().Bool("nextExecution", false, "Alongside the cronjob also retrieve the timestamps for three next executions (for the specified timezone, in milliseconds)")
	spacesJobsUpdateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesJobsUpdateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesJobsUpdateCommand.Flags().String("jobId", "", "")
	spacesJobsUpdateCommand.Flags().String("token", "", "authorization token")
	spacesJobsUpdateCommand.Flags().Bool("nextExecution", false, "Alongside the cronjob also retrieve the timestamps for three next executions (for the specified timezone, in milliseconds)")
	spacesJobsUpdateCommand.Flags().String("body", "{}", "request body as JSON")
	spacesJobsEnabledCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesJobsEnabledCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesJobsEnabledCommand.Flags().String("jobId", "", "")
	spacesJobsEnabledCommand.Flags().String("token", "", "authorization token")
	spacesJobsEnabledCommand.Flags().Bool("nextExecution", false, "Alongside the cronjob also retrieve the timestamps for three next executions (for the specified timezone, in milliseconds)")
	spacesJobsEnabledCommand.Flags().Bool("enabled", false, "Flag indicating whether the job is active and scheduled to run")
	spacesJobsExpressionCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesJobsExpressionCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesJobsExpressionCommand.Flags().String("jobId", "", "")
	spacesJobsExpressionCommand.Flags().String("token", "", "authorization token")
	spacesJobsExpressionCommand.Flags().Bool("nextExecution", false, "Alongside the cronjob also retrieve the timestamps for three next executions (for the specified timezone, in milliseconds)")
	spacesJobsExpressionCommand.Flags().String("body", "{}", "request body as JSON")
	spacesJobsNextCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesJobsNextCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesJobsNextCommand.Flags().String("jobId", "", "")
	spacesJobsNextCommand.Flags().String("token", "", "authorization token")
	spacesJobsNextCommand.Flags().Float32("amount", 0, "Amount of execution-timestamps that should be returned")

	spacesJobsCommand.AddCommand(spacesJobsCreateCommand)
	spacesJobsCommand.AddCommand(spacesJobsSearchCommand)
	spacesJobsCommand.AddCommand(spacesJobsDeleteCommand)
	spacesJobsCommand.AddCommand(spacesJobsGetCommand)
	spacesJobsCommand.AddCommand(spacesJobsUpdateCommand)
	spacesJobsCommand.AddCommand(spacesJobsEnabledCommand)
	spacesJobsCommand.AddCommand(spacesJobsExpressionCommand)
	spacesJobsCommand.AddCommand(spacesJobsLogsCommand)
	spacesJobsCommand.AddCommand(spacesJobsNextCommand)
}

var spacesJobsCommand = &cobra.Command{
	Use:   "jobs",
	Short: "Commands for jobs",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesJobsCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Creates a new cronjob",
	Long:  "Creates a new cronjob in the specified space using either a Quartz cron expression or a future Unix timestamp, which is automatically converted to a one-time Quartz expression",
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

		params := &high5Client.PostV1OrgOrgNameSpacesSpaceNameJobsParams{
			Authorization: jwt,
		}

		if cmd.Flags().Changed("nextExecution") {
			nextExecution, err := cmd.Flags().GetBool("nextExecution")
			if err != nil {
				return err
			}
			params.NextExecution = &nextExecution
		}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body high5Client.PostV1OrgOrgNameSpacesSpaceNameJobsJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameJobsWithResponse(context.TODO(), orgName, spaceName, params, body)
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

var spacesJobsSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Retrieves all cronjobs that match the provided search filter. Will return all cronjobs of the specified space if no search filter is provided",
	Long:  "Returns a paginated list of cronjobs within the specified space that match the provided search filter, returning all cronjobs of the space if no filter is given",
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

		params := &high5Client.PostV1OrgOrgNameSpacesSpaceNameJobsSearchParams{
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

		if cmd.Flags().Changed("nextExecution") {
			nextExecution, err := cmd.Flags().GetBool("nextExecution")
			if err != nil {
				return err
			}
			params.NextExecution = &nextExecution
		}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body high5Client.PostV1OrgOrgNameSpacesSpaceNameJobsSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameJobsSearchWithResponse(context.TODO(), orgName, spaceName, params, body)
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

var spacesJobsDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a cronjob",
	Long:  "Permanently deletes the specified cronjob along with all of its associated execution logs",
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

		jobId, err := command.FlagOrPromptString(cmd, "jobId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1OrgOrgNameSpacesSpaceNameJobsJobIdWithResponse(context.TODO(), orgName, spaceName, jobId, &high5Client.DeleteV1OrgOrgNameSpacesSpaceNameJobsJobIdParams{
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

var spacesJobsGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Retrieve a cronjob by its ID",
	Long:  "Returns the cronjob with the given ID from the specified space, optionally including the timestamps of the next three executions",
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

		jobId, err := command.FlagOrPromptString(cmd, "jobId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &high5Client.GetV1OrgOrgNameSpacesSpaceNameJobsJobIdParams{
			Authorization: jwt,
		}

		if cmd.Flags().Changed("nextExecution") {
			nextExecution, err := cmd.Flags().GetBool("nextExecution")
			if err != nil {
				return err
			}
			params.NextExecution = &nextExecution
		}

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameJobsJobIdWithResponse(context.TODO(), orgName, spaceName, jobId, params)
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

var spacesJobsUpdateCommand = &cobra.Command{
	Use:   "update",
	Short: "Updates a cronjob",
	Long:  "Fully replaces all fields of the specified cronjob with the provided data, accepting either a Quartz cron expression or a future Unix timestamp",
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

		jobId, err := command.FlagOrPromptString(cmd, "jobId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &high5Client.PutV1OrgOrgNameSpacesSpaceNameJobsJobIdParams{
			Authorization: jwt,
		}

		if cmd.Flags().Changed("nextExecution") {
			nextExecution, err := cmd.Flags().GetBool("nextExecution")
			if err != nil {
				return err
			}
			params.NextExecution = &nextExecution
		}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body high5Client.PutV1OrgOrgNameSpacesSpaceNameJobsJobIdJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PutV1OrgOrgNameSpacesSpaceNameJobsJobIdWithResponse(context.TODO(), orgName, spaceName, jobId, params, body)
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

var spacesJobsEnabledCommand = &cobra.Command{
	Use:   "enabled",
	Short: "Sets the cronjob's enabled flag",
	Long:  "Enables or disables the specified cronjob by updating its enabled flag",
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

		jobId, err := command.FlagOrPromptString(cmd, "jobId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &high5Client.PatchV1OrgOrgNameSpacesSpaceNameJobsJobIdEnabledParams{
			Authorization: jwt,
		}

		if cmd.Flags().Changed("nextExecution") {
			nextExecution, err := cmd.Flags().GetBool("nextExecution")
			if err != nil {
				return err
			}
			params.NextExecution = &nextExecution
		}

		enabled, err := command.FlagOrPromptBool(cmd, "enabled")
		if err != nil {
			return err
		}

		body := high5Client.PatchV1OrgOrgNameSpacesSpaceNameJobsJobIdEnabledJSONRequestBody{
			Enabled: enabled,
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameJobsJobIdEnabledWithResponse(context.TODO(), orgName, spaceName, jobId, params, body)
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

var spacesJobsExpressionCommand = &cobra.Command{
	Use:   "expression",
	Short: "Updates the quartz cron expression of the cronjob",
	Long:  "Updates the Quartz cron expression of the specified cronjob without affecting any other fields",
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

		jobId, err := command.FlagOrPromptString(cmd, "jobId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &high5Client.PatchV1OrgOrgNameSpacesSpaceNameJobsJobIdExpressionParams{
			Authorization: jwt,
		}

		if cmd.Flags().Changed("nextExecution") {
			nextExecution, err := cmd.Flags().GetBool("nextExecution")
			if err != nil {
				return err
			}
			params.NextExecution = &nextExecution
		}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body high5Client.PatchV1OrgOrgNameSpacesSpaceNameJobsJobIdExpressionJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameJobsJobIdExpressionWithResponse(context.TODO(), orgName, spaceName, jobId, params, body)
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

var spacesJobsNextCommand = &cobra.Command{
	Use:   "next",
	Short: "Retrieves the timestamps (in milliseconds, UTC+0) of the next n cronjob executions",
	Long:  "Returns up to 10 UTC timestamps (in milliseconds) for the upcoming executions of the specified cronjob, defaulting to 3 if no amount is provided",
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

		jobId, err := command.FlagOrPromptString(cmd, "jobId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &high5Client.GetV1OrgOrgNameSpacesSpaceNameJobsJobIdNextParams{
			Authorization: jwt,
		}

		if cmd.Flags().Changed("amount") {
			amount, err := cmd.Flags().GetFloat32("amount")
			if err != nil {
				return err
			}
			params.Amount = &amount
		}

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameJobsJobIdNextWithResponse(context.TODO(), orgName, spaceName, jobId, params)
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
