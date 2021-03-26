// Package plugins provides a plugin management page for an Ambient application.
package plugins

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

// New returns a new plugins plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "plugins",
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
	p.Router.Get("/dashboard/plugins", p.edit)
	p.Router.Post("/dashboard/plugins", p.update)
	p.Router.Get("/dashboard/plugins/:id/delete", p.destroy)
	p.Router.Get("/dashboard/plugins/:id/settings", p.settingsEdit)
	p.Router.Post("/dashboard/plugins/:id/settings", p.settingsUpdate)
}
