// Package sitemap provides sitemap functionality
// for an Ambient application.
package sitemap

import (
	"github.com/josephspurrier/ambient/app/core"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginMeta
	*core.Toolkit
}

// New returns a new sitemap plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "sitemap",
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
	p.Mux.Get("/sitemap.xml", p.index)
}
