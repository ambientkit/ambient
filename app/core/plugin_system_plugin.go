package core

import (
	"embed"
	"html/template"
	"net/http"
)

// IPlugin represents a plugin.
type IPlugin interface {
	PluginName() string
	PluginVersion() string
	Routes()
	Enable(*Toolkit) error
	Disable() error
	Assets() ([]Asset, *embed.FS, func(r *http.Request) template.FuncMap)
	Fields() []Field
	Middleware() []func(next http.Handler) http.Handler
	SessionManager(ss SessionStorer) (ISession, error)
	Router(te IRender) (IAppRouter, error)
	Storage() (DataStorer, SessionStorer, error)
	TemplateEngine(tm TemplateManager, pi *PluginInjector, pluginNames []string) (IRender, error)
	//Header() string
	//Body() string
	//SetSettings()
	// Deactivate() error
	// Uninstall() error
}
