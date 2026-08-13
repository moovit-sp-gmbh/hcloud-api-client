package mothership

import (
	"context"
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
