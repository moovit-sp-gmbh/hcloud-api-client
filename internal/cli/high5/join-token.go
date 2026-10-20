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
	joinTokenGetCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	joinTokenGetCommand.Flags().String("token", "", "authorization token")
	joinTokenGetCommand.Flags().Float32("limit", 0, "Page size")
	joinTokenGetCommand.Flags().Float32("page", 0, "Page index")
	joinTokenCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	joinTokenCreateCommand.Flags().String("token", "", "authorization token")
	joinTokenCreateCommand.Flags().String("name", "", "Name of the join token (only alphanumeric characters, underscores, hyphens and dots are allowed; it should contain between 1 and 64 characters)")
	joinTokenDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	joinTokenDeleteCommand.Flags().String("tokenId", "", "")
	joinTokenDeleteCommand.Flags().String("token", "", "authorization token")

	joinTokenCommand.AddCommand(joinTokenCreateAllCommand)
	joinTokenCommand.AddCommand(joinTokenGetCommand)
	joinTokenCommand.AddCommand(joinTokenCreateCommand)
	joinTokenCommand.AddCommand(joinTokenDeleteCommand)
	joinCommand.AddCommand(joinTokenCommand)
}

var joinCommand = &cobra.Command{
	Use:   "join",
	Short: "Commands for join",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var joinTokenCommand = &cobra.Command{
	Use:   "token",
	Short: "Commands for token",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var joinTokenCreateAllCommand = &cobra.Command{
	Use:   "create-all",
	Short: "Use a join token to create an agent user",
	Long:  "Create an agent user using a join token in the Authorization header",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewHigh5Client(cmd)
		if err != nil {
			return err
		}

		res, err := client.PostV1JoinTokenWithResponse(context.TODO())
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

		command.SetResult(cmd, text)

		return nil
	},
}

var joinTokenGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Retrieves the join tokens of an organization",
	Long:  "Paginated request to obtain the join tokens of an organization",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &high5Client.GetV1OrgOrgNameJoinTokenParams{
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

		res, err := client.GetV1OrgOrgNameJoinTokenWithResponse(context.TODO(), orgName, params)
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

var joinTokenCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Creates a join token for an organization",
	Long:  "Creates a join token for an organization with a given name that can be used to create agent users",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		name, err := command.FlagOrPromptString(cmd, "name")
		if err != nil {
			return err
		}

		body := high5Client.PostV1OrgOrgNameJoinTokenJSONRequestBody{
			Name: name,
		}

		res, err := client.PostV1OrgOrgNameJoinTokenWithResponse(context.TODO(), orgName, &high5Client.PostV1OrgOrgNameJoinTokenParams{
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

var joinTokenDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Revokes a join token of an organization",
	Long:  "Revoke a join token, making its use impossible",
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

		tokenId, err := command.FlagOrPromptString(cmd, "tokenId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1OrgOrgNameJoinTokenTokenIdWithResponse(context.TODO(), orgName, tokenId, &high5Client.DeleteV1OrgOrgNameJoinTokenTokenIdParams{
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
