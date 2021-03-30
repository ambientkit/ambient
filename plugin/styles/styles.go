// Package styles provides a styles page
// for an Ambient application.
package styles

import (
	"embed"

	"github.com/josephspurrier/ambient/app/core"
)

//go:embed *
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginMeta
	*core.Toolkit
}

// New returns a new styles plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "styles",
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

const (
	// Favicon allows user to set the favicon.
	Favicon = "Favicon"
	// Styles allows user to set the styles.
	Styles = "Styles"
)

// Fields returns a list of user settable fields.
func (p *Plugin) Fields() []core.Field {
	return []core.Field{
		{
			Name: Favicon,
		},
		{
			Name: Styles,
			Type: core.Textarea,
		},
	}
}

// Routes gets routes for the plugin.
func (p *Plugin) Routes() {
	p.Router.Get("/dashboard/styles", p.edit)
	p.Router.Post("/dashboard/styles", p.update)
}
