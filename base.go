package ambient

import (
	"embed"
	"html/template"
	"net/http"
)

// PluginBase represents a base plugin that works with Ambient.
type PluginBase struct {
	*Toolkit
}

// Enable is to enable the plugin. Toolkit should be saved.
func (p *PluginBase) Enable(toolkit *Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

// Disable is to disable the plugin.
func (p *PluginBase) Disable() error {
	return nil
}

// Routes sets routes for the plugin.
func (p *PluginBase) Routes() {}

// Assets returns a list of assets and an embedded filesystem.
func (p *PluginBase) Assets() ([]Asset, *embed.FS) {
	return nil, nil
}

// FuncMap returns a callable function that accepts a request.
func (p *PluginBase) FuncMap() func(r *http.Request) template.FuncMap {
	return nil
}

// Settings returns a list of user settable fields.
func (p *PluginBase) Settings() []Setting {
	return nil
}

// GrantRequests returns a list of grants requested by the plugin.
func (p *PluginBase) GrantRequests() []GrantRequest {
	return nil
}
