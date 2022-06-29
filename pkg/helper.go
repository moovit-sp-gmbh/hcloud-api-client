package pkg

import (
	"encoding/json"
	"os"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
)

func GetHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		PrintErr(&hcloud.ErrorResponse{Code: "000.000.0000", Error: "client.get.homedir", Message: err.Error()})
	}

	return home
}

func ParseError(body string) *hcloud.ErrorResponse {
	resp := &hcloud.ErrorResponse{}
	err := json.Unmarshal([]byte(body), resp)
	if err != nil {
		PrintErr(&hcloud.ErrorResponse{Code: "000.000.0000", Error: "client.parse.error", Message: err.Error()})
	}
	return resp
}
