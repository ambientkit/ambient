package prism

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
			Name:       "prism",
			Version:    "1.0.0",
			AppVersion: "1.0.0",
		},
	}
}

// EmbeddedAssets -
func (p Plugin) EmbeddedAssets() ([]string, *embed.FS) {
	return []string{
		"css/prism-vsc-dark-plus.css",
		"css/clean.css",
	}, &assets
}

// SetPages -
func (p Plugin) Header() string {
	return `<link rel="stylesheet" href="/plugins/prism/css/prism-vsc-dark-plus.css?` + p.Version + `">
	<link rel="stylesheet" href="/plugins/prism/css/clean.css?` + p.Version + `">`
}

// Body -
func (p Plugin) Body() string {
	return `<script src="https://unpkg.com/prismjs@1.23.0/components/prism-core.min.js"></script>
	<script src="https://unpkg.com/prismjs@1.23.0/plugins/autoloader/prism-autoloader.min.js"></script>`
}

// // SetSettings -
// func (pm Prism) SetSettings(s ambsystem.ISettings) error {
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
