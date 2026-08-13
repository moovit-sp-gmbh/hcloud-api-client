package mothership

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/mothershipClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	agentCommand.AddCommand(agentDeleteCommand)

	agentDeleteCommand.Flags().String("token", "", "Token obtained from the hello endpoint")
}

var agentCommand = &cobra.Command{
	Use: "agent",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
var agentDeleteCommand = &cobra.Command{
	Use: "delete",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerMothership(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		helloToken, err := command.GetArg(cmd, command.StringFlagIfSet("hello_token"), command.ContextIfValidString(latest.AgentHelloToken), command.PromptString("hello token"))
		if err != nil {
			return err
		}
		latest.AgentHelloToken = sql.NullString{Valid: true, String: helloToken}

		client, err := mothershipClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		res, err := client.DeleteV1AgentWithResponse(context.TODO(), func(ctx context.Context, req *http.Request) error {
			req.Header.Set("authorization", helloToken)
			return nil
		})

		var text command.Text
		switch res.StatusCode() {
		case 204:
			text = command.Text("Success")
			service.Push(latest)
		default:
			text, err = command.JSON(string(res.Body))
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}
