package cosmo

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/cosmoClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	spacesHigh5LinkSpaceCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesHigh5LinkSpaceCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesHigh5LinkSpaceCommand.Flags().String("token", "", "authorization token")
	spacesHigh5LinkSpaceCommand.Flags().String("executionTarget", "", "Email address or pool name for the execution target")
	spacesHigh5LinkSpaceCommand.Flags().String("high5_spaceName", "", "Name for the space")

	spacesHigh5LinkCommand.AddCommand(spacesHigh5LinkSpaceCommand)
	spacesHigh5Command.AddCommand(spacesHigh5LinkCommand)
	spacesHigh5Command.AddCommand(spacesHigh5UnlinkCommand)
}

var spacesHigh5Command = &cobra.Command{
	Use:   "high5",
	Short: "Commands for high5",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesHigh5LinkCommand = &cobra.Command{
	Use:   "link",
	Short: "Commands for link",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesHigh5LinkSpaceCommand = &cobra.Command{
	Use:   "space",
	Short: "Link a high5 space to a cosmo space",
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

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		executionTarget, err := command.FlagOrPromptString(cmd, "executionTarget")
		if err != nil {
			return err
		}

		high5Spacename, err := command.GetArg(cmd, command.StringFlagIfSet("high5_spaceName"), command.ContextIfValidString(latest.High5SpaceName), command.PromptString("high5_spaceName"))
		if err != nil {
			return err
		}

		body := cosmoClient.PatchV1OrgOrgNameSpacesSpaceNameHigh5LinkSpaceJSONRequestBody{
			ExecutionTarget: executionTarget,
			High5SpaceName:  high5Spacename,
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameHigh5LinkSpaceWithResponse(context.TODO(), orgName, spaceName, body, func(_ context.Context, req *http.Request) error {
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
