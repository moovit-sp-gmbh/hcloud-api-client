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
	spacesExecuteStreamIdCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesExecuteStreamIdCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesExecuteStreamIdCreateCommand.Flags().String("streamId", "", "")
	spacesExecuteStreamIdCreateCommand.Flags().String("token", "", "authorization token")
	spacesExecuteStreamIdCreateCommand.Flags().Bool("design", false, "")
	spacesExecuteStreamIdCreateCommand.Flags().Bool("debug", false, "")
	spacesExecuteStreamIdCreateCommand.Flags().String("body", "{}", "request body as JSON")
	spacesExecuteStreamIdPackageCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesExecuteStreamIdPackageCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesExecuteStreamIdPackageCommand.Flags().String("streamId", "", "")
	spacesExecuteStreamIdPackageCommand.Flags().String("secret", "", "Secret that got generated on execution package creation at one of the POST execute endpoints (it should contain between 3 and 1024 characters)")
	spacesExecuteStreamIdPackageCommand.Flags().String("token", "", "authorization token")
	spacesExecuteStreamIdPackageCommand.Flags().String("hash", "", "SHA-256 hash of the client's cached design, used to determine if a fresh download is required")

	spacesExecuteStreamIdCommand.AddCommand(spacesExecuteStreamIdCreateCommand)
	spacesExecuteStreamIdCommand.AddCommand(spacesExecuteStreamIdPackageCommand)
	spacesExecuteStreamCommand.AddCommand(spacesExecuteStreamIdCommand)
}

var spacesExecuteStreamCommand = &cobra.Command{
	Use:   "stream",
	Short: "Commands for stream",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesExecuteStreamIdCommand = &cobra.Command{
	Use:   "id",
	Short: "Commands for id",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesExecuteStreamIdCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Execute a stream by it's ID",
	Long:  "Creates execution entries in the database and dispatches the stream to the specified target agent via NATS; supports an optional 'design' flag to execute the current design instead of the published snapshot and an optional 'debug' flag to enable step-by-step debug mode",
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

		streamId, err := command.FlagOrPromptString(cmd, "streamId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &high5Client.PostV1OrgOrgNameSpacesSpaceNameExecuteStreamIdStreamIdParams{
			Authorization: jwt,
		}

		if cmd.Flags().Changed("design") {
			design, err := cmd.Flags().GetBool("design")
			if err != nil {
				return err
			}
			params.Design = &design
		}

		if cmd.Flags().Changed("debug") {
			debug, err := cmd.Flags().GetBool("debug")
			if err != nil {
				return err
			}
			params.Debug = &debug
		}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body high5Client.PostV1OrgOrgNameSpacesSpaceNameExecuteStreamIdStreamIdJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameExecuteStreamIdStreamIdWithResponse(context.TODO(), orgName, spaceName, streamId, params, body)
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

var spacesExecuteStreamIdPackageCommand = &cobra.Command{
	Use:   "package",
	Short: "Request the execution package by streamId and a secret",
	Long:  "Returns the complete execution package needed by the agent to run a stream, including the stream build, payload, wave engine version and catalog data; the optional 'hash' query parameter allows the agent to skip re-downloading an unchanged design build",
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

		streamId, err := command.FlagOrPromptString(cmd, "streamId")
		if err != nil {
			return err
		}

		secret, err := command.FlagOrPromptString(cmd, "secret")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &high5Client.GetV1OrgOrgNameSpacesSpaceNameExecuteStreamIdStreamIdPackageSecretParams{
			Authorization: jwt,
		}

		if cmd.Flags().Changed("hash") {
			hash, err := cmd.Flags().GetString("hash")
			if err != nil {
				return err
			}
			params.Hash = &hash
		}

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameExecuteStreamIdStreamIdPackageSecretWithResponse(context.TODO(), orgName, spaceName, streamId, secret, params)
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
