package bearcss

import (
	"embed"

	"github.com/josephspurrier/ambient/app/core"
)

//go:embed *
var assets embed.FS

// Plugin -
type Plugin struct {
	core.PluginMeta
}

// New sets up the plugin.
func New() Plugin {
	return Plugin{
		PluginMeta: core.PluginMeta{
			Name:       "bearcss",
			Version:    "1.0.0",
			AppVersion: "1.0.0",
		},
	}
}

// Assets -
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
