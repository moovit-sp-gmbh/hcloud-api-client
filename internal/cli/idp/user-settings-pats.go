package idp

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/idpClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	userSettingsPatsDeleteAllCommand.Flags().String("token", "", "authorization token")
	userSettingsPatsCreateCommand.Flags().String("token", "", "authorization token")
	userSettingsPatsCreateCommand.Flags().String("body", "{}", "request body as JSON")
	userSettingsPatsSearchCommand.Flags().String("token", "", "authorization token")
	userSettingsPatsSearchCommand.Flags().String("body", "{}", "request body as JSON")
	userSettingsPatsDeleteCommand.Flags().String("patId", "", "")
	userSettingsPatsDeleteCommand.Flags().String("token", "", "authorization token")
	userSettingsPatsGetCommand.Flags().String("patId", "", "")
	userSettingsPatsGetCommand.Flags().String("token", "", "authorization token")
	userSettingsPatsPatchCommand.Flags().String("patId", "", "")
	userSettingsPatsPatchCommand.Flags().String("token", "", "authorization token")
	userSettingsPatsPatchCommand.Flags().Float32("expiration", 0, "UNIX timestamp with new expiration. Use the value -1 to make the token never expire.")
	userSettingsPatsPatchCommand.Flags().String("name", "", "Name of the Personal Access Token (only alphanumeric characters, underscores, hyphens, dots and spaces are allowed; it should contain between 1 and 64 characters)")
	userSettingsPatsPatchCommand.Flags().String("bodyToken", "", "TOTP token. Must be defined if the user has 2FA enabled.")
	userSettingsPatsRegenerateCommand.Flags().String("patId", "", "")
	userSettingsPatsRegenerateCommand.Flags().String("token", "", "authorization token")
	userSettingsPatsRegenerateCommand.Flags().String("bodyToken", "", "TOTP token. Must be defined if the user has 2FA enabled")

	userSettingsPatsCommand.AddCommand(userSettingsPatsDeleteAllCommand)
	userSettingsPatsCommand.AddCommand(userSettingsPatsCreateCommand)
	userSettingsPatsCommand.AddCommand(userSettingsPatsSearchCommand)
	userSettingsPatsCommand.AddCommand(userSettingsPatsDeleteCommand)
	userSettingsPatsCommand.AddCommand(userSettingsPatsGetCommand)
	userSettingsPatsCommand.AddCommand(userSettingsPatsPatchCommand)
	userSettingsPatsCommand.AddCommand(userSettingsPatsRegenerateCommand)
}

var userSettingsPatsCommand = &cobra.Command{
	Use:   "pats",
	Short: "Commands for pats",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var userSettingsPatsDeleteAllCommand = &cobra.Command{
	Use:   "delete-all",
	Short: "Deletes all users pats",
	Long:  "Required scopes: HCLOUD_FULL",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1UserSettingsPatsWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
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

var userSettingsPatsCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Add a new pat for request user",
	Long:  "Required scopes: HCLOUD_FULL",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body idpClient.PostV1UserSettingsPatsJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1UserSettingsPatsWithResponse(context.TODO(), body, func(_ context.Context, req *http.Request) error {
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

var userSettingsPatsSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Search all pats assigned to the request user",
	Long:  "Required scopes: HCLOUD_FULL",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body idpClient.PostV1UserSettingsPatsSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1UserSettingsPatsSearchWithResponse(context.TODO(), body, func(_ context.Context, req *http.Request) error {
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

var userSettingsPatsDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a users pat by ID",
	Long:  "Required scopes: HCLOUD_FULL",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		patId, err := command.FlagOrPromptString(cmd, "patId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1UserSettingsPatsPatIdWithResponse(context.TODO(), patId, func(_ context.Context, req *http.Request) error {
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

var userSettingsPatsGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Get a single pat by ID",
	Long:  "Required scopes: HCLOUD_FULL",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		patId, err := command.FlagOrPromptString(cmd, "patId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.GetV1UserSettingsPatsPatIdWithResponse(context.TODO(), patId, func(_ context.Context, req *http.Request) error {
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

var userSettingsPatsPatchCommand = &cobra.Command{
	Use:   "patch",
	Short: "Change a pats access scopes, expiration or name",
	Long:  "Required scopes: HCLOUD_FULL",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		patId, err := command.FlagOrPromptString(cmd, "patId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		body := idpClient.PatchV1UserSettingsPatsPatIdJSONRequestBody{}

		if cmd.Flags().Changed("expiration") {
			expiration, err := cmd.Flags().GetFloat32("expiration")
			if err != nil {
				return err
			}
			body.Expiration = &expiration
		}

		if cmd.Flags().Changed("name") {
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}
			body.Name = &name
		}

		if cmd.Flags().Changed("bodyToken") {
			bodyToken, err := cmd.Flags().GetString("bodyToken")
			if err != nil {
				return err
			}
			body.Token = &bodyToken
		}

		res, err := client.PatchV1UserSettingsPatsPatIdWithResponse(context.TODO(), patId, body, func(_ context.Context, req *http.Request) error {
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

var userSettingsPatsRegenerateCommand = &cobra.Command{
	Use:   "regenerate",
	Short: "Regenerate an existing pat token with previous defined expiration date",
	Long:  "Required scopes: HCLOUD_FULL",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		patId, err := command.FlagOrPromptString(cmd, "patId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		body := idpClient.PatchV1UserSettingsPatsPatIdRegenerateJSONRequestBody{}

		if cmd.Flags().Changed("bodyToken") {
			bodyToken, err := cmd.Flags().GetString("bodyToken")
			if err != nil {
				return err
			}
			body.Token = &bodyToken
		}

		res, err := client.PatchV1UserSettingsPatsPatIdRegenerateWithResponse(context.TODO(), patId, body, func(_ context.Context, req *http.Request) error {
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
