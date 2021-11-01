// Package azureblobstorage is an Ambient plugin that provides storage in Azure Blob Storage.
package azureblobstorage

import (
	"fmt"
	"os"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/plugin/storage/azureblobstorage/store"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit

	sitePath    string
	sessionPath string
}

// New returns an Ambient plugin that provides storage in Azure Blob Storage.
func New(sitePath string, sessionPath string) *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},

		sitePath:    sitePath,
		sessionPath: sessionPath,
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "azureblobstorage"
}

// PluginVersion returns the plugin version.
func (p *Plugin) PluginVersion() string {
	return "1.0.0"
}

const (
	// ContainerEnv is the Azure Storage container environment variable.
	ContainerEnv = "AMB_AZURE_CONTAINER"
)

// Storage returns data and session storage.
func (p *Plugin) Storage(logger ambient.Logger) (ambient.DataStorer, ambient.SessionStorer, error) {
	// Store the logger so it can be used by the plugin.
	p.Log = logger

	// Load the container from the environment variable.
	container := os.Getenv(ContainerEnv)
	if len(container) == 0 {
		return nil, nil, fmt.Errorf("%v: environment variable, %v, is missing", p.PluginName(), ContainerEnv)
	}

	ds := store.NewAzureBlobStorage(container, p.sitePath)
	ss := store.NewAzureBlobStorage(container, p.sessionPath)

	return ds, ss, nil
}
