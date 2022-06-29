package config

import (
	"fmt"
	"hcloud-api-client/pkg"
	"io/ioutil"
	"os"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"

	"gopkg.in/yaml.v3"
)

var configPath = pkg.GetHomeDir() + "/.hcloud"

var Config Configuration

type Configuration struct {
	Active   *ContextEntry  `json:"active" yml:"active"`
	Contexts []ContextEntry `json:"context" yml:"context"`
}

func (c Configuration) String() (res string) {
	res = ""
	res += fmt.Sprintf("Active:\n%s\n\n", c.Active)
	for _, ce := range c.Contexts {
		res += fmt.Sprintf("%s\n", ce)
	}

	return
}

type ContextEntry struct {
	Identifier string `json:"identifier" yml:"identifier"`
	Server     string `json:"server" yml:"server"`
	Token      string `json:"token" yml:"token"`
	Email      string `json:"email" yml:"email"`
}

func (c ContextEntry) String() string {
	return fmt.Sprintf("Identifier: %s\nServer: %s\nEmail: %s\n", c.Identifier, c.Server, c.Email)
}

func LoadConfig() error {
	c, err := ioutil.ReadFile(configPath + "/config")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(c, &Config)
	if err != nil {
		pkg.PrintErr(&hcloud.ErrorResponse{Code: "000.000.0000", Error: "client.config.parse", Message: err.Error()})
	}

	return nil
}

func (c *Configuration) GetActiveContext() *ContextEntry {
	if Config.Active == nil {
		pkg.PrintErr(&hcloud.ErrorResponse{Code: "000.000.0000", Error: "client.config.context", Message: "no active context found - authenticate first to create a new context (hcloud service idp authenticate -h)"})
	}

	return Config.Active
}

func AddContext(identifier string, server string, email string, token string) {
	for _, config := range Config.Contexts {
		if config.Identifier == identifier {
			config.Server = server
			config.Email = email
			config.Token = token
			writeConfig()
			return
		}
	}

	Config.Contexts = append(Config.Contexts, ContextEntry{Identifier: identifier, Server: server, Email: email, Token: token})

	if len(Config.Contexts) == 1 {
		Config.Active = &Config.Contexts[0]
	}

	writeConfig()
}

func SetContext(identifier string) {
	for _, config := range Config.Contexts {
		if config.Identifier == identifier {
			Config.Active = &config
			writeConfig()
			return
		}
	}

	pkg.PrintErr(&hcloud.ErrorResponse{Code: "000.000.0000", Error: "client.config.context", Message: fmt.Sprintf("context '%s' not found", identifier)})
}

func DelContext(identifier string) {
	for i, config := range Config.Contexts {
		if config.Identifier == identifier {
			Config.Contexts = remove(Config.Contexts, i)
			if Config.Active.Identifier == identifier {
				if len(Config.Contexts) >= 1 {
					Config.Active = &Config.Contexts[0]
				} else {
					Config.Active = nil
				}
			}
			writeConfig()
			return
		}
	}

	pkg.PrintErr(&hcloud.ErrorResponse{Code: "000.000.0000", Error: "client.client.context", Message: fmt.Sprintf("context '%s' not found", identifier)})
}

func remove(slice []ContextEntry, s int) []ContextEntry {
	copy(slice[s:], slice[s+1:])
	slice[len(slice)-1] = ContextEntry{}
	return slice[:len(slice)-1]
}

func writeConfig() {
	os.MkdirAll(pkg.GetHomeDir()+"/.hcloud", 0700)

	b, err := yaml.Marshal(Config)
	if err != nil {
		pkg.PrintErr(&hcloud.ErrorResponse{Code: "000.000.0000", Error: "client.config.write", Message: err.Error()})
	}

	err = ioutil.WriteFile(pkg.GetHomeDir()+"/.hcloud/config", b, 0700)
	if err != nil {
		pkg.PrintErr(&hcloud.ErrorResponse{Code: "000.000.0000", Error: "client.config.write", Message: err.Error()})
	}
}
