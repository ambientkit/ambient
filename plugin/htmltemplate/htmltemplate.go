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
	*core.PluginMeta
	*core.Toolkit
}

// New returns a new htmltemplate plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "htmltemplate",
			Version:    "1.0.0",
			AppVersion: "1.0.0",
		},
	}
}

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *core.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

// TemplateEngine returns a template engine.
func (p *Plugin) TemplateEngine(logger core.ILogger, pi *core.PluginInjector) (core.IRender, error) {
	tmpl := NewTemplateEngine(logger, pi)
	return tmpl, nil
}
