package envinfo

import (
	"net/http"
	"os"
	"sort"
)

func (p *Plugin) showEnv(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "Server Environment Variables"

	// Get environment variables in a sort listed.
	arr := os.Environ()
	sort.Strings(arr)
	vars["envs"] = arr

	return p.Render.Page(w, r, assets, "template/show_env", nil, vars)
}
