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

	plugins, err := p.Site.Plugins()
	if err != nil {
		return p.Site.Error(err)
	}

	settings, ok := plugins[pluginName]
	if !ok {
		settings = core.PluginSettings{}
	}

	arr := make([]Field, 0)
	for index, name := range settings.Fields.SortedNames() {
		v := settings.Fields[name]
		curVal, err := p.Site.NeighborPluginField(pluginName, v.Name)
		if p.Site.ErrorAccessDenied(err) {
			return p.Site.Error(err)
		}

		arr = append(arr, Field{
			Index:       index,
			Name:        v.Name,
			Value:       curVal,
			FieldType:   v.Type,
			Description: v.Description,
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

	// Get list of plugins.
	plugins, err := p.Site.Plugins()
	if err != nil {
		return p.Site.Error(err)
	}

	// Get a list of the settings for the specified plugin.
	settings, ok := plugins[pluginName]
	if !ok {
		settings = core.PluginSettings{}
	}

	// Loop through each plugin to get the settings then save.
	for index, name := range settings.Fields.SortedNames() {
		v := settings.Fields[name]
		val := r.FormValue(fmt.Sprintf("field%v", index))
		err := p.Site.SetNeighborPluginField(pluginName, v.Name, val)
		if p.Site.ErrorAccessDenied(err) {
			return p.Site.Error(err)
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/dashboard/plugins/%v/settings", pluginName), http.StatusFound)
	return
}
