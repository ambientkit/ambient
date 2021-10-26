package simplelogin

import (
	"net/http"
)

func (p *Plugin) home(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "Home"
	return p.Render.Page(w, r, assets, "template/content/home", p.funcMap(r), vars)
}

func (p *Plugin) dashboard(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "Dashboard"
	return p.Render.Page(w, r, assets, "template/content/dashboard", p.funcMap(r), vars)
}
