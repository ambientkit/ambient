// Package ambient is a modular web application framework.
package ambient

import (
	"embed"
	"html/template"
	"net/http"
)

// PluginLoader contains the plugins for the Ambient application.
type PluginLoader struct {
	Plugins    []Plugin
	Middleware []Middleware
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

	// These are called before the plugin is enabled so they only have access to the logger.
	Logger(appName string, appVersion string) (AppLogger, error)                   // optional
	Storage(logger Logger) (DataStorer, SessionStorer, error)                      // optional
	SessionManager(logger Logger, sessionStorer SessionStorer) (AppSession, error) // optional
	TemplateEngine(logger Logger, injector AssetInjector) (Renderer, error)        // optional
	Router(logger Logger, render Renderer) (AppRouter, error)                      // optional

	// These should all have access to the toolkit.
	Enable(toolkit *Toolkit) error                   // optional, called during enable
	Disable() error                                  // optional, called during disable
	Routes()                                         // optional, called during enable
	Assets() ([]Asset, *embed.FS)                    // optional, called during enable
	Settings() []Setting                             // optional, called during special operations
	GrantRequests() []GrantRequest                   // optional, called during every plugin operation against data provider
	FuncMap() func(r *http.Request) template.FuncMap // optional, called on every render
}

// Middleware represents middleware.
type Middleware interface {
	Plugin

	Middleware() []func(next http.Handler) http.Handler // optional, called during enable
}
