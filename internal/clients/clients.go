package clients

import (
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/cosmoClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/daliClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/high5Client"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/idpClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/mothershipClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func NewIDPClient(cmd *cobra.Command) (*idpClient.ClientWithResponses, error) {
	server, err := command.ServerIDP(cmd)
	if err != nil {
		return nil, err
	}

	client, err := idpClient.NewClientWithResponses(server)
	if err != nil {
		return nil, err
	}

	return client, nil
}
func NewHigh5Client(cmd *cobra.Command) (*high5Client.ClientWithResponses, error) {
	server, err := command.ServerHigh5(cmd)
	if err != nil {
		return nil, err
	}

	client, err := high5Client.NewClientWithResponses(server)
	if err != nil {
		return nil, err
	}

	return client, nil
}
func NewDaliClient(cmd *cobra.Command) (*daliClient.ClientWithResponses, error) {
	server, err := command.ServerDali(cmd)
	if err != nil {
		return nil, err
	}

	client, err := daliClient.NewClientWithResponses(server)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewMothershipClient(cmd *cobra.Command) (*mothershipClient.ClientWithResponses, error) {
	server, err := command.ServerMothership(cmd)
	if err != nil {
		return nil, err
	}

	client, err := mothershipClient.NewClientWithResponses(server)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewCosmoClient(cmd *cobra.Command) (*cosmoClient.ClientWithResponses, error) {
	server, err := command.ServerCosmo(cmd)
	if err != nil {
		return nil, err
	}

	client, err := cosmoClient.NewClientWithResponses(server)
	if err != nil {
		return nil, err
	}

	return client, nil
}
