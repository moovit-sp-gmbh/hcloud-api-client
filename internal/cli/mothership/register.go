package mothership

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/mothershipClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/spf13/cobra"
)

func init() {
	registerCommand.Flags().String("cpu", "", "")
	registerCommand.Flags().Float64("cpuUtilization", 0, "")
	registerCommand.Flags().String("hostname", "", "")
	registerCommand.Flags().Float64("memoryTotal", 0, "")
	registerCommand.Flags().Float64("memoryUsed", 0, "")
	registerCommand.Flags().String("nickname", "", "")
	registerCommand.Flags().String("osPlatform", "", "")
	registerCommand.Flags().String("osRelease", "", "")
	registerCommand.Flags().String("osVersion", "", "")
	registerCommand.Flags().String("publicKey", "", "")
	registerCommand.Flags().String("status", "", "")
	registerCommand.Flags().Float64("uptime", 0, "")
	registerCommand.Flags().String("uuid", "", "")
}

var registerCommand = &cobra.Command{
	Use:   "register",
	Short: "Register an agent on the mothership",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewMothershipClient(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		cpuInfo, err := cpu.Info()
		if err != nil {
			return err
		}
		cpuModel := cpuInfo[0].ModelName
		cpuUtilization, err := cpu.Percent(time.Second/2, false)
		if err != nil {
			return err
		}

		hostname, err := os.Hostname()
		vm, err := mem.VirtualMemory()
		if err != nil {
			return err
		}

		memoryTotal := vm.Total
		memoryUsed := vm.Used

		osPlatform, osRelease, osVersion, err := host.PlatformInformation()
		if err != nil {
			return err
		}

		uptime, err := host.Uptime()
		if err != nil {
			return err
		}

		nickname, err := command.FlagOrPromptString(cmd, "nickname")
		if err != nil {
			return err
		}

		publicKey, privateKey, err := generateRSAKeyPair(4096)
		if err != nil {
			return err
		}
		latest.AgentPrivateKey = sql.NullString{Valid: true, String: privateKey}
		latest.AgentPublicKey = sql.NullString{Valid: true, String: publicKey}

		status, err := command.FlagOrPromptString(cmd, "status")
		if err != nil {
			return err
		}

		var uuidValid uuid.UUID
		if cmd.Flags().Changed("uuid") {
			uuidStr, err := cmd.Flags().GetString("uuid")
			if err != nil {
				return err
			}
			uuidValid, err = uuid.Parse(uuidStr)
			if err != nil {
				return err
			}
		} else {
			uuidValid = uuid.New()
		}

		res, err := client.PostV1RegisterWithResponse(context.TODO(), mothershipClient.AgentRegisterDto{
			Cpu:            cpuModel,
			CpuUtilization: float32(cpuUtilization[0]),
			Hostname:       hostname,
			MemoryTotal:    float32(memoryTotal),
			MemoryUsed:     float32(memoryUsed),
			Nickname:       nickname,
			OsPlatform:     osPlatform,
			OsRelease:      osRelease,
			OsVersion:      osVersion,
			PublicKey:      publicKey,
			Status:         mothershipClient.AgentStatus(status),
			Uptime:         float32(uptime),
			Uuid:           uuidValid,
		})
		if err != nil {
			return err
		}

		var text command.Text
		switch res.StatusCode() {
		case 200, 201:
			var data struct {
				Secret string `json:"secret"`
			}
			err := json.Unmarshal([]byte(res.Body), &data)
			if err != nil {
				return err
			}
			text, err = command.JSON(data)
			if err != nil {
				return err
			}

			encSecret, err := base64.StdEncoding.DecodeString(data.Secret)
			if err != nil {
				return fmt.Errorf("unable to decode secret: %w", err)
			}

			block, _ := pem.Decode([]byte(privateKey))
			key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return fmt.Errorf("unable to parse private key: %w", err)
			}
			pKey, ok := key.(*rsa.PrivateKey)
			if !ok {
				return errors.New("unable to decrypt secret")
			}

			b, err := rsa.DecryptOAEP(sha1.New(), nil, pKey, encSecret, nil)
			if err != nil {
				return fmt.Errorf("unable to decrypt secret: %w", err)
			}
			latest.AgentSecret = sql.NullString{Valid: true, String: base64.StdEncoding.EncodeToString(b)}

			latest.AgentUuid = sql.NullString{Valid: true, String: uuidValid.String()}
			service.Push(latest)
		default:
			var data any
			err := json.Unmarshal([]byte(res.Body), &data)
			if err != nil {
				return err
			}
			text, err = command.JSON(data)
			if err != nil {
				return err
			}
		}

		command.SetResult(cmd, text)

		return nil
	},
}

func generateRSAKeyPair(modulusLength int) (string, string, error) {
	// Generate the private key
	privateKey, err := rsa.GenerateKey(rand.Reader, modulusLength)
	if err != nil {
		return "", "", err
	}

	// Encode private key as PKCS8 PEM (matches Node's privateKeyEncoding: pkcs8)
	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", err
	}
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	})

	// Encode public key as SPKI PEM (matches Node's publicKeyEncoding: spki)
	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})

	return string(pubPEM), string(privPEM), nil
}
