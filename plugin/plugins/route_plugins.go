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
	Grants   []core.GrantRequest
}

func (p *Plugin) edit(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "Plugins"
	vars["token"] = p.Site.SetCSRF(r)

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
		if err != nil {
			return p.Site.Error(err)
		}

		// Get the list of settings.
		settingsList, err := p.Site.PluginNeighborSettingsList(pluginName)
		if err != nil {
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
	ok := p.Site.CSRF(r)
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
			err = p.Site.EnablePlugin(name, true)
			if err != nil {
				return p.Site.Error(err)
			}
		} else if !enable && info.Enabled {
			// Disable the plugin.
			err = p.Site.DisablePlugin(name, true)
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
