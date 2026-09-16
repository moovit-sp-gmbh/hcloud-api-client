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
	organizationSettingsDomainsSamlCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	organizationSettingsDomainsSamlCommand.Flags().String("token", "", "authorization token")
	organizationSettingsDomainsOidcCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	organizationSettingsDomainsOidcCommand.Flags().String("token", "", "authorization token")
	organizationSettingsDomainsGetCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	organizationSettingsDomainsGetCommand.Flags().String("token", "", "authorization token")
	organizationSettingsDomainsCommand.AddCommand(organizationSettingsDomainsGetCommand)
	organizationSettingsDomainsCommand.AddCommand(organizationSettingsDomainsOidcCommand)
	organizationSettingsDomainsCommand.AddCommand(organizationSettingsDomainsSamlCommand)
}

var organizationSettingsDomainsCommand = &cobra.Command{
	Use:   "domains",
	Short: "Commands for domains",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var organizationSettingsDomainsGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Get all domains associated with an organization",
	Long:  "\nRequired role: ADMIN",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.GetV1OrgOrgNameSettingsDomainsWithResponse(context.TODO(), orgName, func(_ context.Context, req *http.Request) error {
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

var organizationSettingsDomainsOidcCommand = &cobra.Command{
	Use:   "oidc",
	Short: "Gets all providers for an organization",
	Long:  "\nRequired organization role: ADMIN",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.GetV1OrgOrgNameSettingsDomainsOidcWithResponse(context.TODO(), orgName, func(_ context.Context, req *http.Request) error {
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

var organizationSettingsDomainsSamlCommand = &cobra.Command{
	Use:   "saml",
	Short: "Gets all providers for an organization",
	Long:  "\nRequired organization role: ADMIN",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.GetV1OrgOrgNameSettingsDomainsSamlWithResponse(context.TODO(), orgName, func(_ context.Context, req *http.Request) error {
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
