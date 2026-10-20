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
	userSettingsGeneralPatchCommand.Flags().String("token", "", "authorization token")
	userSettingsGeneralPatchCommand.Flags().String("dateFormat", "", "New format to use for dates")
	userSettingsGeneralPatchCommand.Flags().String("entrypoint", "", "New entrypoint in helmut.cloud")
	userSettingsGeneralPatchCommand.Flags().String("language", "", "New language")
	userSettingsGeneralPatchCommand.Flags().String("lastUrl", "", "Last URL visited by the user")
	userSettingsGeneralPatchCommand.Flags().String("lastView", "", "Last view the user was in")
	userSettingsGeneralPatchCommand.Flags().String("theme", "", "Theme to use")
	userSettingsGeneralPatchCommand.Flags().String("timezone", "", "IANA timezone")
	userSettingsGeneralLastVisitCommand.Flags().String("token", "", "authorization token")
	userSettingsGeneralLastVisitCommand.Flags().String("lastUrl", "", "Last URL visited by the user")
	userSettingsGeneralLastVisitCommand.Flags().String("lastView", "", "Last view the user was in")
}

var userSettingsGeneralPatchCommand = &cobra.Command{
	Use:   "patch",
	Short: "Patch the current general settings of the logged in user",
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

		body := idpClient.PatchV1UserSettingsGeneralJSONRequestBody{}

		if cmd.Flags().Changed("dateFormat") {
			dateFormat, err := cmd.Flags().GetString("dateFormat")
			if err != nil {
				return err
			}
			dateFormatValue := idpClient.DateFormat(dateFormat)
			body.DateFormat = &dateFormatValue
		}

		if cmd.Flags().Changed("entrypoint") {
			entrypoint, err := cmd.Flags().GetString("entrypoint")
			if err != nil {
				return err
			}
			entrypointValue := idpClient.Entrypoint(entrypoint)
			body.Entrypoint = &entrypointValue
		}

		if cmd.Flags().Changed("language") {
			language, err := cmd.Flags().GetString("language")
			if err != nil {
				return err
			}
			languageValue := idpClient.Language(language)
			body.Language = &languageValue
		}

		if cmd.Flags().Changed("lastUrl") {
			lastUrl, err := cmd.Flags().GetString("lastUrl")
			if err != nil {
				return err
			}
			body.LastUrl = &lastUrl
		}

		if cmd.Flags().Changed("lastView") {
			lastView, err := cmd.Flags().GetString("lastView")
			if err != nil {
				return err
			}
			lastViewValue := idpClient.Views(lastView)
			body.LastView = &lastViewValue
		}

		if cmd.Flags().Changed("theme") {
			theme, err := cmd.Flags().GetString("theme")
			if err != nil {
				return err
			}
			themeValue := idpClient.Theme(theme)
			body.Theme = &themeValue
		}

		if cmd.Flags().Changed("timezone") {
			timezone, err := cmd.Flags().GetString("timezone")
			if err != nil {
				return err
			}
			body.Timezone = &timezone
		}

		res, err := client.PatchV1UserSettingsGeneralWithResponse(context.TODO(), body, func(_ context.Context, req *http.Request) error {
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

var userSettingsGeneralLastVisitCommand = &cobra.Command{
	Use:   "last-visit",
	Short: "Patch the last view and / or last URL of the logged in user",
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

		body := idpClient.PatchV1UserSettingsGeneralLastVisitJSONRequestBody{}

		if cmd.Flags().Changed("lastUrl") {
			lastUrl, err := cmd.Flags().GetString("lastUrl")
			if err != nil {
				return err
			}
			body.LastUrl = &lastUrl
		}

		if cmd.Flags().Changed("lastView") {
			lastView, err := cmd.Flags().GetString("lastView")
			if err != nil {
				return err
			}
			lastViewValue := idpClient.Views(lastView)
			body.LastView = &lastViewValue
		}

		res, err := client.PatchV1UserSettingsGeneralLastVisitWithResponse(context.TODO(), body, func(_ context.Context, req *http.Request) error {
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
