// Package prism provides code highlighting through Prism for an Ambient
// application.
package prism

import (
	"embed"

	"github.com/josephspurrier/ambient/app/core"
)

//go:embed *
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginMeta
	*core.Toolkit
}

// New returns a new prism plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "prism",
			Version:    "1.0.0",
			AppVersion: "1.0.0",
		},
	}
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
			Path:     "css/prism-vsc-dark-plus.css",
			Filetype: core.AssetStylesheet,
			Location: core.LocationHead,
			Embedded: true,
		},
		{
			Path:     "css/clean.css",
			Filetype: core.AssetStylesheet,
			Location: core.LocationHead,
			Embedded: true,
		},
		{
			Path:     "https://unpkg.com/prismjs@1.23.0/components/prism-core.min.js",
			Filetype: core.AssetJavaScript,
			Location: core.LocationBody,
			Embedded: false,
		},
		{
			Path:     "https://unpkg.com/prismjs@1.23.0/plugins/autoloader/prism-autoloader.min.js",
			Filetype: core.AssetJavaScript,
			Location: core.LocationBody,
			Embedded: false,
		},
	}, &assets
}

//
// func (p Plugin) Header() string {
// 	return `<link rel="stylesheet" href="/plugins/prism/css/prism-vsc-dark-plus.css?` + p.Version + `">
// 	<link rel="stylesheet" href="/plugins/prism/css/clean.css?` + p.Version + `">`
// }

// // Body -
// func (p Plugin) Body() string {
// 	return `<script src="https://unpkg.com/prismjs@1.23.0/components/prism-core.min.js"></script>
// 	<script src="https://unpkg.com/prismjs@1.23.0/plugins/autoloader/prism-autoloader.min.js"></script>`
// }

// // SetSettings -
// func (pm Prism) SetSettings(s core.ISettings) error {
// 	s.Add("name string", fieldType string, defaultValue string)

// }

// Deactivate deactivates the plugin, but leaves the state in the system.
// func Deactivate() error {
// 	return nil
// }

// // Uninstall removes all plugin state from the system.
// func Uninstall() error {
// 	return nil
// }
