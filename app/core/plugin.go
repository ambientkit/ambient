package core

import (
	"embed"
	"html/template"
	"net/http"
)

// IPlugin represents a plugin.
type IPlugin interface {
	// PluginName should be globally unique. Only lowercase letters, numbers,
	// and hypens are permitted. Must start with with a letter.
	PluginName() string
	// PluginVersion must follow https://semver.org/.
	PluginVersion() string
	Routes()
	Enable(*Toolkit) error
	Disable() error
	Assets() ([]Asset, *embed.FS)
	FuncMap() func(r *http.Request) template.FuncMap
	Settings() []Setting
	Middleware() []func(next http.Handler) http.Handler

	// These are called before the plugin is enabled so they only have access to the logger.
	SessionManager(logger ILogger, ss SessionStorer) (ISession, error)
	Router(logger ILogger, te IRender) (IAppRouter, error)
	Storage(logger ILogger) (DataStorer, SessionStorer, error)
	TemplateEngine(logger ILogger, injector AssetInjector) (IRender, error)
	Grants() []Grant
}

// IPluginList is a list of IPlugins.
type IPluginList []IPlugin

// PluginNames return an list of plugin names.
func (arr IPluginList) PluginNames() []string {
	pluginNames := make([]string, 0)
	for _, v := range arr {
		pluginNames = append(pluginNames, v.PluginName())
	}

	return pluginNames
}
