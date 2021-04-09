package plugins

import (
	"fmt"
	"net/http"

	"github.com/josephspurrier/ambient/app/core"
)

// Field -
type Field struct {
	Index       int
	Name        string
	Value       string
	FieldType   core.FieldType
	Description core.FieldDescription
}

func (p *Plugin) settingsEdit(w http.ResponseWriter, r *http.Request) (status int, err error) {
	pluginName := p.Mux.Param(r, "id")

	vars := make(map[string]interface{})
	vars["title"] = "Edit settings for: " + pluginName
	vars["token"] = p.Security.SetCSRF(r)

	settings, err := p.Site.PluginNeighborSettingsList(pluginName)
	if err != nil {
		return p.Site.Error(err)
	}

	arr := make([]Field, 0)
	for index, setting := range settings {
		curVal, err := p.Site.NeighborPluginSettingString(pluginName, setting.Name)
		if p.Site.ErrorAccessDenied(err) {
			return p.Site.Error(err)
		}

		arr = append(arr, Field{
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
	ok := p.Security.CSRF(r)
	if !ok {
		return http.StatusBadRequest, nil
	}

	settings, err := p.Site.PluginNeighborSettingsList(pluginName)
	if err != nil {
		return p.Site.Error(err)
	}

	// Loop through each plugin to get the settings then save.
	for index, setting := range settings {
		val := r.FormValue(fmt.Sprintf("field%v", index))
		err := p.Site.SetNeighborPluginSetting(pluginName, setting.Name, val)
		if p.Site.ErrorAccessDenied(err) {
			return p.Site.Error(err)
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/dashboard/plugins/%v/settings", pluginName), http.StatusFound)
	return
}
