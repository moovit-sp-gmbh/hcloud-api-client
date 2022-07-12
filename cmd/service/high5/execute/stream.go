package execute

import (
	"fmt"
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"
	"runtime"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/spf13/cobra"
)

var executeStreamCmd = &cobra.Command{
	Use:   "stream",
	Short: "execute a stream by it's ID",
	Run:   executeStream,
}

func init() {
	executeStreamCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "the id of the stream")
	executeStreamCmd.MarkPersistentFlagRequired("id")

	executeStreamCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "the name of the target to execute the stream on")
	executeStreamCmd.MarkPersistentFlagRequired("name")

	executeStreamCmd.PersistentFlags().StringVarP(&data, "data", "d", "", "the data to be send as payload for the stream")
	if runtime.GOOS != "windows" {
		executeStreamCmd.PersistentFlags().BoolVarP(&dataStdin, "data-stdin", "", false, "read the data to be send as payload for the stream from stdin")
	}
	executeStreamCmd.PersistentFlags().BoolVarP(&wait, "wait", "w", true, "wait for the response of the execution")
	executeStreamCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", 10000, "the maximum time to wait for a successful execution in ms")

	executeCmd.AddCommand(executeStreamCmd)
}

func executeStream(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.ClientConfig{Api: ctx.Server, Token: ctx.Token}))

	if dataStdin {
		fmt.Scan(&data)
	}

	if data == "" {
		data = "{}"
	}

	apps, err := high5.ExecuteStreamById(id, name, []byte(data), timeout, wait)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(apps)
}
