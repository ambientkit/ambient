// Package htmltemplate provides a HTML template engine
// for an Ambient application.
package htmltemplate

import (
	"embed"

	"github.com/josephspurrier/ambient/app/core"
)

//go:embed layout/*.tmpl
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginBase
	*core.Toolkit
}

// New returns a new htmltemplate plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &core.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "htmltemplate"
}

// PluginVersion returns the plugin version.
func (p *Plugin) PluginVersion() string {
	return "1.0.0"
}

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *core.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

// TemplateEngine returns a template engine.
func (p *Plugin) TemplateEngine(logger core.ILogger, injector core.AssetInjector) (core.IRender, error) {
	tmpl := NewTemplateEngine(logger, injector)
	return tmpl, nil
}
