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
	spaceGetCommand.Flags().String("orgName", "", "")
	spaceGetCommand.Flags().String("spaceName", "", "")
	spaceGetCommand.Flags().String("token", "", "")
	spaceCreateCommand.Flags().String("orgName", "", "Name of the organization")
	spaceCreateCommand.Flags().String("token", "", "Authorization token to use in the API call")
	spaceCreateCommand.Flags().Bool("noAvatar", false, "Skip avatar for space")
	spaceCreateCommand.Flags().String("name", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spaceCreateCommand.Flags().Bool("enableDev", false, "Specifies to the front-end whether to display the versions under development of Wave Engines and Node Catalogs")
	spaceCommand.AddCommand(spaceCreateCommand)
	spaceCommand.AddCommand(spaceEditCommand)
	spaceCommand.AddCommand(spaceGetCommand)
}

var spaceCommand = &cobra.Command{
	Use:   "space",
	Short: "Commands for space",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spaceCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new space in High5",
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

		params := &high5Client.PostV1OrgOrgNameSpacesParams{
			Authorization: jwt,
		}

		if cmd.Flags().Changed("noAvatar") {
			noAvatar, err := cmd.Flags().GetBool("noAvatar")
			if err != nil {
				return err
			}
			params.NoAvatar = &noAvatar
		}

		name, err := command.GetArg(cmd, command.StringFlagIfSet("name"), command.ContextIfValidString(latest.High5SpaceName), command.PromptString("name"))
		if err != nil {
			return err
		}
		latest.High5SpaceName = sql.NullString{Valid: true, String: name}

		body := high5Client.PostV1OrgOrgNameSpacesJSONRequestBody{
			Name: name,
		}

		if cmd.Flags().Changed("enableDev") {
			enableDev, err := cmd.Flags().GetBool("enableDev")
			if err != nil {
				return err
			}
			body.EnableDev = &enableDev
		}

		res, err := client.PostV1OrgOrgNameSpacesWithResponse(context.TODO(), orgName, params, body)
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
			if res.JSON201 != nil {
				latest.High5SpaceID = sql.NullString{Valid: true, String: res.JSON201.UnderscoreId}
			}
			if res.JSON201 != nil {
				latest.High5SpaceName = sql.NullString{Valid: true, String: res.JSON201.Name}
			}
			service.Push(latest)
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var spaceGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Get a High5 space by its name",
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

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameWithResponse(context.TODO(), orgName, spaceName, &high5Client.GetV1OrgOrgNameSpacesSpaceNameParams{
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
