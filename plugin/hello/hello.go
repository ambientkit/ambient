// Package hello provides a hello page for an Ambient application.
package hello

import (
	"embed"
	"net/http"

	"github.com/josephspurrier/ambient/app/core"
)

//go:embed *
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	core.PluginMeta
	*core.Toolkit
}

// New returns a new hello plugin.
func New() Plugin {
	return Plugin{
		PluginMeta: core.PluginMeta{
			Name:       "hello",
			Version:    "1.0.0",
			AppVersion: "1.0.0",
		},
	}
}

// SetPages sets pages on a router.
func (p Plugin) SetPages(toolkit *core.Toolkit) error {
	p.Toolkit = toolkit
	p.Router.Get("/dashboard/hello", p.index)
	return nil
}

func (p *Plugin) index(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "Plugins"
	return p.Render.PluginTemplate(w, r, assets, "template/hello.tmpl", vars)
}
