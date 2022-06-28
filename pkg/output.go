package pkg

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
)

var format = "plain"

type OkResponse struct {
	Result string `json:"result"`
}

func (o OkResponse) String() string {
	return o.Result
}

func Print(msg interface{}) {
	if format == "json" {
		b, err := json.Marshal(msg)
		if err != nil {
			PrintErr(&hcloud.ErrorResponse{Code: -1, Message: err.Error()})
		}
		fmt.Printf("%s", string(b))
	} else {
		fmt.Printf("%s\n", msg)
	}
	os.Exit(0)
}

func PrintErr(err *hcloud.ErrorResponse) {
	if format == "json" {
		b, _ := json.Marshal(err)
		fmt.Fprintf(os.Stderr, "%s\n", string(b))
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", err.ToString())
	}
	os.Exit(1)
}

func SetFormat(fmt string) {
	format = fmt
}
