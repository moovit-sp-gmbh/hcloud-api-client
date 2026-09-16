package idp

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	organizationServiceAccountsGetCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	organizationServiceAccountsGetCommand.Flags().String("serviceAccountId", "", "ID of the service account")
	organizationServiceAccountsGetCommand.Flags().String("token", "", "authorization token")
	organizationServiceAccountsCommand.AddCommand(organizationServiceAccountsGetCommand)
}

var organizationServiceAccountsCommand = &cobra.Command{
	Use:   "service-accounts",
	Short: "Commands for service-accounts",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var organizationServiceAccountsGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Get a service account by its ID",
	Long:  "\nRequired organization role: ADMIN\nRequired scopes: IDP_ORGANIZATION_READ",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("organization name"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		serviceAccountId, err := command.FlagOrDifferentPromptString(cmd, "serviceAccountId", "service account ID")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.GetV1OrgOrgNameServiceAccountsServiceAccountIdWithResponse(context.TODO(), orgName, serviceAccountId, func(_ context.Context, req *http.Request) error {
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
