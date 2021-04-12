package plugins

import (
	"fmt"
	"net/http"

	"github.com/josephspurrier/ambient/app/core"
)

type pluginGrant struct {
	Index   int
	Name    core.Grant
	Granted bool
	//Description core.SettingDescription
}

func (p *Plugin) grantsEdit(w http.ResponseWriter, r *http.Request) (status int, err error) {
	pluginName := p.Mux.Param(r, "id")

	vars := make(map[string]interface{})
	vars["title"] = "Edit grants for: " + pluginName
	vars["token"] = p.Security.SetCSRF(r)

	grantList, err := p.Site.NeighborPluginGrantList(pluginName)
	if p.Site.ErrorAccessDenied(err) {
		return p.Site.Error(err)
	}

	grants, err := p.Site.NeighborPluginGrants(pluginName)
	if p.Site.ErrorAccessDenied(err) {
		return p.Site.Error(err)
	}

	arr := make([]pluginGrant, 0)
	for index, name := range grantList {
		arr = append(arr, pluginGrant{
			Index:   index,
			Name:    name,
			Granted: grants[name],
		})
	}

	vars["grants"] = arr

	return p.Render.Page(w, r, assets, "template/grants_edit", nil, vars)
}

func (p *Plugin) grantsUpdate(w http.ResponseWriter, r *http.Request) (status int, err error) {
	pluginName := p.Mux.Param(r, "id")
	r.ParseForm()

	// CSRF protection.
	ok := p.Security.CSRF(r)
	if !ok {
		return http.StatusBadRequest, nil
	}

	grantList, err := p.Site.NeighborPluginGrantList(pluginName)
	if p.Site.ErrorAccessDenied(err) {
		return p.Site.Error(err)
	}

	// Loop through each plugin to get the grants then save.
	for index, name := range grantList {
		val := r.FormValue(fmt.Sprintf("field%v", index))
		err := p.Site.SetNeighborPluginGrant(pluginName, name, val == "true")
		if p.Site.ErrorAccessDenied(err) {
			return p.Site.Error(err)
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/dashboard/plugins/%v/grants", pluginName), http.StatusFound)
	return
}
