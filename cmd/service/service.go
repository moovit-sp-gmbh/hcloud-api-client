package service

import (
	"hcloud-api-client/cmd/service/idp"

	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "use a helmut.cloud service",
}

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(serviceCmd)
	idp.Init(serviceCmd)
}
