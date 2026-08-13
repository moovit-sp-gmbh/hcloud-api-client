package mothership

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/oapi-codegen/runtime/types"
	"github.com/spf13/cobra"
)

func init() {
	disconnectCommand.Flags().String("orgName", "", "Name of the organization")
	disconnectCommand.Flags().String("email", "", "Email to disconnect")
	disconnectCommand.Flags().String("token", "", "token gotten from the hello phase")

	disconnectCommand.AddCommand(disconnectAllCommand)
	disconnectAllCommand.Flags().String("email", "", "Email to disconnect")
	disconnectAllCommand.Flags().String("token", "", "token gotten from the hello phase")
}

var disconnectCommand = &cobra.Command{
	Use:   "disconnect",
	Short: "Disconnect from an organization",
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

		helloToken, err := command.GetArg(cmd, command.StringFlagIfSet("hello_token"), command.ContextIfValidString(latest.AgentHelloToken), command.PromptString("hello token"))

		if err != nil {
			return err
		}

		res, err := client.DeleteV1OrgOrgNameConnectEmailWithResponse(context.TODO(), orgName, types.Email(email), func(ctx context.Context, req *http.Request) error {
			req.Header.Set("authorization", helloToken)
			return nil
		})
		if err != nil {
			return err
		}

		var text command.Text
		if res.StatusCode() < 200 || res.StatusCode() >= 300 {
			var data any
			err = json.Unmarshal(res.Body, &data)
			if err != nil {
				return err
			}
			text, err = command.JSON(data)
			if err != nil {
				return err
			}
		} else {
			text = command.Text("Success")
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var disconnectAllCommand = &cobra.Command{
	Use:   "all",
	Short: "Disconnect completely",
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

		email, err := command.GetArg(cmd, command.StringFlagIfSet("email"), command.PromptString("email"))
		if err != nil {
			return err
		}

		helloToken, err := command.GetArg(cmd, command.StringFlagIfSet("hello_token"), command.ContextIfValidString(latest.AgentHelloToken), command.PromptString("hello token"))

		if err != nil {
			return err
		}

		res, err := client.DeleteV1ConnectEmailWithResponse(context.TODO(), types.Email(email), func(ctx context.Context, req *http.Request) error {
			req.Header.Set("authorization", helloToken)
			return nil
		})
		if err != nil {
			return err
		}

		var text command.Text
		if res.StatusCode() < 200 || res.StatusCode() >= 300 {
			var data any
			err = json.Unmarshal(res.Body, &data)
			if err != nil {
				return err
			}
			text, err = command.JSON(data)
			if err != nil {
				return err
			}
		} else {
			text = command.Text("Success")
		}

		command.SetResult(cmd, text)

		return nil
	},
}
