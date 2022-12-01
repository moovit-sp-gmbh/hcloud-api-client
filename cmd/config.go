package cmd

import (
	"hcloud-api-client/cmd/config/context"
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "manage the local configuration file",
}

type ConfigPath struct {
	Path string `json:"path"`
}

func (c ConfigPath) String() string {
	return c.Path
}

var identifier string

var configSubCmd = []*cobra.Command{
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

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSubCmd...)

	// init context subcommand
	context.Init(configCmd)
}

func printLocalConfigPath(cmd *cobra.Command, args []string) {
	pkg.Print(ConfigPath{Path: pkg.GetHomeDir() + "/.hcloud/config.yml"})
}

func printLocalConfig(cmd *cobra.Command, args []string) {
	pkg.Print(config.Config)
}
