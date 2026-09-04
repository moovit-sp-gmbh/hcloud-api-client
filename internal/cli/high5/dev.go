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
	spaceEditDevCommand.Flags().String("orgName", "", "Name of the organization")
	spaceEditDevCommand.Flags().String("spaceName", "", "Name of the space")
	spaceEditDevCommand.Flags().String("token", "", "")
	spaceEditDevCommand.Flags().Bool("enableDev", false, "Specifies to the front-end whether to display the versions under development of Wave Engines and Node Catalogs")
	spaceEditCommand.AddCommand(spaceEditDevCommand)
}

var spaceEditCommand = &cobra.Command{
	Use:   "edit",
	Short: "Commands for edit",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spaceEditDevCommand = &cobra.Command{
	Use:   "dev",
	Short: "Edit the enable dev flag on a space",
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

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("organization name"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		spaceName, err := command.GetArg(cmd, command.StringFlagIfSet("spaceName"), command.ContextIfValidString(latest.High5SpaceName), command.PromptString("space name"))
		if err != nil {
			return err
		}
		latest.High5SpaceName = sql.NullString{Valid: true, String: spaceName}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		enableDev, err := command.FlagOrPromptBool(cmd, "enableDev")
		if err != nil {
			return err
		}

		body := high5Client.PatchV1OrgOrgNameSpacesSpaceNameDevJSONRequestBody{
			EnableDev: enableDev,
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameDevWithResponse(context.TODO(), orgName, spaceName, &high5Client.PatchV1OrgOrgNameSpacesSpaceNameDevParams{
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
