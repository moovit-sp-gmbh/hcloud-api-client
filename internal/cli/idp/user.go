package idp

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/idpClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	patchCommand.Flags().String("name", "", "New user name")
	patchPwCommand.Flags().String("old", "", "Old password")
	patchPwCommand.Flags().String("new", "", "New password")
	whoamiCommand.AddCommand(patchCommand)
	whoamiCommand.AddCommand(getUserEmailCommand)
	whoamiCommand.AddCommand(getUserNatsPermissionsCommand)
	whoamiCommand.AddCommand(patchPwCommand)

	whoamiCommand.AddCommand(pingCommand)
	whoamiCommand.AddCommand(logoutCommand)

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


var pingCommand = &cobra.Command{
	Use:   "ping",
	Short: "Ping to keep the user session alive",
	Long:  "Notify the server that the user is online, keeping the session alive",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := command.Server(cmd)
		if err != nil {
			return err
		}

		client, err := idpClient.NewClientWithResponses(server + "/api/account")
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

		res, err := client.GetV1UserPing(context.TODO(), func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", ctx.Jwt.String)
			return nil
		})
		if err != nil {
			return err
		}
		defer res.Body.Close()

		var text command.Text

		if res.StatusCode != http.StatusNoContent && (res.StatusCode < 200 || res.StatusCode >= 300) {
			data, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}
			text, err := command.JSON(string(data))
			if err != nil {
				return err
			}
			return errors.New(string(text))
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var logoutCommand = &cobra.Command{
    Use:   "logout",
    Short: "Log out of hcloud on all devices",
    Long:  "Resets all user sessions and therewith logs the user out of hcloud on all devices",
    RunE: func(cmd *cobra.Command, args []string) error {
        server, err := command.Server(cmd)
        if err != nil {
            return err
        }

        client, err := idpClient.NewClientWithResponses(server + "/api/account")
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

        res, err := client.DeleteV1UserSessionsWithResponse(context.TODO(), func(_ context.Context, req *http.Request) error {
            req.Header.Set("authorization", ctx.Jwt.String)
            return nil
        })
        if err != nil {
            return err
        }

		log.Println(res.StatusCode(), string(res.Body))

		var text command.Text

		if res.StatusCode() != http.StatusNoContent && (res.StatusCode() < 200 || res.StatusCode() >= 300) {
			text, err := command.JSON(string(res.Body))
			if err != nil {
				return err
			}
			return errors.New(string(text))
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
		log.Println(res.JSON200.Name)

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
