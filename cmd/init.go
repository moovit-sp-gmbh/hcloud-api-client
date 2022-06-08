package cmd

import (
	"fmt"
	"hcloud-api-client/cmd/config"
	"hcloud-api-client/cmd/service"
	"hcloud-api-client/cmd/update"
	"hcloud-api-client/cmd/version"
	"hcloud-api-client/pkg"
	"os"

	"github.com/spf13/cobra"
)

var format string

var rootCmd = &cobra.Command{
	Use:               "hcloud",
	Short:             "hcloud is an command line interface for the helmut.cloud platform",
	Long:              "hcloud is an command line interface to communicate with the helmut.cloud platform API",
	CompletionOptions: cobra.CompletionOptions{},
}

func Execute(args []string) {
	rootCmd.DisableAutoGenTag = true

	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "plain", "set the output format to plain or json")
	rootCmd.ParseFlags(args)
	pkg.SetFormat(format)

	update.Init(rootCmd)
	service.Init(rootCmd)
	version.Init(rootCmd)
	config.Init(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
