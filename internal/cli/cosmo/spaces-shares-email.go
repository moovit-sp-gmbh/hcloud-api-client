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
	"github.com/oapi-codegen/runtime/types"
	"github.com/spf13/cobra"
)

func init() {
	spacesSharesEmailCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesSharesEmailCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesSharesEmailCreateCommand.Flags().String("shareId", "", "")
	spacesSharesEmailCreateCommand.Flags().String("token", "", "authorization token")
	spacesSharesEmailCreateCommand.Flags().StringSlice("emails", nil, "Emails addresses to resend to")
	spacesSharesEmailAllCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesSharesEmailAllCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesSharesEmailAllCommand.Flags().String("shareId", "", "")
	spacesSharesEmailAllCommand.Flags().String("token", "", "authorization token")
	spacesSharesEmailPublicLinkCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesSharesEmailPublicLinkCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesSharesEmailPublicLinkCommand.Flags().String("shareId", "", "")
	spacesSharesEmailPublicLinkCommand.Flags().String("email", "", "Email addresses to send the share link to")
	spacesSharesEmailPublicLinkCommand.Flags().String("shareIdHMAC", "", "HMAC of the shareId provided in the public link")

	spacesSharesEmailCommand.AddCommand(spacesSharesEmailCreateCommand)
	spacesSharesEmailCommand.AddCommand(spacesSharesEmailAllCommand)
	spacesSharesEmailCommand.AddCommand(spacesSharesEmailPublicLinkCommand)
}

var spacesSharesEmailCommand = &cobra.Command{
	Use:   "email",
	Short: "Commands for email",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesSharesEmailCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Resend share emails",
	Long:  "Resends the share email to a set of the users that previously received emails relative to this share.",
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

		shareId, err := command.FlagOrPromptString(cmd, "shareId")
		if err != nil {
			return err
		}
		shareIdParsed, err := uuid.Parse(shareId)
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		emails, err := cmd.Flags().GetStringSlice("emails")
		if err != nil {
			return err
		}

		body := cosmoClient.PostV1OrgOrgNameSpacesSpaceNameSharesShareIdEmailJSONRequestBody{
			Emails: emails,
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameSharesShareIdEmailWithResponse(context.TODO(), orgName, spaceName, shareIdParsed, body, func(_ context.Context, req *http.Request) error {
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

var spacesSharesEmailAllCommand = &cobra.Command{
	Use:   "all",
	Short: "Resend all share emails",
	Long:  "Resends the share email to all the users that previously received emails relative to this share.",
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

		shareId, err := command.FlagOrPromptString(cmd, "shareId")
		if err != nil {
			return err
		}
		shareIdParsed, err := uuid.Parse(shareId)
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameSharesShareIdEmailAllWithResponse(context.TODO(), orgName, spaceName, shareIdParsed, func(_ context.Context, req *http.Request) error {
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

var spacesSharesEmailPublicLinkCommand = &cobra.Command{
	Use:   "public-link",
	Short: "Send the share email after public link got clicked",
	Long:  "Sends a share email to the users after the public link has been clicked.",
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

		shareId, err := command.FlagOrPromptString(cmd, "shareId")
		if err != nil {
			return err
		}
		shareIdParsed, err := uuid.Parse(shareId)
		if err != nil {
			return err
		}

		email, err := command.GetArg(cmd, command.StringFlagIfSet("email"), command.ContextIfValidString(latest.Email), command.PromptString("email"))
		if err != nil {
			return err
		}

		shareIdHMAC, err := command.FlagOrPromptString(cmd, "shareIdHMAC")
		if err != nil {
			return err
		}

		body := cosmoClient.PostV1OrgOrgNameSpacesSpaceNameSharesShareIdEmailPublicLinkJSONRequestBody{
			Email:       types.Email(email),
			ShareIdHMAC: shareIdHMAC,
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameSharesShareIdEmailPublicLinkWithResponse(context.TODO(), orgName, spaceName, shareIdParsed, body)
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
