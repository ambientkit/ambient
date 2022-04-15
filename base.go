package ambient

import (
	"context"
	"html/template"
	"net/http"
)

// PluginBase represents a base plugin that works with Ambient.
type PluginBase struct {
	*Toolkit
}

// Enable is to enable the plugin. Toolkit should be saved.
func (p *PluginBase) Enable(_ context.Context, toolkit *Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

// Disable is to disable the plugin.
func (p *PluginBase) Disable(context.Context) error {
	return nil
}

// Routes sets routes for the plugin.
func (p *PluginBase) Routes(context.Context) {}

// Assets returns a list of assets and an embedded filesystem.
func (p *PluginBase) Assets(context.Context) ([]Asset, FileSystemReader) {
	return nil, nil
}

// FuncMap returns a callable function that accepts a request.
func (p *PluginBase) FuncMap(context.Context) func(r *http.Request) template.FuncMap {
	return nil
}

// Settings returns a list of user settable fields.
func (p *PluginBase) Settings(context.Context) []Setting {
	return nil
}

// GrantRequests returns a list of grants requested by the plugin.
func (p *PluginBase) GrantRequests(context.Context) []GrantRequest {
	return nil
}

// Middleware returns a list of middleware.
func (p *PluginBase) Middleware(context.Context) []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{}
}
