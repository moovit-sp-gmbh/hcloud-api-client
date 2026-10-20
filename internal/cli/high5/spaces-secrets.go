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
	spacesSecretsCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesSecretsCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesSecretsCreateCommand.Flags().String("token", "", "authorization token")
	spacesSecretsCreateCommand.Flags().String("key", "", "Key of the secret (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesSecretsCreateCommand.Flags().String("value", "", "Secret value (it should contain between 1 and 1024 characters)")
	spacesSecretsCreateCommand.Flags().String("description", "", "Secret description (maximum length is 80 symbols)")
	spacesSecretsCreateCommand.Flags().Bool("encrypted", false, "Flag to determine if the value should be stored encrypted. Only unencrypted values can be read again outside of the StreamExecutionPackage. Defaults to false")
	spacesSecretsDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesSecretsDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesSecretsDeleteCommand.Flags().String("token", "", "authorization token")
	spacesSecretsListCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesSecretsListCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesSecretsListCommand.Flags().String("token", "", "authorization token")
	spacesSecretsUpdateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesSecretsUpdateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesSecretsUpdateCommand.Flags().String("token", "", "authorization token")
	spacesSecretsUpdateCommand.Flags().String("value", "", "Secret value (it should contain between 1 and 1024 characters)")
	spacesSecretsUpdateCommand.Flags().String("description", "", "Secret description (maximum length is 80 symbols)")
	spacesSecretsUpdateCommand.Flags().Bool("encrypted", false, "Flag to determine if the value should be stored encrypted. Only unencrypted values can be read again outside of the StreamExecutionPackage. Defaults to false")
	spacesSecretsSearchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesSecretsSearchCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesSecretsSearchCommand.Flags().String("token", "", "authorization token")
	spacesSecretsSearchCommand.Flags().Float32("limit", 0, "Page size")
	spacesSecretsSearchCommand.Flags().Float32("page", 0, "Page index")
	spacesSecretsSearchCommand.Flags().String("body", "{}", "request body as JSON")

	spacesSecretsCommand.AddCommand(spacesSecretsCreateCommand)
	spacesSecretsCommand.AddCommand(spacesSecretsDeleteCommand)
	spacesSecretsCommand.AddCommand(spacesSecretsListCommand)
	spacesSecretsCommand.AddCommand(spacesSecretsUpdateCommand)
	spacesSecretsCommand.AddCommand(spacesSecretsSearchCommand)
}

var spacesSecretsCommand = &cobra.Command{
	Use:   "secrets",
	Short: "Commands for secrets",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesSecretsCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new secret for the Space",
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

		key, err := command.FlagOrPromptString(cmd, "key")
		if err != nil {
			return err
		}

		value, err := command.FlagOrPromptString(cmd, "value")
		if err != nil {
			return err
		}

		body := high5Client.PostV1OrgOrgNameSpacesSpaceNameSecretsJSONRequestBody{
			Key:   key,
			Value: value,
		}

		if cmd.Flags().Changed("description") {
			description, err := cmd.Flags().GetString("description")
			if err != nil {
				return err
			}
			body.Description = &description
		}

		if cmd.Flags().Changed("encrypted") {
			encrypted, err := cmd.Flags().GetBool("encrypted")
			if err != nil {
				return err
			}
			body.Encrypted = &encrypted
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameSecretsWithResponse(context.TODO(), orgName, spaceName, &high5Client.PostV1OrgOrgNameSpacesSpaceNameSecretsParams{
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

var spacesSecretsDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete a secret of the Space",
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

		res, err := client.DeleteV1OrgOrgNameSpacesSpaceNameSecretsSecretKeyWithResponse(context.TODO(), orgName, spaceName, &high5Client.DeleteV1OrgOrgNameSpacesSpaceNameSecretsSecretKeyParams{
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

var spacesSecretsListCommand = &cobra.Command{
	Use:   "list",
	Short: "Get a secret by its key",
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

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameSecretsSecretKeyWithResponse(context.TODO(), orgName, spaceName, &high5Client.GetV1OrgOrgNameSpacesSpaceNameSecretsSecretKeyParams{
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

var spacesSecretsUpdateCommand = &cobra.Command{
	Use:   "update",
	Short: "Update a secret for the Space",
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

		value, err := command.FlagOrPromptString(cmd, "value")
		if err != nil {
			return err
		}

		body := high5Client.PutV1OrgOrgNameSpacesSpaceNameSecretsSecretKeyJSONRequestBody{
			Value: value,
		}

		if cmd.Flags().Changed("description") {
			description, err := cmd.Flags().GetString("description")
			if err != nil {
				return err
			}
			body.Description = &description
		}

		if cmd.Flags().Changed("encrypted") {
			encrypted, err := cmd.Flags().GetBool("encrypted")
			if err != nil {
				return err
			}
			body.Encrypted = &encrypted
		}

		res, err := client.PutV1OrgOrgNameSpacesSpaceNameSecretsSecretKeyWithResponse(context.TODO(), orgName, spaceName, &high5Client.PutV1OrgOrgNameSpacesSpaceNameSecretsSecretKeyParams{
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

var spacesSecretsSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Retrieve data about the space's secrets",
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

		params := &high5Client.PostV1OrgOrgNameSpacesSpaceNameSecretsSearchParams{
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

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body high5Client.PostV1OrgOrgNameSpacesSpaceNameSecretsSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameSecretsSearchWithResponse(context.TODO(), orgName, spaceName, params, body)
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
