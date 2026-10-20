package cosmo

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/cosmoClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	spacesNamespacesRefMetadataCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesRefMetadataCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNamespacesRefMetadataCommand.Flags().String("namespaceName", "", "Name of the namespace (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesRefMetadataCommand.Flags().String("refId", "", "")
	spacesNamespacesRefMetadataCommand.Flags().String("token", "", "authorization token")
	spacesNamespacesRefMetadataCommand.Flags().String("body", "{}", "request body as JSON")
	spacesNamespacesRefStatusCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesRefStatusCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNamespacesRefStatusCommand.Flags().String("namespaceName", "", "Name of the namespace (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesRefStatusCommand.Flags().String("refId", "", "")
	spacesNamespacesRefStatusCommand.Flags().String("token", "", "authorization token")
	spacesNamespacesRefStatusCommand.Flags().String("status", "", "")

	spacesNamespacesRefCommand.AddCommand(spacesNamespacesRefMetadataCommand)
	spacesNamespacesRefCommand.AddCommand(spacesNamespacesRefStatusCommand)
}

var spacesNamespacesRefCommand = &cobra.Command{
	Use:   "ref",
	Short: "Commands for ref",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesNamespacesRefMetadataCommand = &cobra.Command{
	Use:   "metadata",
	Short: "Update the metadata of a reference (e.g. Media_Asset)",
	Long:  "Replace the current metadata of a reference with a new metadata set. To add a single field one must send the current metadata set with the extra field, not just the extra field. For deleting, one must send the current metadata set without the field to be deleted.",
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

		namespaceName, err := command.FlagOrPromptString(cmd, "namespaceName")
		if err != nil {
			return err
		}

		refId, err := command.FlagOrPromptString(cmd, "refId")
		if err != nil {
			return err
		}
		refIdParsed, err := uuid.Parse(refId)
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
		var body cosmoClient.PutV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameRefRefIdMetadataJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PutV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameRefRefIdMetadataWithResponse(context.TODO(), orgName, spaceName, namespaceName, refIdParsed, body, func(_ context.Context, req *http.Request) error {
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

var spacesNamespacesRefStatusCommand = &cobra.Command{
	Use:   "status",
	Short: "Change the status of a reference (e.g. Media_Asset)",
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

		namespaceName, err := command.FlagOrPromptString(cmd, "namespaceName")
		if err != nil {
			return err
		}

		refId, err := command.FlagOrPromptString(cmd, "refId")
		if err != nil {
			return err
		}
		refIdParsed, err := uuid.Parse(refId)
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		status, err := command.FlagOrPromptString(cmd, "status")
		if err != nil {
			return err
		}

		body := cosmoClient.PatchV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameRefRefIdStatusJSONRequestBody{
			Status: cosmoClient.ReferencePatchDtoStatus(status),
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameRefRefIdStatusWithResponse(context.TODO(), orgName, spaceName, namespaceName, refIdParsed, body, func(_ context.Context, req *http.Request) error {
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
