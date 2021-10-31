// Package azureblobstorage is an Ambient plugin that provides storage in Azure Blob Storage.
package azureblobstorage

import (
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

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *ambient.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

const (
	// Container allows user to set the Blob container.
	Container = "Container"
)

// Settings returns a list of user settable fields.
func (p *Plugin) Settings() []ambient.Setting {
	return []ambient.Setting{
		{
			Name: Container,
		},
	}
}

// Storage returns data and session storage.
func (p *Plugin) Storage(logger ambient.Logger) (ambient.DataStorer, ambient.SessionStorer, error) {
	container := os.Getenv("AMB_AZURE_CONTAINER")
	if len(container) == 0 {
		var err error
		container, err = p.Site.PluginSettingString(Container)
		if err != nil {
			return nil, nil, err
		}
	}

	ds := store.NewAzureBlobStorage(container, p.sitePath)
	ss := store.NewAzureBlobStorage(container, p.sessionPath)

	return ds, ss, nil
}
