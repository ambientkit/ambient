// Package ambient is a modular web application framework.
package ambient

import (
	"embed"
	"html/template"
	"net/http"
)

// PluginLoader contains the plugins for the Ambient application.
type PluginLoader struct {
	Router         RouterPlugin
	TemplateEngine TemplateEnginePlugin
	TrustedPlugins map[string]bool
	Plugins        []Plugin
	Middleware     []MiddlewarePlugin
}

// PluginCore represents the core of any plugin.
type PluginCore interface {
	// PluginName should be globally unique. Only lowercase letters, numbers,
	// and hypens are permitted. Must start with with a letter.
	PluginName() string // required, read frequently
	// PluginVersion must follow https://semver.org/.
	PluginVersion() string // required, read frequently
}

// Plugin represents a plugin.
type Plugin interface {
	PluginCore

	// These should all have access to the toolkit.
	Enable(toolkit *Toolkit) error                   // optional, called during enable
	Disable() error                                  // optional, called during disable
	Routes()                                         // optional, called during enable
	Assets() ([]Asset, *embed.FS)                    // optional, called during enable
	Settings() []Setting                             // optional, called during special operations
	GrantRequests() []GrantRequest                   // optional, called during every plugin operation against data provider
	FuncMap() func(r *http.Request) template.FuncMap // optional, called on every render

	// Session manager should have middleware with it.
	SessionManager(logger Logger, sessionStorer SessionStorer) (AppSession, error) // optional
}

// LoggingPlugin represents a logging plugin.
type LoggingPlugin interface {
	PluginCore

	Logger(appName string, appVersion string) (AppLogger, error)
}

// StoragePlugin represents a storage plugin.
type StoragePlugin interface {
	PluginCore

	Storage(logger Logger) (DataStorer, SessionStorer, error)
}

// RouterPlugin represents a router engine plugin.
type RouterPlugin interface {
	PluginCore

	Router(logger Logger, render Renderer) (AppRouter, error)
}

// TemplateEnginePlugin represents a template engine plugin.
type TemplateEnginePlugin interface {
	PluginCore

	TemplateEngine(logger Logger, injector AssetInjector) (Renderer, error)
}

// MiddlewarePlugin represents a middleware plugin.
type MiddlewarePlugin interface {
	Plugin

	Middleware() []func(next http.Handler) http.Handler // optional, called during enable
}
