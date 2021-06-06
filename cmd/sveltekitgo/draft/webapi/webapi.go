// Package webapi provides a JSON API for an Ambient application.
package webapi

import (
	"github.com/josephspurrier/ambient"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit
}

// New returns a new webapi plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "webapi"
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

// Disable the plugin background tasks.
func (p *Plugin) Disable() error {
	return nil
}

// Routes gets routes for the plugin.
func (p *Plugin) Routes() {
	p.Mux.Get("/", p.index)
	p.Mux.Get("/v1/auth/session", p.session)
	p.Mux.Post("/v1/auth/login", p.login)
	p.Mux.Post("/v1/auth/logout", p.logout)
	p.Mux.Post("/v1/auth/register", p.register)

	p.Mux.Get("/v1/user/profile", p.userProfile)
	p.Mux.Post("/v1/user/profile", p.updateUserProfile)

	p.Mux.Get("/v1/note", p.loadNotes)
	p.Mux.Post("/v1/note", p.createNote)
	p.Mux.Put("/v1/note/:noteid", p.updateNote)
	p.Mux.Delete("/v1/note/:noteid", p.deleteNote)
}
