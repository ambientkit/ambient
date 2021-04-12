// Package plugins provides a plugin management page for an Ambient application.
package plugins

import (
	"embed"

	"github.com/josephspurrier/ambient/app/core"
)

//go:embed template/*.tmpl
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginBase
	*core.Toolkit
}

// New returns a new plugins plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &core.PluginBase{},
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
func (p *Plugin) Enable(toolkit *core.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

// GrantRequests returns a list of grants requested by the plugin.
func (p *Plugin) GrantRequests() []core.GrantRequest {
	return []core.GrantRequest{
		{Grant: core.GrantSitePluginRead, Description: "Access to read the plugins."},
		{Grant: core.GrantSitePluginEnable, Description: "Access to enable plugins."},
		{Grant: core.GrantSitePluginDisable, Description: "Access to disable plugins."},
		{Grant: core.GrantSitePluginDelete, Description: "Access to delete plugin storage."},
		{Grant: core.GrantPluginNeighborSettingRead, Description: "Access to read other plugin settings."},
		{Grant: core.GrantPluginNeighborSettingWrite, Description: "Access to write to other plugin settings"},
		{Grant: core.GrantPluginNeighborGrantRead, Description: "Access to read grant requests for plugins"},
		{Grant: core.GrantPluginNeighborGrantWrite, Description: "Access to approve grants for plugins."},
		{Grant: core.GrantRouterNeighborRouteClear, Description: "Access to clear routes for plugins."},
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
