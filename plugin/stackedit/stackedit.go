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

// EmbeddedAssets -
func (p Plugin) EmbeddedAssets() ([]string, *embed.FS) {
	return []string{
		"js/stackedit.js",
	}, &assets
}

// Body -
func (p Plugin) Body() string {
	return `<script src="https://unpkg.com/stackedit-js@1.0.7/docs/lib/stackedit.min.js"></script>
	<script src="/plugins/stackedit/js/stackedit.js?` + p.Version + `"></script>`
}
