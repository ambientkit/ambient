// Package ambient is a pluggable web app framework.
package ambient

import (
	"context"
	"html/template"
	"io"
	"net/http"
)

// PluginCore represents the core of any plugin.
type PluginCore interface {
	// PluginName should be globally unique. It must start with a lowercase
	// letter and then contain only lowercase letters and numbers.
	PluginName(context.Context) string
	// PluginVersion must follow https://semver.org/.
	PluginVersion() string
}

// Plugin represents a plugin.
type Plugin interface {
	PluginCore

	// These should all have access to the toolkit.
	Enable(*Toolkit) error                           // optional, called during enable
	Disable() error                                  // optional, called during disable
	Routes()                                         // optional, called during enable
	Assets() ([]Asset, FileSystemReader)             // optional, called during enable
	Settings() []Setting                             // optional, called during special operations
	GrantRequests() []GrantRequest                   // optional, called during every plugin operation against data provider
	FuncMap() func(r *http.Request) template.FuncMap // optional, called on every render
}

// LoggingPlugin represents a logging plugin.
type LoggingPlugin interface {
	PluginCore

	Logger(appName string, appVersion string, writer io.Writer) (AppLogger, error)
}

// StoragePluginGroup represents a storage plugin and an optional encryption
// package.
type StoragePluginGroup struct {
	Storage    StoragePlugin
	Encryption StorageEncryption
}

// StoragePlugin represents a storage plugin.
type StoragePlugin interface {
	PluginCore

	Storage(logger Logger) (DataStorer, SessionStorer, error)
}

// StorageEncryption represents a encryption/decryption for a storage
// plugin.
type StorageEncryption interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(enc []byte) ([]byte, error)
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

// SessionManagerPlugin represents a session manager plugin.
type SessionManagerPlugin interface {
	PluginCore

	// Session manager should have middleware with it.
	SessionManager(logger Logger, sessionStorer SessionStorer) (AppSession, error)
	Middleware() []func(next http.Handler) http.Handler
}

// MiddlewarePlugin represents a middleware plugin.
type MiddlewarePlugin interface {
	Plugin

	Middleware() []func(next http.Handler) http.Handler // optional, called during enable
}
