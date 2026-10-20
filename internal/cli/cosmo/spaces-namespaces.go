package cosmo

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/cosmoClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	spacesNamespacesCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNamespacesCreateCommand.Flags().String("token", "", "authorization token")
	spacesNamespacesCreateCommand.Flags().String("name", "", "Namespace name")
	spacesNamespacesCreateCommand.Flags().String("defaultStatus", "", "Namespace name")

	spacesNamespacesCommand.AddCommand(spacesNamespacesCreateCommand)
	spacesNamespacesCommand.AddCommand(spacesNamespacesAssetsCommand)
	spacesNamespacesCommand.AddCommand(spacesNamespacesCommentsCommand)
	spacesNamespacesCommand.AddCommand(spacesNamespacesRefCommand)
	spacesNamespacesCommand.AddCommand(spacesNamespacesTagsCommand)
}

var spacesNamespacesCommand = &cobra.Command{
	Use:   "namespaces",
	Short: "Commands for namespaces",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesNamespacesCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create Namespace",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewCosmoClient(cmd)
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

		name, err := command.FlagOrPromptString(cmd, "name")
		if err != nil {
			return err
		}

		body := cosmoClient.PostV1OrgOrgNameSpacesSpaceNameNamespacesJSONRequestBody{
			Name: name,
		}

		if cmd.Flags().Changed("defaultStatus") {
			defaultStatus, err := cmd.Flags().GetString("defaultStatus")
			if err != nil {
				return err
			}
			body.DefaultStatus = &defaultStatus
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameNamespacesWithResponse(context.TODO(), orgName, spaceName, body, func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", jwt)
			return nil
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
