package core

import (
	"embed"
	"net/http"
)

// IPlugin represents a plugin.
type IPlugin interface {
	PluginName() string
	PluginVersion() string
	Routes()
	Enable(*Toolkit) error
	Disable() error
	Assets() ([]Asset, *embed.FS)
	Fields() []Field
	Middleware() []func(next http.Handler) http.Handler
	SessionManager(ss SessionStorer, secretKey string) (ISession, error)
	Router(te IRender) (IAppRouter, error)
	//Header() string
	//Body() string
	//SetSettings()
	// Deactivate() error
	// Uninstall() error
}
