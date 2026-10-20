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
	membersCommand.AddCommand(membersDeleteAllCommand)
	membersCommand.AddCommand(membersInvitationsCommand)
	membersCommand.AddCommand(membersInvitationsListCommand)
	membersCommand.AddCommand(membersSearchCommand)
	membersCommand.AddCommand(membersDeleteCommand)
	membersCommand.AddCommand(membersExecutionTargetCommand)
	membersCommand.AddCommand(membersRoleCommand)
	organizationMembersGetCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	organizationMembersGetCommand.Flags().String("userId", "", "")
	organizationMembersGetCommand.Flags().String("token", "", "authorization token")
	membersDeleteAllCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	membersDeleteAllCommand.Flags().String("token", "", "authorization token")
	membersInvitationsListCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	membersInvitationsListCommand.Flags().String("token", "", "authorization token")
	membersInvitationsListCommand.Flags().String("body", "{}", "request body as JSON")
	membersSearchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	membersSearchCommand.Flags().String("token", "", "authorization token")
	membersSearchCommand.Flags().Float32("limit", 0, "Page size")
	membersSearchCommand.Flags().Float32("page", 0, "Page index")
	membersSearchCommand.Flags().String("body", "{}", "request body as JSON")
	membersDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	membersDeleteCommand.Flags().String("userId", "", "")
	membersDeleteCommand.Flags().String("token", "", "authorization token")
	membersExecutionTargetCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	membersExecutionTargetCommand.Flags().String("userId", "", "")
	membersExecutionTargetCommand.Flags().String("token", "", "authorization token")
	membersExecutionTargetCommand.Flags().Bool("executionTarget", false, "Boolean flag that determines if this user can be an execution target for this organization.")
	membersRoleCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	membersRoleCommand.Flags().String("userId", "", "")
	membersRoleCommand.Flags().String("token", "", "authorization token")
	membersRoleCommand.Flags().String("role", "", "")
	organizationMembersCommand.AddCommand(organizationMembersGetCommand)
}

var organizationMembersCommand = &cobra.Command{
	Use:   "members",
	Short: "Commands for members",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var organizationMembersGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Get an organization member by the user ID",
	Long:  "\nRequired role: MEMBER\nRequired scopes: IDP_ORGANIZATION_READ",
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

		userId, err := command.FlagOrDifferentPromptString(cmd, "userId", "user ID")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.GetV1OrgOrgNameMembersUserIdWithResponse(context.TODO(), orgName, userId, func(_ context.Context, req *http.Request) error {
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

var membersCommand = &cobra.Command{
	Use:   "members",
	Short: "Commands for members",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var membersDeleteAllCommand = &cobra.Command{
	Use:   "delete-all",
	Short: "Leave shared organization by the request of member",
	Long:  "Required role: MEMBER\nRequired scopes: IDP_ORGANIZATION_READ",
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

		res, err := client.DeleteV1OrgOrgNameMembersWithResponse(context.TODO(), orgName, func(_ context.Context, req *http.Request) error {
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

var membersInvitationsListCommand = &cobra.Command{
	Use:   "invitations-list",
	Short: "Invite multiple members to an organization",
	Long:  "Sends invitations to a list of users to join the specified organization\nRequired role: ADMIN\nRequired scopes: IDP_ORGANIZATION_WRITE",
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

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body idpClient.PostV1OrgOrgNameMembersInvitationsListJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameMembersInvitationsListWithResponse(context.TODO(), orgName, body, func(_ context.Context, req *http.Request) error {
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

var membersSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Search organization members with search filters",
	Long:  "Required role: MEMBER\nRequired scopes: IDP_ORGANIZATION_READ",
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

		params := &idpClient.PostV1OrgOrgNameMembersSearchParams{}

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
		var body idpClient.PostV1OrgOrgNameMembersSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameMembersSearchWithResponse(context.TODO(), orgName, params, body, func(_ context.Context, req *http.Request) error {
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

var membersDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Remove a user from an organization by the organization user ID",
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

		userId, err := command.FlagOrPromptString(cmd, "userId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1OrgOrgNameMembersUserIdWithResponse(context.TODO(), orgName, userId, func(_ context.Context, req *http.Request) error {
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

var membersExecutionTargetCommand = &cobra.Command{
	Use:   "execution-target",
	Short: "Set whether or not this user should be considered an execution target",
	Long:  "A member can set its execution target to true or false.\nA manager of an organization can set the execution target of any other member to false, but not to true.\nOnly the member itself (when userId is the same as the ID of the request user) can set to true\nRequired role: MEMBER\nRequired scopes: IDP_ORGANIZATION_WRITE",
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

		userId, err := command.FlagOrPromptString(cmd, "userId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		executionTarget, err := command.FlagOrPromptBool(cmd, "executionTarget")
		if err != nil {
			return err
		}

		body := idpClient.PatchV1OrgOrgNameMembersUserIdExecutionTargetJSONRequestBody{
			ExecutionTarget: executionTarget,
		}

		res, err := client.PatchV1OrgOrgNameMembersUserIdExecutionTargetWithResponse(context.TODO(), orgName, userId, body, func(_ context.Context, req *http.Request) error {
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

var membersRoleCommand = &cobra.Command{
	Use:   "role",
	Short: "Patch an organization member's role in the Organization",
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

		userId, err := command.FlagOrPromptString(cmd, "userId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		role, err := command.FlagOrPromptString(cmd, "role")
		if err != nil {
			return err
		}

		body := idpClient.PatchV1OrgOrgNameMembersUserIdRoleJSONRequestBody{
			Role: idpClient.OrgMemberPatchDtoRole(role),
		}

		res, err := client.PatchV1OrgOrgNameMembersUserIdRoleWithResponse(context.TODO(), orgName, userId, body, func(_ context.Context, req *http.Request) error {
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
