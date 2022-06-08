package main

import (
	"hcloud-api-client/cmd"
	"hcloud-api-client/config"
	"os"
)

func main() {
	config.LoadConfig()

	cmd.Execute(os.Args)
}
