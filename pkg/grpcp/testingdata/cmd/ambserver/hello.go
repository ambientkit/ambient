package main

import (
	"github.com/ambientkit/ambient"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
}

// New returns a new hello plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "hello"
}

// PluginVersion returns the plugin version.
func (p *Plugin) PluginVersion() string {
	return "1.0.0"
}

// GrantRequests returns a list of grants requested by the plugin.
func (p *Plugin) GrantRequests() []ambient.GrantRequest {
	return []ambient.GrantRequest{
		{Grant: ambient.GrantUserAuthenticatedRead, Description: "Show different menus to authenticated vs unauthenticated users."},
		{Grant: ambient.GrantUserAuthenticatedWrite, Description: "Access to login and logout the user."},
		{Grant: ambient.GrantRouterRouteWrite, Description: "Access to add routes for pages."},
	}
}

// Cool -
type Cool struct {
	Message string
}

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *ambient.Toolkit) error {
	err := p.PluginBase.Enable(toolkit)
	if err != nil {
		return err
	}

	c := Cool{
		Message: "interesting",
	}

	p.Log.Debug("this is a debug log")
	p.Log.Info("this is an info log: %#v", c)
	p.Log.Warn("this is a warn log")
	p.Log.Warn("this is an error log")

	return nil

}

// Routes sets routes for the plugin.
func (p *Plugin) Routes() {
	//p.Log.Error("plugin: routes called")
	p.Mux.Get("/", p.index)
	p.Mux.Get("/another", p.another)
	p.Mux.Get("/name/{name}", p.name)
	p.Mux.Get("/nameold/{name}", p.Mux.Wrap(p.nameOld))
	p.Mux.Get("/error", p.errorFunc)
	p.Mux.Get("/created", p.created)
	p.Mux.Get("/headers", p.headers)
	p.Mux.Post("/headers", p.headersPOST)
	p.Mux.Get("/login", p.login)
}
