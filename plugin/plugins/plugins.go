// Package plugins provides a plugin management page for an Ambient application.
package plugins

import (
	"embed"
	"net/http"

	"github.com/josephspurrier/ambient/app/core"
)

//go:embed *
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginMeta
	*core.Toolkit
}

// New returns a new plugins plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "plugins",
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

// Routes gets routes for the plugin.
func (p *Plugin) Routes() error {
	p.Router.Get("/dashboard/plugins", p.edit)
	p.Router.Post("/dashboard/plugins", p.update)
	p.Router.Get("/dashboard/plugins/:id/delete", p.destroy)
	p.Router.Get("/dashboard/plugins/:id/settings", p.settingsEdit)
	p.Router.Post("/dashboard/plugins/:id/settings", p.settingsUpdate)

	return nil
}

func (p *Plugin) edit(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "Plugins"
	vars["token"] = p.Security.SetCSRF(r)

	plugins, err := p.Site.Plugins()
	if err != nil {
		return p.Site.Error(err)
	}
	vars["plugins"] = plugins

	return p.Render.PluginTemplate(w, r, assets, "template/plugins_edit.tmpl", vars)
}

func (p *Plugin) update(w http.ResponseWriter, r *http.Request) (status int, err error) {
	r.ParseForm()

	// CSRF protection.
	ok := p.Security.CSRF(r)
	if !ok {
		return http.StatusBadRequest, nil
	}

	// Get list of plugins.
	plugins, err := p.Site.Plugins()
	if err != nil {
		return p.Site.Error(err)
	}

	// Loop through each plugin to get the settings then save.
	for name, info := range plugins {
		enable := (r.FormValue(name) == "on")
		if enable && !info.Enabled {
			err = p.Site.EnablePlugin(name)
			if err != nil {
				return p.Site.Error(err)
			}

			// Load the plugin routes.
			err = p.PluginLoader.LoadSinglePlugin(name)
			if err != nil {
				return http.StatusInternalServerError, err
			}
		} else if !enable && info.Enabled {
			err = p.Site.DisablePlugin(name)
			if err != nil {
				return p.Site.Error(err)
			}

			// Clear the plugin routes.
			err = p.Site.ClearRoutePlugin(name)
			if err != nil {
				return p.Site.Error(err)
			}
		}
	}

	http.Redirect(w, r, "/dashboard/plugins", http.StatusFound)
	return
}

func (p *Plugin) destroy(w http.ResponseWriter, r *http.Request) (status int, err error) {
	ID := p.Router.Param(r, "id")

	plugins, err := p.Site.Plugins()
	if err != nil {
		return p.Site.Error(err)
	}

	if _, ok := plugins[ID]; !ok {
		return http.StatusNotFound, nil
	}

	err = p.Site.DeletePlugin(ID)
	if err != nil {
		return p.Site.Error(err)
	}

	http.Redirect(w, r, "/dashboard/plugins", http.StatusFound)
	return
}
