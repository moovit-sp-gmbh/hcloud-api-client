package pkg

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hokaccha/go-prettyjson"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"gopkg.in/yaml.v2"
)

var format = "plain"

type OkResponse struct {
	Result string `json:"result"`
}

func (o OkResponse) String() string {
	return o.Result
}

func Print(msg interface{}) {
	PrintContinously(msg)
	os.Exit(0)
}

func PrintContinously(msg interface{}) {
	if format == "json" {
		b, err := json.Marshal(msg)
		if err != nil {
			PrintErr(&hcloud.ErrorResponse{Code: "000.000.0000", Error: "client.parse.json", Message: err.Error()})
		}
		fmt.Printf("%s", string(b))
	} else if format == "json-indent" {
		b, err := json.MarshalIndent(msg, "", "    ")
		if err != nil {
			PrintErr(&hcloud.ErrorResponse{Code: "000.000.0000", Error: "client.parse.json", Message: err.Error()})
		}
		fmt.Printf("%s", string(b))
	} else if format == "json-pretty" {
		b, err := prettyjson.Marshal(msg)
		if err != nil {
			PrintErr(&hcloud.ErrorResponse{Code: "000.000.0000", Error: "client.parse.json", Message: err.Error()})
		}
		fmt.Printf("%s", string(b))
	} else if format == "yaml" {
		yamlData, err := yaml.Marshal(&msg)
		if err != nil {
			PrintErr(&hcloud.ErrorResponse{Code: "000.000.0000", Error: "client.parse.yaml", Message: err.Error()})
		}
		fmt.Printf("%s\n", yamlData)
	} else {
		fmt.Printf("%s\n", msg)
	}
}

func PrintErr(err *hcloud.ErrorResponse) {
	if format == "json" {
		b, _ := json.Marshal(err)
		fmt.Fprintf(os.Stderr, "%s\n", string(b))
	} else if format == "json-indent" {
		b, _ := json.MarshalIndent(err, "", "    ")
		fmt.Fprintf(os.Stderr, "%s\n", string(b))
	} else if format == "json-pretty" {
		b, _ := prettyjson.Marshal(err)
		fmt.Fprintf(os.Stderr, "%s\n", string(b))
	} else if format == "yaml" {
		yamlData, _ := yaml.Marshal(&err)
		fmt.Fprintf(os.Stderr, "%s\n", yamlData)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", err.ToString())
	}
	os.Exit(1)
}

func SetFormat(fmt string) {
	format = fmt
}
