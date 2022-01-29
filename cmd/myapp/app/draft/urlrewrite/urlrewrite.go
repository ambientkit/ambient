// Package urlrewrite removes trailing slash from requests for an Ambient app.
package urlrewrite

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/ambientkit/ambient"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit
}

// New returns a new urlrewrite plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "urlrewrite"
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

// Middleware returns router middleware.
func (p *Plugin) Middleware() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{
		p.HandlePrefix,
	}
}

// HandlePrefix will handle URLs behind a proxy.
func (p *Plugin) HandlePrefix(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlprefix := os.Getenv("AMB_URL_PREFIX")

		// If there is a prefix, then strip it out for all requests.
		if len(urlprefix) > 0 {
			r.URL.Path = path.Join("/", strings.TrimPrefix(r.URL.Path, urlprefix))
			p.Log.Debug("Rewrote URL: %v", r.URL.Path)
		}

		next.ServeHTTP(w, r)
	})
}
