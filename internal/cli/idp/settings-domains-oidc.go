package idp

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/idpClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	settingsDomainsOidcDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	settingsDomainsOidcDeleteCommand.Flags().String("domainName", "", "")
	settingsDomainsOidcDeleteCommand.Flags().String("token", "", "authorization token")
	settingsDomainsOidcPatchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	settingsDomainsOidcPatchCommand.Flags().String("domainName", "", "")
	settingsDomainsOidcPatchCommand.Flags().String("token", "", "authorization token")
	settingsDomainsOidcPatchCommand.Flags().String("clientId", "", "")
	settingsDomainsOidcPatchCommand.Flags().String("clientSecret", "", "")
	settingsDomainsOidcCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	settingsDomainsOidcCreateCommand.Flags().String("domainName", "", "")
	settingsDomainsOidcCreateCommand.Flags().String("token", "", "authorization token")
	settingsDomainsOidcCreateCommand.Flags().String("clientId", "", "")
	settingsDomainsOidcCreateCommand.Flags().String("clientSecret", "", "")

	settingsDomainsOidcCommand.AddCommand(settingsDomainsOidcDeleteCommand)
	settingsDomainsOidcCommand.AddCommand(settingsDomainsOidcPatchCommand)
	settingsDomainsOidcCommand.AddCommand(settingsDomainsOidcCreateCommand)
}

var settingsDomainsOidcCommand = &cobra.Command{
	Use:   "oidc",
	Short: "Commands for oidc",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var settingsDomainsOidcDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a OIDC IdP",
	Long:  "Required organization role: ADMIN",
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

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("orgName"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		domainName, err := command.FlagOrPromptString(cmd, "domainName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1OrgOrgNameSettingsDomainsDomainNameOidcWithResponse(context.TODO(), orgName, domainName, func(_ context.Context, req *http.Request) error {
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

var settingsDomainsOidcPatchCommand = &cobra.Command{
	Use:   "patch",
	Short: "Edit the current configuration of the OIDC provider linked with a domain",
	Long:  "Required organization role: ADMIN",
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

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("orgName"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		domainName, err := command.FlagOrPromptString(cmd, "domainName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		body := idpClient.PatchV1OrgOrgNameSettingsDomainsDomainNameOidcJSONRequestBody{}

		if cmd.Flags().Changed("clientId") {
			clientId, err := cmd.Flags().GetString("clientId")
			if err != nil {
				return err
			}
			body.ClientId = &clientId
		}

		if cmd.Flags().Changed("clientSecret") {
			clientSecret, err := cmd.Flags().GetString("clientSecret")
			if err != nil {
				return err
			}
			body.ClientSecret = &clientSecret
		}

		res, err := client.PatchV1OrgOrgNameSettingsDomainsDomainNameOidcWithResponse(context.TODO(), orgName, domainName, body, func(_ context.Context, req *http.Request) error {
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

var settingsDomainsOidcCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Creates a new provider for the specified domain in an organization",
	Long:  "Creates a new provider for the specified domain in an organization\nThere are multiple ways the process can fail besides authentication/authorization:\n- The domain is not verified\n- There is already an IdP associated with this organization and domain\n\nRequired organization role: ADMIN",
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

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("orgName"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		domainName, err := command.FlagOrPromptString(cmd, "domainName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		clientId, err := command.FlagOrPromptString(cmd, "clientId")
		if err != nil {
			return err
		}

		clientSecret, err := command.FlagOrPromptString(cmd, "clientSecret")
		if err != nil {
			return err
		}

		body := idpClient.PostV1OrgOrgNameSettingsDomainsDomainNameOidcJSONRequestBody{
			ClientId:     clientId,
			ClientSecret: clientSecret,
		}

		res, err := client.PostV1OrgOrgNameSettingsDomainsDomainNameOidcWithResponse(context.TODO(), orgName, domainName, body, func(_ context.Context, req *http.Request) error {
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
