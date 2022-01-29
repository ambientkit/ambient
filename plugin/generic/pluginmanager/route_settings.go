package pluginmanager

import (
	"fmt"
	"net/http"

	"github.com/ambientkit/ambient"
)

type pluginSetting struct {
	Index       int
	Name        string
	Value       string
	FieldType   ambient.SettingType
	Description ambient.SettingDescription
}

func (p *Plugin) settingsEdit(w http.ResponseWriter, r *http.Request) (status int, err error) {
	pluginName := p.Mux.Param(r, "id")

	vars := make(map[string]interface{})
	vars["title"] = "Edit settings for: " + pluginName
	vars["token"] = p.Site.SetCSRF(r)

	settings, err := p.Site.PluginNeighborSettingsList(pluginName)
	if err != nil {
		return p.Site.Error(err)
	}

	arr := make([]pluginSetting, 0)
	for index, setting := range settings {
		if setting.Hide {
			continue
		}

		curVal, err := p.Site.NeighborPluginSettingString(pluginName, setting.Name)
		if err != nil {
			return p.Site.Error(err)
		}

		arr = append(arr, pluginSetting{
			Index:       index,
			Name:        setting.Name,
			Value:       curVal,
			FieldType:   setting.Type,
			Description: setting.Description,
		})
	}

	vars["settings"] = arr

	return p.Render.Page(w, r, assets, "template/settings_edit", nil, vars)
}

func (p *Plugin) settingsUpdate(w http.ResponseWriter, r *http.Request) (status int, err error) {
	pluginName := p.Mux.Param(r, "id")
	r.ParseForm()

	// CSRF protection.
	ok := p.Site.CSRF(r, r.FormValue("token"))
	if !ok {
		return http.StatusBadRequest, nil
	}

	settings, err := p.Site.PluginNeighborSettingsList(pluginName)
	if err != nil {
		return p.Site.Error(err)
	}

	// Loop through each plugin to get the settings then save.
	for index, setting := range settings {
		if setting.Hide {
			continue
		}

		val := r.FormValue(fmt.Sprintf("field%v", index))
		err := p.Site.SetNeighborPluginSetting(pluginName, setting.Name, val)
		if err != nil {
			return p.Site.Error(err)
		}
	}

	// Disable the plugin.
	err = p.Site.DisablePlugin(pluginName, true)
	if err != nil {
		return p.Site.Error(err)
	}

	// Re-enable the plugin to get any change in routes.
	err = p.Site.EnablePlugin(pluginName, true)
	if err != nil {
		return p.Site.Error(err)
	}

	p.Redirect(w, r, fmt.Sprintf("/dashboard/plugins/%v/settings", pluginName), http.StatusFound)
	return
}
