// Package redirecttourl redirects to the correct site URL
// for an Ambient application.
package redirecttourl

import (
	"net/http"

	"github.com/josephspurrier/ambient/app/core"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginBase
	*core.Toolkit
}

// New returns a new notrailingslash plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &core.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "redirecttourl"
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

const (
	// SiteScheme allows user to set scheme to redirect to.
	SiteScheme = "Site Scheme"
	// SiteURL allows user to set the URL to redirect to.
	SiteURL = "Site URL"
)

// Settings returns a list of user settable fields.
func (p *Plugin) Settings() []core.Setting {
	return []core.Setting{
		{
			Name: SiteScheme,
			Description: core.SettingDescription{
				Text: "http or https",
			},
		},
		{
			Name: SiteURL,
			Description: core.SettingDescription{
				Text: "example: domain.com",
			},
		},
	}
}

// Middleware returns router middleware.
func (p *Plugin) Middleware() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{
		p.stripSlash,
	}
}
