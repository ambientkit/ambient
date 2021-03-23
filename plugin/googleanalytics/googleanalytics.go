// Package googleanalytics provides Google Analytics tracking
// for an Ambient application.
package googleanalytics

import (
	"embed"

	"github.com/josephspurrier/ambient/app/core"
)

//go:embed *
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	core.PluginMeta
}

// New returns a new googleanalytics plugin.
func New() Plugin {
	return Plugin{
		PluginMeta: core.PluginMeta{
			Name:       "googleanalytics",
			Version:    "1.0.0",
			AppVersion: "1.0.0",
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p Plugin) Assets() ([]core.Asset, *embed.FS) {
	return []core.Asset{
		{
			Path:     "https://www.googletagmanager.com/gtag/js?id=12345",
			Filetype: core.FiletypeJavaScript,
			Location: core.LocationBody,
			Embedded: false,
		},
		{
			Path:     "js/ga.js",
			Filetype: core.FiletypeJavaScript,
			Location: core.LocationBody,
			Embedded: true,
		},
	}, &assets
}
