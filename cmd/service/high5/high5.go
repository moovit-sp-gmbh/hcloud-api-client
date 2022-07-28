package high5

import (
	"hcloud-api-client/cmd/service/high5/app"
	"hcloud-api-client/cmd/service/high5/event"
	"hcloud-api-client/cmd/service/high5/execute"
	"hcloud-api-client/cmd/service/high5/stream"

	"github.com/spf13/cobra"
)

var high5Cmd = &cobra.Command{
	Use:   "high5",
	Short: "use helmut.cloud high5 (stream engine)",
}

func Init(serviceCmd *cobra.Command) {
	serviceCmd.AddCommand(high5Cmd)

	// init app
	app.Init(high5Cmd)

	// init event
	event.Init(high5Cmd)

	// init stream
	stream.Init(high5Cmd)

	// init execute
	execute.Init(high5Cmd)
}
