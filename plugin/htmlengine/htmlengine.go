// Package htmlengine provides a HTML template engine for an Ambient application.
package htmlengine

import (
	"embed"

	"github.com/josephspurrier/ambient"
)

//go:embed layout/*.tmpl
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit
}

// New returns a new htmlengine plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "htmlengine"
}

// PluginVersion returns the plugin version.
func (p *Plugin) PluginVersion() string {
	return "1.0.0"
}

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *ambient.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

// TemplateEngine returns a template engine.
func (p *Plugin) TemplateEngine(logger ambient.Logger, injector ambient.AssetInjector) (ambient.Renderer, error) {
	tmpl := NewTemplateEngine(logger, injector)
	return tmpl, nil
}
