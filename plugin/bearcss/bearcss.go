// Package bearcss provides styles from the Bear Blog (https://bearblog.dev/)
// for an Ambient application.
package bearcss

import (
	"embed"

	"github.com/josephspurrier/ambient/app/core"
)

//go:embed css/*.css
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginBase
	*core.Toolkit
}

// New returns a new bearcss plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &core.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "bearcss"
}

// PluginVersion returns the plugin version.
func (p *Plugin) PluginVersion() string {
	return "1.0.0"
}

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *core.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]core.Asset, *embed.FS) {
	return []core.Asset{
		{
			Path:     "css/bear.css",
			Filetype: core.AssetStylesheet,
			Location: core.LocationHead,
		},
	}, &assets
}
