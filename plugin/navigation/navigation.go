// Package navigation provides a navigation page
// for an Ambient application.
package navigation

import (
	"embed"

	"github.com/josephspurrier/ambient/app/core"
)

//go:embed template/*.tmpl
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginMeta
	*core.Toolkit
}

// New returns a new navigation plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "navigation",
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

// Routes gets routes for the plugin.
func (p *Plugin) Routes() {
	p.Mux.Get("/dashboard/plugins/navigation", p.index)
}
