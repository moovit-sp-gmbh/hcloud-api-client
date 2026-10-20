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
	spacesRolesUserUpdateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesRolesUserUpdateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesRolesUserUpdateCommand.Flags().String("roleId", "", "Role ID")
	spacesRolesUserUpdateCommand.Flags().String("token", "", "authorization token")
	spacesRolesUserUpdateCommand.Flags().String("userId", "", "ID of the user")
	spacesRolesUserDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesRolesUserDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesRolesUserDeleteCommand.Flags().String("roleId", "", "Role ID")
	spacesRolesUserDeleteCommand.Flags().String("userId", "", "User ID")
	spacesRolesUserDeleteCommand.Flags().String("token", "", "authorization token")

	spacesRolesUserCommand.AddCommand(spacesRolesUserUpdateCommand)
	spacesRolesUserCommand.AddCommand(spacesRolesUserDeleteCommand)
}

var spacesRolesUserCommand = &cobra.Command{
	Use:   "user",
	Short: "Commands for user",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesRolesUserUpdateCommand = &cobra.Command{
	Use:   "update",
	Short: "Assign a user to a custom role",
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

		roleId, err := command.FlagOrPromptString(cmd, "roleId")
		if err != nil {
			return err
		}
		roleIdParsed, err := uuid.Parse(roleId)
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		userId, err := command.FlagOrPromptString(cmd, "userId")
		if err != nil {
			return err
		}

		body := cosmoClient.PutV1OrgOrgNameSpacesSpaceNameRolesRoleIdUserJSONRequestBody{
			UserId: userId,
		}

		res, err := client.PutV1OrgOrgNameSpacesSpaceNameRolesRoleIdUserWithResponse(context.TODO(), orgName, spaceName, roleIdParsed, body, func(_ context.Context, req *http.Request) error {
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

var spacesRolesUserDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Unassign a user from a custom role",
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

		roleId, err := command.FlagOrPromptString(cmd, "roleId")
		if err != nil {
			return err
		}
		roleIdParsed, err := uuid.Parse(roleId)
		if err != nil {
			return err
		}

		userId, err := command.FlagOrPromptString(cmd, "userId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1OrgOrgNameSpacesSpaceNameRolesRoleIdUserUserIdWithResponse(context.TODO(), orgName, spaceName, roleIdParsed, userId, func(_ context.Context, req *http.Request) error {
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
