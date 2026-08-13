package mothership

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/mothershipClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/oapi-codegen/runtime/types"
	"github.com/spf13/cobra"
)

func init() {
	connectCommand.Flags().String("orgName", "", "Name of the organization")
	connectCommand.Flags().String("email", "", "Email of the user")
	connectCommand.Flags().String("token", "", "JWT of the user")

	connectCommand.AddCommand(searchConnectedCommand)

	searchConnectedCommand.Flags().String("token", "", "JWT of the user")
	searchConnectedCommand.Flags().String("orgName", "", "Name of the organization")
	searchConnectedCommand.Flags().String("body", "{}", "Search filters and sorting order")
	searchConnectedCommand.Flags().Float32("limit", 10, "Page size")
	searchConnectedCommand.Flags().Float32("page", 0, "Page number")
}

var connectCommand = &cobra.Command{
	Use:   "connect",
	Short: "Connect to an organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewMothershipClient(cmd)
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

		email, err := command.GetArg(cmd, command.StringFlagIfSet("email"), command.PromptString("email"))
		if err != nil {
			return err
		}

		token, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("token"))
		if err != nil {
			return err
		}

		helloToken, err := command.GetArg(cmd, command.StringFlagIfSet("hello_token"), command.ContextIfValidString(latest.AgentHelloToken), command.PromptString("hello token"))
		if err != nil {
			return err
		}
		res, err := client.PostV1OrgOrgNameConnectWithResponse(context.TODO(), orgName, mothershipClient.PostV1OrgOrgNameConnectJSONRequestBody{
			Email:       types.Email(email),
			MemberToken: token,
		}, func(ctx context.Context, req *http.Request) error {
			req.Header.Set("authorization", helloToken)
			return nil
		})
		if err != nil {
			return err
		}

		var text command.Text
		if res.StatusCode() >= 200 && res.StatusCode() < 300 {
			text = command.Text("Success")
		} else {
			var data any
			err = json.Unmarshal(res.Body, &data)
			if err != nil {
				return err
			}
			text, err = command.JSON(data)
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var searchConnectedCommand = &cobra.Command{
	Use:   "search",
	Short: "Search the available execution targets for a given organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewMothershipClient(cmd)
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

		limit, err := cmd.Flags().GetFloat32("limit")
		if err != nil {
			return err
		}

		page, err := cmd.Flags().GetFloat32("page")
		if err != nil {
			return err
		}

		bodyStr, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}

		var body mothershipClient.PostV1OrgOrgNameConnectSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyStr), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameConnectSearchWithResponse(context.TODO(), orgName, &mothershipClient.PostV1OrgOrgNameConnectSearchParams{
			Limit: &limit,
			Page:  &page,
		}, body, func(ctx context.Context, req *http.Request) error {
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
