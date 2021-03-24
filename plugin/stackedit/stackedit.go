// Package stackedit provides a markdown editor to content blocks for an Ambient
// application.
package stackedit

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

// New returns a new stackedit plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "stackedit",
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
			Path:     "https://unpkg.com/stackedit-js@1.0.7/docs/lib/stackedit.min.js",
			Filetype: core.FiletypeJavaScript,
			Location: core.LocationBody,
			Embedded: false,
			Auth:     core.AuthenticatedOnly,
		},
		{
			Path:     "js/stackedit.js",
			Filetype: core.FiletypeJavaScript,
			Location: core.LocationBody,
			Embedded: true,
			Auth:     core.AuthenticatedOnly,
		},
	}, &assets
}

// // Body -
// func (p Plugin) Body() string {
// 	return `<script src="https://unpkg.com/stackedit-js@1.0.7/docs/lib/stackedit.min.js"></script>
// 	<script src="/plugins/stackedit/js/stackedit.js?` + p.Version + `"></script>`
// }
