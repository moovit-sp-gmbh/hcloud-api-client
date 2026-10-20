package cosmo

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	spacesHigh5UnlinkSpaceCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesHigh5UnlinkSpaceCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesHigh5UnlinkSpaceCommand.Flags().String("token", "", "authorization token")

	spacesHigh5UnlinkCommand.AddCommand(spacesHigh5UnlinkSpaceCommand)
}

var spacesHigh5UnlinkCommand = &cobra.Command{
	Use:   "unlink",
	Short: "Commands for unlink",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesHigh5UnlinkSpaceCommand = &cobra.Command{
	Use:   "space",
	Short: "Unlink a high5 space from a cosmo space",
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

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameHigh5UnlinkSpaceWithResponse(context.TODO(), orgName, spaceName, func(_ context.Context, req *http.Request) error {
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
