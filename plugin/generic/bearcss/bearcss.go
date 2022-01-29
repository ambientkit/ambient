// Package bearcss is an Ambient plugin that provides styles from the Bear Blog (https://bearblog.dev/).
package bearcss

import (
	"embed"

	"github.com/ambientkit/ambient"
)

//go:embed css/*.css
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit
}

// New returns an Ambient plugin that provides styles from the Bear Blog (https://bearblog.dev/).
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
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
func (p *Plugin) Enable(toolkit *ambient.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

// GrantRequests returns a list of grants requested by the plugin.
func (p *Plugin) GrantRequests() []ambient.GrantRequest {
	return []ambient.GrantRequest{
		{Grant: ambient.GrantSiteAssetWrite, Description: "Access to add a stylesheet to the header."},
		{Grant: ambient.GrantRouterRouteWrite, Description: "Access to create routes for a stylesheet."},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]ambient.Asset, *embed.FS) {
	return []ambient.Asset{
		{
			Filetype: ambient.AssetStylesheet,
			Path:     "css/bear.css",
			Location: ambient.LocationHead,
		},
	}, &assets
}
