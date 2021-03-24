// Package disqus provides Disqus commenting
// for an Ambient application.
package disqus

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

// New returns a new disqus plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "disqus",
			Version:    "1.0.0",
			AppVersion: "1.0.0",
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]core.Asset, *embed.FS) {
	return []core.Asset{
		{
			Path:     "css/disqus.css",
			Filetype: core.FiletypeStylesheet,
			Location: core.LocationHeader,
			Embedded: true,
		},
		{
			Path:     "js/disqus.js",
			Filetype: core.FiletypeJavaScript,
			Location: core.LocationBody,
			Embedded: true,
		},
	}, &assets
}
