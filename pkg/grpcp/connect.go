package grpcp

import (
	"fmt"
	"os/exec"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/hclogadapter"
	plugin "github.com/hashicorp/go-plugin"
)

// ConnectPlugin will connect to a plugin over gRPC.
func ConnectPlugin(logger ambient.AppLogger, pluginName string, pluginPath string) (ambient.Plugin, *plugin.Client, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]plugin.Plugin{
			pluginName: &GenericPlugin{},
		},
		Cmd:    exec.Command(pluginPath),
		Logger: hclogadapter.New(pluginName, logger),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC,
		},
	})

	// Connect via RPC.
	rpcClient, err := client.Client()
	if err != nil {
		return nil, client, fmt.Errorf("server: could not get gRPC client: %v", err.Error())
	}

	// Request the plugin.
	raw, err := rpcClient.Dispense(pluginName)
	if err != nil {
		return nil, client, fmt.Errorf("server: could not get connect to plugin (%v): %v", pluginName, err.Error())
	}

	p := raw.(ambient.Plugin)
	// if !ok {
	// 	fmt.Println("The plugin is not the right format.")
	// 	return
	// }

	return p, client, nil
}
