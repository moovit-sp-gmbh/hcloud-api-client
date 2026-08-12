package mothership

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/cli"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/mothershipClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/spf13/cobra"
)

func init() {
	helloCommand.AddCommand(helloAgainCommand)

	helloCommand.Flags().String("secret", "", "secret gotten from the register phase")
	helloCommand.Flags().String("uuid", "", "UUID used in the register phase")

	helloAgainCommand.Flags().String("token", "", "token gotten from the hello phase")
}

var helloCommand = &cobra.Command{
	Use:   "hello",
	Short: "Say hello to the mothership to start reporting",
	RunE: func(cmd *cobra.Command, args []string) error {
		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		var secret string
		if cmd.Flags().Changed("secret") {
			secret, err = cmd.Flags().GetString("secret")
			if err != nil {
				return err
			}
		} else {
			if ctx.AgentSecret.Valid {
				secret = ctx.AgentSecret.String
			} else {
				secret, err = cli.Prompt("secret")
				if err != nil {
					return err
				}
			}
		}

		var uuidInst uuid.UUID
		if cmd.Flags().Changed("uuid") {
			uuidStr, err := cmd.Flags().GetString("uuid")
			if err != nil {
				return err
			}
			uuidInst, err = uuid.Parse(uuidStr)
			if err != nil {
				return err
			}
		} else {
			if ctx.AgentUuid.Valid {
				uuidInst = uuid.MustParse(ctx.AgentUuid.String)
			} else {
				uuidStr, err := cli.Prompt("uuid")
				if err != nil {
					return err
				}
				uuidInst, err = uuid.Parse(uuidStr)
				if err != nil {
					return err
				}
			}
		}

		client, err := clients.NewMothershipClient(cmd)
		if err != nil {
			return err
		}

		cpuUtilization, err := cpu.Percent(time.Second/2, false)
		if err != nil {
			return err
		}

		vm, err := mem.VirtualMemory()
		if err != nil {
			return err
		}

		memoryUsed := vm.Used

		uptime, err := host.Uptime()
		if err != nil {
			return err
		}

		secretBytes, err := base64.StdEncoding.DecodeString(secret)
		if err != nil {
			return err
		}
		mac := hmac.New(sha256.New, secretBytes)
		mac.Write([]byte(uuidInst.String()))
		sum := mac.Sum(nil)

		uuidSignature := base64.StdEncoding.EncodeToString(sum)

		res, err := client.PostV1HelloWithResponse(context.TODO(), mothershipClient.AgentHelloDto{
			CpuUtilization: float32(cpuUtilization[0]),
			MemoryUsed:     float32(memoryUsed),
			Status:         mothershipClient.AgentStatus("IDLE"),
			Uptime:         float32(uptime),
		}, func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", fmt.Sprintf("Bearer %s.%s", uuidInst, uuidSignature))
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
			ctx.AgentHelloToken = sql.NullString{Valid: true, String: res.HTTPResponse.Header.Get("authorization")}
			service.Push(ctx)
		case 401:
			text, err = command.JSON(res.JSON401)
		case 500:
			text, err = command.JSON(res.JSON500)
		}
		if err != nil {
			return err
		}

		command.SetResult(cmd, text)

		return nil
	},
}

var helloAgainCommand = &cobra.Command{
	Use:   "again",
	Short: "Say hello again to the mothership to keep reporting",
	RunE: func(cmd *cobra.Command, args []string) error {
		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		client, err := clients.NewMothershipClient(cmd)
		if err != nil {
			return err
		}

		var token string
		if cmd.Flags().Changed("token") {
			token, err = cmd.Flags().GetString("token")
			if err != nil {
				return err
			}
		} else {
			if ctx.AgentHelloToken.Valid {
				token = ctx.AgentHelloToken.String
			} else {
				token, err = cli.Prompt("token")
				if err != nil {
					return err
				}
			}
		}

		cpuUtilization, err := cpu.Percent(time.Second/2, false)
		if err != nil {
			return err
		}

		vm, err := mem.VirtualMemory()
		if err != nil {
			return err
		}

		memoryUsed := vm.Used

		uptime, err := host.Uptime()
		if err != nil {
			return err
		}

		res, err := client.PatchV1HelloAgainWithResponse(context.TODO(), mothershipClient.AgentHelloAgainDto{
			CpuUtilization: float32(cpuUtilization[0]),
			MemoryUsed:     float32(memoryUsed),
			Status:         mothershipClient.AgentStatus("IDLE"),
			Uptime:         float32(uptime),
		}, func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", token)
			return nil
		})
		if err != nil {
			return err
		}

		var data any
		err = json.Unmarshal([]byte(res.Body), &data)
		if err != nil {
			return err
		}
		text, err := command.JSON(data)
		if err != nil {
			return err
		}

		command.SetResult(cmd, text)

		return nil
	},
}
