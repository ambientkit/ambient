package core

import (
	"embed"
	"html/template"
	"net/http"
)

// PluginBase represents a base plugin that works with Ambient.
type PluginBase struct{}

// Enable is to enable the plugin. Toolkit should be saved.
func (p *PluginBase) Enable(*Toolkit) error {
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

// FuncMap returns a callable function when passed in a request.
func (p *PluginBase) FuncMap() func(r *http.Request) template.FuncMap {
	return nil
}

// Settings returns a list of user settable fields.
func (p *PluginBase) Settings() []Setting {
	return nil
}

// Middleware returns handler wrapped in middleware.
func (p *PluginBase) Middleware() []func(next http.Handler) http.Handler {
	return nil
}

// SessionManager returns a session manager.
func (p *PluginBase) SessionManager(logger ILogger, ss SessionStorer) (IAppSession, error) {
	return nil, nil
}

// Router returns a request router.
func (p *PluginBase) Router(logger ILogger, te IRender) (IAppRouter, error) {
	return nil, nil
}

// Storage returns data and session storage.
func (p *PluginBase) Storage(logger ILogger) (DataStorer, SessionStorer, error) {
	return nil, nil, nil
}

// TemplateEngine returns a template engine.
func (p *PluginBase) TemplateEngine(logger ILogger, injector AssetInjector) (IRender, error) {
	return nil, nil
}

// Grants returns a list of grants requested by the plugin.
func (p *PluginBase) Grants() []Grant {
	return nil
}

// Logger -
func (p *PluginBase) Logger(appName string, appVersion string) (IAppLogger, error) {
	return nil, nil
}
