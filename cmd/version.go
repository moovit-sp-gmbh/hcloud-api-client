package cmd

import (
	_ "embed"
	"fmt"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"

	"github.com/spf13/cobra"
)

type Version struct {
	Version string `json:"version"`
}

func (v Version) String() string {
	return fmt.Sprint(v.Version)
}

//go:embed version.txt
var hcloudVersion string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output the current version",
}
var localVersionCmd = &cobra.Command{
	Use:   "local",
	Short: "output local hcloud client version",
	Run:   localVersion,
}
var remoteVersionCmd = &cobra.Command{
	Use:   "remote",
	Short: "output remote hcloud platform version",
	Run:   remoteVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.AddCommand(localVersionCmd)
	versionCmd.AddCommand(remoteVersionCmd)
}

func localVersion(cmd *cobra.Command, args []string) {
	pkg.Print(Version{Version: hcloudVersion})
}

func remoteVersion(cmd *cobra.Command, args []string) {
	pkg.PrintErr(&hcloud.ErrorResponse{Code: -1, Message: "Not implemented"})
}
