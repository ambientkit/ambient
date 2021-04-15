// Package ambient is a modular web application framework.
package ambient

import (
	"embed"
	"html/template"
	"net/http"
)

// IPlugin represents a plugin.
type IPlugin interface {
	// PluginName should be globally unique. Only lowercase letters, numbers,
	// and hypens are permitted. Must start with with a letter.
	PluginName() string // required, read frequently
	// PluginVersion must follow https://semver.org/.
	PluginVersion() string // required, read frequently

	// These are called before the plugin is enabled so they only have access to the logger.
	Logger(appName string, appVersion string) (IAppLogger, error)                    // optional
	Storage(logger ILogger) (IDataStorer, SessionStorer, error)                      // optional
	SessionManager(logger ILogger, sessionStorer SessionStorer) (IAppSession, error) // optional
	TemplateEngine(logger ILogger, injector AssetInjector) (IRender, error)          // optional
	Router(logger ILogger, render IRender) (IAppRouter, error)                       // optional

	// These should all have access to the toolkit.
	Enable(toolkit *Toolkit) error                      // optional, called during enable
	Disable() error                                     // optional, called during disable
	Routes()                                            // optional, called during enable
	Assets() ([]Asset, *embed.FS)                       // optional, called during enable
	Middleware() []func(next http.Handler) http.Handler // optional, called during enable
	Settings() []Setting                                // optional, called during special operations
	GrantRequests() []GrantRequest                      // optional, called during every plugin operation against data provider
	FuncMap() func(r *http.Request) template.FuncMap    // optional, called on every render
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
