package mock

import "github.com/ambientkit/ambient"

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase

	pluginName    string
	pluginVersion string

	MockGrants []ambient.GrantRequest
	MockRoutes func(p *ambient.PluginBase)
}

// NewPlugin returns a new mock plugin.
func NewPlugin(name string, version string) *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},

		pluginName:    name,
		pluginVersion: version,
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return p.pluginName
}

// PluginVersion returns the plugin version.
func (p *Plugin) PluginVersion() string {
	return p.pluginVersion
}

// GrantRequests returns a list of grants requested by the plugin.
func (p *Plugin) GrantRequests() []ambient.GrantRequest {
	return p.MockGrants
}

// Routes gets routes for the plugin.
func (p *Plugin) Routes() {
	p.MockRoutes(p.PluginBase)
}
