// Package securedashboard is an Ambient plugins that prevents unauthenticated access to the /dashboard routes.
package securedashboard

import (
	"net/http"
	"strings"

	"github.com/ambientkit/ambient"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit
}

// New returns an Ambient plugins that prevents unauthenticated access to the /dashboard routes.
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

// DisallowAnon does not allow anonymous users to access the page.
func (p *Plugin) DisallowAnon(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Don't allow anon users to access the dashboard.
		if strings.HasPrefix(r.URL.Path, p.Path("/dashboard")) {
			// If user is not authenticated, don't allow them to access the page.
			_, err := p.Site.AuthenticatedUser(r)
			if err != nil {
				p.Redirect(w, r, "/", http.StatusFound)
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}
