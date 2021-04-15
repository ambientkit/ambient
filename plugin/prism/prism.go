// Package prism provides code highlighting through Prism for an Ambient
// application.
package prism

import (
	"embed"

	"github.com/josephspurrier/ambient"
)

//go:embed css/*.css
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit
}

// New returns a new prism plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "prism"
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
		{Grant: ambient.GrantSiteAssetWrite, Description: "Access to add stylesheets and javascript to each page."},
		{Grant: ambient.GrantRouterRouteWrite, Description: "Access to create routes for accessing stylesheets."},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]ambient.Asset, *embed.FS) {
	return []ambient.Asset{
		{
			Path:     "css/prism-vsc-dark-plus.css",
			Filetype: ambient.AssetStylesheet,
			Location: ambient.LocationHead,
		},
		{
			Path:     "css/clean.css",
			Filetype: ambient.AssetStylesheet,
			Location: ambient.LocationHead,
		},
		{
			Path:     "https://unpkg.com/prismjs@1.23.0/components/prism-core.min.js",
			Filetype: ambient.AssetJavaScript,
			Location: ambient.LocationBody,
			External: true,
		},
		{
			Path:     "https://unpkg.com/prismjs@1.23.0/plugins/autoloader/prism-autoloader.min.js",
			Filetype: ambient.AssetJavaScript,
			Location: ambient.LocationBody,
			External: true,
		},
	}, &assets
}
