package stackedit

import (
	"embed"

	"github.com/josephspurrier/ambient/app/lib/ambsystem"
)

//go:embed *
var assets embed.FS

// Plugin -
type Plugin struct {
	ambsystem.PluginMeta
}

// New sets up the plugin.
func New() Plugin {
	return Plugin{
		PluginMeta: ambsystem.PluginMeta{
			Name:       "stackedit",
			Version:    "1.0.0",
			AppVersion: "1.0.0",
		},
	}
}

// Assets -
func (p Plugin) Assets() ([]ambsystem.Asset, *embed.FS) {
	return []ambsystem.Asset{
		{
			Path:     "https://unpkg.com/stackedit-js@1.0.7/docs/lib/stackedit.min.js",
			Filetype: ambsystem.FiletypeJavaScript,
			Location: ambsystem.LocationBody,
			Embedded: false,
		},
		{
			Path:     "js/stackedit.js",
			Filetype: ambsystem.FiletypeJavaScript,
			Location: ambsystem.LocationBody,
			Embedded: true,
		},
	}, &assets
}

// Body -
func (p Plugin) Body() string {
	return `<script src="https://unpkg.com/stackedit-js@1.0.7/docs/lib/stackedit.min.js"></script>
	<script src="/plugins/stackedit/js/stackedit.js?` + p.Version + `"></script>`
}
