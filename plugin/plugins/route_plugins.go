package plugins

import (
	"net/http"
	"sort"

	"github.com/josephspurrier/ambient/app/core"
)

type pluginWithSettings struct {
	Name string
	core.PluginData
	Settings []core.Setting
	Grants   []core.Grant
}

func (p *Plugin) edit(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "Plugins"
	vars["token"] = p.Security.SetCSRF(r)

	plugins, err := p.Site.Plugins()
	if err != nil {
		return p.Site.Error(err)
	}

	pluginNames, err := p.Site.PluginNames()
	if err != nil {
		return p.Site.Error(err)
	}
	sort.Strings(pluginNames)

	arr := make([]pluginWithSettings, 0)
	for _, pluginName := range pluginNames {
		// Get the list of grants.
		grantList, err := p.Site.NeighborPluginGrantList(pluginName)
		if p.Site.ErrorAccessDenied(err) {
			return p.Site.Error(err)
		}

		// Get the list of settings.
		settingsList, err := p.Site.PluginNeighborSettingsList(pluginName)
		if p.Site.ErrorAccessDenied(err) {
			return p.Site.Error(err)
		}

		arr = append(arr, pluginWithSettings{
			Name:       pluginName,
			PluginData: plugins[pluginName],
			Grants:     grantList,
			Settings:   settingsList,
		})
	}

	vars["plugins"] = arr

	return p.Render.Page(w, r, assets, "template/plugins_edit", nil, vars)
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
			// Call the Disable() function on the plugin.
			err = p.PluginLoader.DisableSinglePlugin(name)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			// Disable the plugin in the system.
			err = p.Site.DisablePlugin(name)
			if err != nil {
				return p.Site.Error(err)
			}

			// Clear the plugin routes.
			err = p.Site.ClearAllRoutesForPlugin(name)
			if err != nil {
				return p.Site.Error(err)
			}
		}
	}

	http.Redirect(w, r, "/dashboard/plugins", http.StatusFound)
	return
}

func (p *Plugin) destroy(w http.ResponseWriter, r *http.Request) (status int, err error) {
	ID := p.Mux.Param(r, "id")

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
