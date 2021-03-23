// Package bearcss provides styles from the Bear Blog (https://bearblog.dev/)
// for an Ambient application.
package bearcss

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

// New returns a new bearcss plugin.
func New() Plugin {
	return Plugin{
		PluginMeta: core.PluginMeta{
			Name:       "bearcss",
			Version:    "1.0.0",
			AppVersion: "1.0.0",
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p Plugin) Assets() ([]core.Asset, *embed.FS) {
	return []core.Asset{
		{
			Path:     "css/bear.css",
			Filetype: core.FiletypeStylesheet,
			Location: core.LocationHeader,
			Embedded: true,
		},
	}, &assets
}
