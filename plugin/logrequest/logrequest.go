// Package logrequest provides a request logging
// for an Ambient application.
package logrequest

import (
	"net/http"

	"github.com/josephspurrier/ambient/app/core"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginBase
	*core.Toolkit
}

// New returns a new logrequest plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &core.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "logrequest"
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

// Middleware returns router middleware.
func (p *Plugin) Middleware() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{
		p.LogRequest,
	}
}
