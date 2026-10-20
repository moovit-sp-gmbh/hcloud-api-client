package idp

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/idpClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	userPingCommand.Flags().String("token", "", "authorization token")
	patchCommand.Flags().String("name", "", "New user name")
	patchPwCommand.Flags().String("old", "", "Old password")
	patchPwCommand.Flags().String("new", "", "New password")
	userDeleteCommand.Flags().String("token", "", "authorization token")
	userSessionsCommand.Flags().String("token", "", "authorization token")
	whoamiCommand.AddCommand(patchCommand)
	whoamiCommand.AddCommand(getUserEmailCommand)
	whoamiCommand.AddCommand(getUserNatsPermissionsCommand)
	whoamiCommand.AddCommand(patchPwCommand)
	whoamiCommand.AddCommand(userLicenseCommand)
	whoamiCommand.AddCommand(userSettingsCommand)
	whoamiCommand.AddCommand(userPingCommand)
	whoamiCommand.AddCommand(userDeleteCommand)
	whoamiCommand.AddCommand(userOrgsCommand)
	whoamiCommand.AddCommand(userPasskeyCommand)
	whoamiCommand.AddCommand(userPasswordCommand)
	whoamiCommand.AddCommand(userSessionsCommand)

}

var whoamiCommand = &cobra.Command{
	Use:   "user",
	Short: "Get user info",
	Long:  "Get user info",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		res, err := client.GetV1UserWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}

		var text command.Text
		switch res.StatusCode() {
		case 200:

			text, err = command.JSON(res.JSON200)
			if err != nil {
				return err
			}
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

var getUserEmailCommand = &cobra.Command{
	Use:   "email",
	Short: "Get user email",
	Long:  "Get user email",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)

		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		res, err := client.GetV1UserEmailWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}

		var text command.Text
		switch res.StatusCode() {
		case 200:

			text, err = command.JSON(res.JSON200)
			if err != nil {
				return err
			}
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

var getUserNatsPermissionsCommand = &cobra.Command{
	Use:   "nats",
	Short: "Get user Nats permissions",
	Long:  "Get user Nats permissions",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		res, err := client.GetV1UserNatsWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}
		var text command.Text

		switch res.StatusCode() {
		case 200:

			text, err = command.JSON(res.JSON200)
			if err != nil {
				return err
			}
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

var patchCommand = &cobra.Command{
	Use:   "patch",
	Short: "Patch user info",
	Long:  "Patch user info",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		newName, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		res, err := client.PatchV1UserWithResponse(context.TODO(), idpClient.PatchV1UserJSONRequestBody{
			Name: &newName,
		}, func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}

		ctx.NewName = sql.NullString{Valid: true, String: res.JSON200.Name}
		err = service.Push(ctx)
		if err != nil {
			log.Println(err.Error())
		}

		text, err := command.JSON(res.JSON200)
		if err != nil {
			return err
		}
		command.SetResult(cmd, text)

		return nil
	},
}

var patchPwCommand = &cobra.Command{
	Use:   "patch-password",
	Short: "Patch user password",
	Long:  "Patch user password",
	RunE: func(cmd *cobra.Command, args []string) error {

		server, err := command.ServerIDP(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server)
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		oldPw, err := cmd.Flags().GetString("old")
		if err != nil {
			return err
		}

		newPw, err := cmd.Flags().GetString("new")
		if err != nil {
			return err
		}

		res, err := client.PatchV1UserPasswordWithResponse(context.TODO(), idpClient.PatchV1UserPasswordJSONRequestBody{
			Old: oldPw,
			New: newPw,
		}, func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})

		if err != nil {
			return err
		}

		var text command.Text
		switch res.StatusCode() {
		case 200:
			text, err = command.JSON(res.JSON200)
			if err != nil {
				return err
			}
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

var userPingCommand = &cobra.Command{
	Use:   "ping",
	Short: "Keep user online",
	Long:  "\nRequired scopes: HCLOUD_FULL",
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

		res, err := client.GetV1UserPingWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
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

var userDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete the current request user",
	Long:  "Required scopes: IDP_USER_DELETE",
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

		res, err := client.DeleteV1UserWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
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

var userSessionsCommand = &cobra.Command{
	Use:   "sessions",
	Short: "Resets all user sessions and therewith logs him out of hcloud on all devices",
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

		res, err := client.DeleteV1UserSessionsWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
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
