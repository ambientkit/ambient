package hello

import "net/http"

func (p *Plugin) index(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "Plugins"
	return p.Render.PluginDashboard(w, r, assets, "template/hello", vars)
}
