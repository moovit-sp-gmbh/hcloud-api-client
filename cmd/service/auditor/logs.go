package auditor

import (
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/auditor"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "request audit logs from the auditor",
	Run:   requestLogs,
}

func init() {
	logsCmd.PersistentFlags().StringVarP(&organization, "organization", "o", "", "the id of the organization to get logs from")
	logsCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 500, "the maximum amount of logs")
	logsCmd.PersistentFlags().IntVarP(&page, "page", "p", 0, "the amount of entries to skip (page * limit)")
	auditorCmd.AddCommand(logsCmd)
}

func requestLogs(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	auditor := auditor.New(hcloud.New(&hcloud.ClientConfig{Api: ctx.Server, Token: ctx.Token}))

	logs, _, err := auditor.GetAuditLogs(organization, limit, page)
	if err != nil {
		pkg.PrintErr(err)
	}

	for _, log := range *logs {
		pkg.PrintContinously(log)
	}

}
