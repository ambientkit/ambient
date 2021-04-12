package core

import (
	"embed"
	"html/template"
	"net/http"
)

// PluginBase represents metadata for a plugin that works with the Ambient system.
type PluginBase struct{}

// Enable -
func (p *PluginBase) Enable(*Toolkit) error {
	return nil
}

// Disable -
func (p *PluginBase) Disable() error {
	return nil
}

// Routes -
func (p *PluginBase) Routes() {}

// Assets -
func (p *PluginBase) Assets() ([]Asset, *embed.FS) {
	return nil, nil
}

// FuncMap -
func (p *PluginBase) FuncMap() func(r *http.Request) template.FuncMap {
	return nil
}

// Settings -
func (p *PluginBase) Settings() []Setting {
	return nil
}

// Middleware -
func (p *PluginBase) Middleware() []func(next http.Handler) http.Handler {
	return nil
}

// SessionManager -
func (p *PluginBase) SessionManager(logger ILogger, ss SessionStorer) (ISession, error) {
	return nil, nil
}

// Router -
func (p *PluginBase) Router(logger ILogger, te IRender) (IAppRouter, error) {
	return nil, nil
}

// Storage -
func (p *PluginBase) Storage(logger ILogger) (DataStorer, SessionStorer, error) {
	return nil, nil, nil
}

// TemplateEngine -
func (p *PluginBase) TemplateEngine(logger ILogger, injector AssetInjector) (IRender, error) {
	return nil, nil
}

// Grants -
func (p *PluginBase) Grants() []Grant {
	return nil
}
