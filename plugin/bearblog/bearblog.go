// Package bearblog provides basic blog functionality
// for an Ambient application.
package bearblog

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/josephspurrier/ambient/app/core"
)

//go:embed template/partial/*.tmpl template/content/*.tmpl
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginMeta
	*core.Toolkit
}

// New returns a new bearblog plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "bearblog",
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

const (
	// LoginURL allows user to set the login URL.
	LoginURL = "Login URL"
	// Subtitle allows user to set the Subtitle.
	Subtitle = "Subtitle"
	// Footer allows user to set the footer.
	Footer = "Footer"
)

// Fields returns a list of user settable fields.
func (p *Plugin) Fields() []core.Field {
	return []core.Field{
		{
			Name:    LoginURL,
			Default: "admin", // FIXME: Need to add logic for this.
		},
		{
			Name: Subtitle,
		},
		{
			Name: Footer,
			Type: core.Textarea,
		},
	}
}

// Routes gets routes for the plugin.
func (p *Plugin) Routes() {
	p.Mux.Get("/blog", p.postIndex)
	p.Mux.Get("/:slug", p.postShow)

	p.Mux.Get("/login/:slug", p.login)
	p.Mux.Post("/login/:slug", p.loginPost)
	p.Mux.Get("/dashboard/logout", p.logout)

	p.Mux.Get("/", p.index)
	p.Mux.Get("/dashboard", p.edit)
	p.Mux.Post("/dashboard", p.update)
	p.Mux.Get("/dashboard/reload", p.reload)
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]core.Asset, *embed.FS, func(r *http.Request) template.FuncMap) {
	return []core.Asset{
		{
			Path:     "template/partial/nav.tmpl",
			Filetype: core.AssetGeneric,
			Location: core.LocationHeader,
			Inline:   true,
		},
		{
			Path:     "template/partial/footer.tmpl",
			Filetype: core.AssetGeneric,
			Location: core.LocationFooter,
			Inline:   true,
		},
	}, &assets, p.FuncMap
}
