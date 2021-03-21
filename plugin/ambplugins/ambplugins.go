package ambplugins

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/josephspurrier/ambient/app/lib/ambsystem"
)

//go:embed *
var assets embed.FS

// Plugin -
type Plugin struct {
	ambsystem.PluginMeta
	*ambsystem.Toolkit
}

// New sets up the plugin.
func New() Plugin {
	return Plugin{
		PluginMeta: ambsystem.PluginMeta{
			Name:       "ambsystem",
			Version:    "1.0.0",
			AppVersion: "1.0.0",
		},
	}
}

// Assets -
// func (p Plugin) Assets() ([]ambsystem.Asset, *embed.FS) {
// 	return []ambsystem.Asset{
// 		{
// 			Path:     "css/prism-vsc-dark-plus.css",
// 			Filetype: ambsystem.FiletypeStylesheet,
// 			Location: ambsystem.LocationHeader,
// 			Embedded: true,
// 		},
// 		{
// 			Path:     "css/clean.css",
// 			Filetype: ambsystem.FiletypeStylesheet,
// 			Location: ambsystem.LocationHeader,
// 			Embedded: true,
// 		},
// 		{
// 			Path:     "https://unpkg.com/prismjs@1.23.0/components/prism-core.min.js",
// 			Filetype: ambsystem.FiletypeJavaScript,
// 			Location: ambsystem.LocationBody,
// 			Embedded: false,
// 		},
// 		{
// 			Path:     "https://unpkg.com/prismjs@1.23.0/plugins/autoloader/prism-autoloader.min.js",
// 			Filetype: ambsystem.FiletypeJavaScript,
// 			Location: ambsystem.LocationBody,
// 			Embedded: false,
// 		},
// 	}, nil
// }

func (p Plugin) SetPages(toolkit *ambsystem.Toolkit) error {
	p.Toolkit = toolkit

	fmt.Println("pages loaded for ambplugins")

	p.Router.Get("/dashboard/plugins2", p.edit)
	//p.Router.Post("/dashboard/plugins2", p.update)
	//p.Router.Get("/dashboard/plugins2/:id/delete", p.destroy)

	return nil
}

// edit -
func (p *Plugin) edit(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "Plugins"
	vars["token"] = p.Security.SetCSRF(r)
	//vars["plugins"] = p.Storage.Site.Plugins

	//return p.Render.Template(w, r, "dashboard", "plugins_edit", vars)
	return p.Render.PluginTemplate(w, r, assets, "template/plugins_edit.tmpl", vars)
}

// func (p *Plugin) update(w http.ResponseWriter, r *http.Request) (status int, err error) {
// 	r.ParseForm()

// 	// CSRF protection.
// 	success := p.Sess.CSRF(r)
// 	if !success {
// 		return http.StatusBadRequest, nil
// 	}

// 	// Loop through each plugin to get the setting.
// 	for name, plugin := range p.Storage.Site.Plugins {
// 		plugin.Enabled = (r.FormValue(name) == "on")
// 		c.Storage.Site.Plugins[name] = plugin
// 	}

// 	// Save to storage.
// 	err = p.Storage.Save()
// 	if err != nil {
// 		return http.StatusInternalServerError, err
// 	}

// 	http.Redirect(w, r, "/dashboard/plugins", http.StatusFound)
// 	return
// }

// func (p *Plugin) destroy(w http.ResponseWriter, r *http.Request) (status int, err error) {
// 	ID := way.Param(r.Context(), "id")

// 	var ok bool
// 	if _, ok = p.Storage.Site.Plugins[ID]; !ok {
// 		return http.StatusNotFound, nil
// 	}

// 	delete(c.Storage.Site.Plugins, ID)

// 	err = p.Storage.Save()
// 	if err != nil {
// 		return http.StatusInternalServerError, err
// 	}

// 	http.Redirect(w, r, "/dashboard/plugins", http.StatusFound)
// 	return
// }
