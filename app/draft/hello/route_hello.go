package hello

import "net/http"

func (p *Plugin) index(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "Plugins"
	return p.Render.Page(w, r, assets, "template/hello", nil, vars)
}
