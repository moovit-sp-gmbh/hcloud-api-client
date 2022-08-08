package execute

import (
	"bufio"
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"
	"os"
	"runtime"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/spf13/cobra"
)

var executeEventCmd = &cobra.Command{
	Use:   "event",
	Short: "execute an event by it's NAME",
	Run:   executeEvent,
}

func init() {
	executeEventCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "the id of the app")
	executeEventCmd.MarkPersistentFlagRequired("id")

	executeEventCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "the name of the event to execute")
	executeEventCmd.MarkPersistentFlagRequired("name")

	executeEventCmd.PersistentFlags().StringVarP(&target, "target", "t", "", "the name of the target to execute the stream on")
	executeEventCmd.MarkPersistentFlagRequired("target")

	executeEventCmd.PersistentFlags().StringVarP(&data, "data", "d", "", "the data to be send as payload for the stream")
	if runtime.GOOS != "windows" {
		executeEventCmd.PersistentFlags().BoolVarP(&dataStdin, "data-stdin", "", false, "read the data to be send as payload for the stream from stdin")
	}
	executeEventCmd.PersistentFlags().BoolVarP(&wait, "wait", "w", true, "wait for the response of the execution")
	executeEventCmd.PersistentFlags().IntVarP(&timeout, "timeout", "", 10000, "the maximum time to wait for a successful execution in ms")

	executeCmd.AddCommand(executeEventCmd)
}

func executeEvent(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.ClientConfig{Api: ctx.Server, Token: ctx.Token}))

	if dataStdin {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			data += scanner.Text()
		}
	}

	if data == "" {
		data = "{}"
	}

	apps, err := high5.ExecuteEventByName(id, name, target, []byte(data), timeout, wait)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(apps)
}
