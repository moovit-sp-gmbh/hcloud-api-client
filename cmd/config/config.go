package config

import (
	"hcloud-api-client/cmd/config/context"
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/spf13/cobra"
)

type ConfigPath struct {
	Path string `json:"path"`
}

func (c ConfigPath) String() string {
	return c.Path
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "manage the local configuration file",
}

var identifier string

var configCmds = []*cobra.Command{
	{
		Use:   "path",
		Short: "print local configuration path",
		Run:   printLocalConfigPath,
	},
	{
		Use:   "print",
		Short: "print local configuration",
		Run:   printLocalConfig,
	},
}

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(configCmd)

	for _, cm := range configCmds {
		configCmd.AddCommand(cm)
	}

	// init context
	context.Init(configCmd)
}

func printLocalConfigPath(cmd *cobra.Command, args []string) {
	pkg.Print(ConfigPath{Path: pkg.GetHomeDir() + "/.hcloud/config"})
}

func printLocalConfig(cmd *cobra.Command, args []string) {
	pkg.Print(config.Config)
}
