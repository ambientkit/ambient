// Package securedashboard prevent unauthenticated access to the /dashboard routes
// for an Ambient application.
package securedashboard

import (
	"net/http"

	"github.com/josephspurrier/ambient"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit
}

// New returns a new securedashboard plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "securedashboard"
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
		{Grant: ambient.GrantUserAuthenticatedRead, Description: "Access to read the plugin settings."},
	}
}

// Middleware returns router middleware.
func (p *Plugin) Middleware() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{
		p.DisallowAnon,
	}
}
