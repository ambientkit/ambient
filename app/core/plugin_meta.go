package core

import (
	"embed"
	"html"
	"html/template"
	"net/http"
)

// PluginMeta represents metadata for a plugin that works with the Ambient
// system.
type PluginMeta struct {
	// Name should be globally unique. Only lowercase letters, numbers,
	// and hypens are permitted. Must start with with a letter.
	Name string `json:"name"`
	// Version must follow https://semver.org/.
	Version string `json:"version"`
	// AppVersion is the first compatible version of Ambient that the
	// plugin works with.
	AppVersion string `json:"appversion"`
	// Permissions is which permissions the plugin requests.
	//Permissions []string `json:"permissions"`
}

// Enable -
func (p *PluginMeta) Enable(*Toolkit) error {
	return nil
}

// Disable -
func (p *PluginMeta) Disable() error {
	return nil
}

// Routes -
func (p *PluginMeta) Routes() {}

// Assets -
func (p *PluginMeta) Assets() ([]Asset, *embed.FS, func(r *http.Request) template.FuncMap) {
	return nil, nil, nil
}

// Fields -
func (p *PluginMeta) Fields() []Field {
	return nil
}

// PluginName -
func (p *PluginMeta) PluginName() string {
	return html.EscapeString(p.Name)
}

// PluginVersion -
func (p *PluginMeta) PluginVersion() string {
	return p.Version
}

// Middleware -
func (p *PluginMeta) Middleware() []func(next http.Handler) http.Handler {
	return nil
}

// SessionManager -
func (p *PluginMeta) SessionManager(logger ILogger, ss SessionStorer) (ISession, error) {
	return nil, nil
}

// Router -
func (p *PluginMeta) Router(logger ILogger, te IRender) (IAppRouter, error) {
	return nil, nil
}

// Storage -
func (p *PluginMeta) Storage(logger ILogger) (DataStorer, SessionStorer, error) {
	return nil, nil, nil
}

// TemplateEngine -
func (p *PluginMeta) TemplateEngine(logger ILogger, pi *PluginInjector) (IRender, error) {
	return nil, nil
}
