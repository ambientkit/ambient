// Package stackedit is an Ambient plugin that provides a markdown editor using StackEdit.
package stackedit

import (
	"embed"

	"github.com/ambientkit/ambient"
)

//go:embed js/*.js
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit
}

// New returns an Ambient plugin that provides a markdown editor using StackEdit.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "stackedit"
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

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]ambient.Asset, *embed.FS) {
	return []ambient.Asset{
		{
			Path:     "https://unpkg.com/stackedit-js@1.0.7/docs/lib/stackedit.min.js",
			Filetype: ambient.AssetJavaScript,
			Location: ambient.LocationBody,
			External: true,
			Auth:     ambient.AuthOnly,
		},
		{
			Path:     "js/stackedit.js",
			Filetype: ambient.AssetJavaScript,
			Location: ambient.LocationBody,
			Auth:     ambient.AuthOnly,
		},
	}, &assets
}

// // Body -
// func (p Plugin) Body() string {
// 	return `<script src="https://unpkg.com/stackedit-js@1.0.7/docs/lib/stackedit.min.js"></script>
// 	<script src="/plugins/stackedit/js/stackedit.js?` + p.Version + `"></script>`
// }
