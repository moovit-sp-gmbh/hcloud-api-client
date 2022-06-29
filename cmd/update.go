package cmd

import (
	"encoding/json"
	"fmt"
	"hcloud-api-client/pkg"
	"io"
	"net/http"

	"github.com/hashicorp/go-version"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "check for updates",
	Run:   checkForUpdates,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func checkForUpdates(cmd *cobra.Command, args []string) {
	resp, err := http.Get("https://api.github.com/repos/moovit-sp-gmbh/hcloud-api-client/releases/latest")
	if err != nil {
		pkg.PrintErr(&hcloud.ErrorResponse{Code: "000.000.0000", Error: "client.http.error", Message: err.Error()})
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		pkg.PrintErr(&hcloud.ErrorResponse{Code: "000.000.0000", Error: "client.http.error", Message: err.Error()})
	}

	type GithubRelease struct {
		Tag string `json:"tag_name"`
	}

	t := &GithubRelease{}
	err = json.Unmarshal(b, t)
	if err != nil {
		pkg.PrintErr(&hcloud.ErrorResponse{Code: "000.000.0000", Error: "client.parse.error", Message: err.Error()})
	}

	remote, _ := version.NewVersion(t.Tag)
	local, _ := version.NewVersion(hcloudVersion)

	if remote.GreaterThan(local) {
		pkg.Print(pkg.OkResponse{Result: fmt.Sprintf("a new version is available (current: %s, latest: %s)\nVisit https://github.com/moovit-sp-gmbh/hcloud-api-client/releases/tag/v%s to download it", local.String(), remote.String(), remote.String())})
	} else {
		pkg.Print(pkg.OkResponse{Result: fmt.Sprintf("local verison is up to date (current: v%s, latest: v%s)", local.String(), remote.String())})
	}
}
