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

var executeStreamCmd = &cobra.Command{
	Use:   "stream",
	Short: "execute a stream by it's ID",
	Run:   executeStream,
}

func init() {
	executeStreamCmd.PersistentFlags().StringVarP(&id, "streamId", "i", "", "the id of the stream")
	executeStreamCmd.MarkPersistentFlagRequired("streamId")

	executeStreamCmd.PersistentFlags().StringVarP(&target, "target", "t", "", "the name of the target to execute the stream on")
	executeStreamCmd.MarkPersistentFlagRequired("target")

	executeStreamCmd.PersistentFlags().StringVarP(&dataType, "type", "y", "JSON", "the type of data to be send as payload for the stream, JSON or GENERIC")
	executeStreamCmd.PersistentFlags().StringVarP(&data, "data", "d", "", "the data to be send as payload for the stream")
	if runtime.GOOS != "windows" {
		executeStreamCmd.PersistentFlags().BoolVarP(&dataStdin, "data-stdin", "", false, "read the data to be send as payload for the stream from stdin")
	}
	executeStreamCmd.PersistentFlags().BoolVarP(&wait, "wait", "w", true, "wait for the response of the execution")
	executeStreamCmd.PersistentFlags().IntVarP(&timeout, "timeout", "", 10000, "the maximum time to wait for a successful execution in ms")

	executeCmd.AddCommand(executeStreamCmd)
}

func executeStream(cmd *cobra.Command, args []string) {
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.Config{Api: ctx.Server, Token: ctx.Token}))

	if dataStdin {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			data += scanner.Text()
		}
	}

	if data == "" {
		data = "{}"
	}

	apps, err := high5.ExecuteStreamById(id, target, dataType, []byte(data), timeout, wait)
	if err != nil {
		pkg.PrintErr(err)
	}

	pkg.Print(apps)
}
