// Package securedashboard prevent unauthenticated access to the /dashboard routes
// for an Ambient application.
package securedashboard

import (
	"net/http"

	"github.com/josephspurrier/ambient/app/core"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginMeta
	*core.Toolkit
}

// New returns a new securedashboard plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "securedashboard",
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

// Middleware returns router middleware.
func (p *Plugin) Middleware() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{
		p.DisallowAnon,
	}
}
