// Package plugins provides a plugin management page for an Ambient application.
package plugins

import (
	"embed"

	"github.com/josephspurrier/ambient"
)

//go:embed template/*.tmpl
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit
}

// New returns a new plugins plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "plugins"
}

// PluginVersion returns the plugin version.
func (p *Plugin) PluginVersion() string {
	return "1.0.0"
}

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *ambient.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

// GrantRequests returns a list of grants requested by the plugin.
func (p *Plugin) GrantRequests() []ambient.GrantRequest {
	return []ambient.GrantRequest{
		{Grant: ambient.GrantSitePluginRead, Description: "Access to read the plugins."},
		{Grant: ambient.GrantSitePluginEnable, Description: "Access to enable plugins."},
		{Grant: ambient.GrantSitePluginDisable, Description: "Access to disable plugins."},
		{Grant: ambient.GrantSitePluginDelete, Description: "Access to delete plugin storage."},
		{Grant: ambient.GrantPluginNeighborSettingRead, Description: "Access to read other plugin settings."},
		{Grant: ambient.GrantPluginNeighborSettingWrite, Description: "Access to write to other plugin settings"},
		{Grant: ambient.GrantPluginNeighborGrantRead, Description: "Access to read grant requests for plugins"},
		{Grant: ambient.GrantPluginNeighborGrantWrite, Description: "Access to approve grants for plugins."},
		{Grant: ambient.GrantRouterNeighborRouteClear, Description: "Access to clear routes for plugins."},
	}
}

// Routes gets routes for the plugin.
func (p *Plugin) Routes() {
	p.Mux.Get("/dashboard/plugins", p.edit)
	p.Mux.Post("/dashboard/plugins", p.update)
	p.Mux.Get("/dashboard/plugins/:id/delete", p.destroy)
	p.Mux.Get("/dashboard/plugins/:id/settings", p.settingsEdit)
	p.Mux.Post("/dashboard/plugins/:id/settings", p.settingsUpdate)
	p.Mux.Get("/dashboard/plugins/:id/grants", p.grantsEdit)
	p.Mux.Post("/dashboard/plugins/:id/grants", p.grantsUpdate)
}
