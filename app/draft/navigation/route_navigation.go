package navigation

import (
	"net/http"
)

func (p *Plugin) index(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	return p.Render.Page(w, r, assets, "template/index", nil, vars)
}
