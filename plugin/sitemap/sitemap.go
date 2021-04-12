// Package sitemap provides sitemap functionality
// for an Ambient application.
package sitemap

import (
	"github.com/josephspurrier/ambient/app/core"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginBase
	*core.Toolkit
}

// New returns a new sitemap plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &core.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "sitemap"
}

// PluginVersion returns the plugin version.
func (p *Plugin) PluginVersion() string {
	return "1.0.0"
}

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *core.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

// GrantRequests returns a list of grants requested by the plugin.
func (p *Plugin) GrantRequests() []core.GrantRequest {
	return []core.GrantRequest{
		{Grant: core.GrantSiteURLRead, Description: "Access to read the site URL."},
		{Grant: core.GrantSiteSchemeRead, Description: "Access to read the site scheme."},
		{Grant: core.GrantSiteUpdatedRead, Description: "Access to read the last updated date."},
		{Grant: core.GrantSitePostRead, Description: "Access to read all the posts."},
	}
}

// Routes gets routes for the plugin.
func (p *Plugin) Routes() {
	p.Mux.Get("/sitemap.xml", p.index)
}
