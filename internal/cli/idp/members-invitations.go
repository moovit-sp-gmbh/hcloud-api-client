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
	"github.com/oapi-codegen/runtime/types"
	"github.com/spf13/cobra"
)

func init() {
	membersInvitationsCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	membersInvitationsCreateCommand.Flags().String("token", "", "authorization token")
	membersInvitationsCreateCommand.Flags().String("email", "", "Email of the user")
	membersInvitationsCreateCommand.Flags().String("role", "", "Role the user shall have in the Organization")
	membersInvitationsCreateCommand.Flags().Bool("allowNonRegisteredUsers", false, "If non-registered users are allowed and the email is not bound to an HCloud user, then a registration email will be sent to this email. Upon registration, the user will be added automatically to the organization.")
	membersInvitationsCreateCommand.Flags().String("targetUrl", "", "URL to create the redirect link in the email. It will default to https://app.helmut.cloud.")
	membersInvitationsSearchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	membersInvitationsSearchCommand.Flags().String("token", "", "authorization token")
	membersInvitationsSearchCommand.Flags().Float32("limit", 0, "Page size")
	membersInvitationsSearchCommand.Flags().Float32("page", 0, "Page index")
	membersInvitationsSearchCommand.Flags().String("body", "{}", "request body as JSON")
	membersInvitationsDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	membersInvitationsDeleteCommand.Flags().String("invitationId", "", "")
	membersInvitationsDeleteCommand.Flags().String("token", "", "authorization token")
	membersInvitationsApprovalCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	membersInvitationsApprovalCommand.Flags().String("invitationId", "", "")
	membersInvitationsApprovalCommand.Flags().String("token", "", "authorization token")
	membersInvitationsApprovalCommand.Flags().Bool("accept", false, "True to accept. False to reject")
	membersInvitationsResendCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	membersInvitationsResendCommand.Flags().String("invitationId", "", "")
	membersInvitationsResendCommand.Flags().String("token", "", "authorization token")

	membersInvitationsCommand.AddCommand(membersInvitationsCreateCommand)
	membersInvitationsCommand.AddCommand(membersInvitationsSearchCommand)
	membersInvitationsCommand.AddCommand(membersInvitationsDeleteCommand)
	membersInvitationsCommand.AddCommand(membersInvitationsApprovalCommand)
	membersInvitationsCommand.AddCommand(membersInvitationsResendCommand)
}

var membersInvitationsCommand = &cobra.Command{
	Use:   "invitations",
	Short: "Commands for invitations",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var membersInvitationsCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Invite a member to an organization",
	Long:  "Required role: ADMIN\nRequired scopes: IDP_ORGANIZATION_WRITE",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		email, err := command.GetArg(cmd, command.StringFlagIfSet("email"), command.ContextIfValidString(latest.Email), command.PromptString("email"))
		if err != nil {
			return err
		}

		role, err := command.FlagOrPromptString(cmd, "role")
		if err != nil {
			return err
		}

		body := idpClient.PostV1OrgOrgNameMembersInvitationsJSONRequestBody{
			Email: types.Email(email),
			Role:  idpClient.OrganizationRole(role),
		}

		if cmd.Flags().Changed("allowNonRegisteredUsers") {
			allowNonRegisteredUsers, err := cmd.Flags().GetBool("allowNonRegisteredUsers")
			if err != nil {
				return err
			}
			body.AllowNonRegisteredUsers = &allowNonRegisteredUsers
		}

		if cmd.Flags().Changed("targetUrl") {
			targetUrl, err := cmd.Flags().GetString("targetUrl")
			if err != nil {
				return err
			}
			body.TargetUrl = &targetUrl
		}

		res, err := client.PostV1OrgOrgNameMembersInvitationsWithResponse(context.TODO(), orgName, body, func(_ context.Context, req *http.Request) error {
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

var membersInvitationsSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Retrieves all invitations to the specified organization that match the provided search filter. Will return all organization invitations if no filter provided",
	Long:  "Required role: MANAGER\nRequired scopes: IDP_ORGANIZATION_READ",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &idpClient.PostV1OrgOrgNameMembersInvitationsSearchParams{}

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
		var body idpClient.PostV1OrgOrgNameMembersInvitationsSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameMembersInvitationsSearchWithResponse(context.TODO(), orgName, params, body, func(_ context.Context, req *http.Request) error {
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

var membersInvitationsDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Cancel an invitation",
	Long:  "Required role: ADMIN\nRequired scopes: IDP_ORGANIZATION_WRITE",
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

		invitationId, err := command.FlagOrPromptString(cmd, "invitationId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1OrgOrgNameMembersInvitationsInvitationIdWithResponse(context.TODO(), orgName, invitationId, func(_ context.Context, req *http.Request) error {
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

var membersInvitationsApprovalCommand = &cobra.Command{
	Use:   "approval",
	Short: "Accept/Reject an invitation to an organization",
	Long:  "Required scopes: IDP_ORGANIZATION_WRITE",
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

		invitationId, err := command.FlagOrPromptString(cmd, "invitationId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		accept, err := command.FlagOrPromptBool(cmd, "accept")
		if err != nil {
			return err
		}

		body := idpClient.PatchV1OrgOrgNameMembersInvitationsInvitationIdApprovalJSONRequestBody{
			Accept: accept,
		}

		res, err := client.PatchV1OrgOrgNameMembersInvitationsInvitationIdApprovalWithResponse(context.TODO(), orgName, invitationId, body, func(_ context.Context, req *http.Request) error {
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

var membersInvitationsResendCommand = &cobra.Command{
	Use:   "resend",
	Short: "Resend the invitation email",
	Long:  "Required role: ADMIN\nRequired scopes: IDP_ORGANIZATION_WRITE",
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

		invitationId, err := command.FlagOrPromptString(cmd, "invitationId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.PostV1OrgOrgNameMembersInvitationsInvitationIdResendWithResponse(context.TODO(), orgName, invitationId, func(_ context.Context, req *http.Request) error {
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
