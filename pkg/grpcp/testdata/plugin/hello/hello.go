// Package hello provides a hello page for an Ambient app.
package hello

import (
	"embed"

	"github.com/ambientkit/ambient/pkg/grpcp"
)

//go:embed template/*.tmpl
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*grpcp.PluginBase
}

// New returns a new hello plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &grpcp.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() (string, error) {
	return "hello", nil
}

// PluginVersion returns the plugin version.
func (p *Plugin) PluginVersion() (string, error) {
	return "1.0.0", nil
}

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *grpcp.Toolkit) error {
	err := p.PluginBase.Enable(toolkit)
	if err != nil {
		return err
	}

	p.Log.Info("plugin: enabled called")

	return nil
}

// Routes sets routes for the plugin.
func (p *Plugin) Routes() error {
	p.Log.Error("plugin: routes called")
	p.Mux.Get("/", p.index)
	p.Mux.Get("/another", p.another)
	p.Mux.Get("/name/{name}", p.name)
	p.Mux.Get("/nameold/{name}", p.Mux.Wrap(p.nameOld))
	p.Mux.Get("/error", p.errorFunc)
	p.Mux.Get("/created", p.created)
	p.Mux.Get("/headers", p.headers)
	p.Mux.Get("/form", p.formGet)
	p.Mux.Post("/form", p.formPOST)
	p.Mux.Get("/login", p.login)
	p.Mux.Get("/loggedin", p.loggedin)
	p.Mux.Get("/errors", p.errorsFunc)
	return nil
}
