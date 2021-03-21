package route

import (
	"fmt"
	"net/http"

	"github.com/josephspurrier/ambient/app/lib/ambsystem"
	"github.com/josephspurrier/ambient/app/lib/datastorage"
	"github.com/josephspurrier/ambient/app/lib/router"
	"github.com/josephspurrier/ambient/plugin/stackedit"
)

// PluginPage -
type PluginPage struct {
	*Core
}

func registerPluginPage(c *PluginPage) {
	c.Router.Get("/dashboard/plugins", c.edit)
	c.Router.Post("/dashboard/plugins", c.update)
}

// edit -
func (c *PluginPage) edit(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "Plugins"
	vars["token"] = c.Sess.SetCSRF(r)
	vars["plugins"] = c.Storage.Site.Plugins
	fmt.Println(c.Storage.Site.Plugins)

	return c.Render.Template(w, r, "dashboard", "plugins_edit", vars)
}

func (c *PluginPage) update(w http.ResponseWriter, r *http.Request) (status int, err error) {
	r.ParseForm()

	// CSRF protection.
	success := c.Sess.CSRF(r)
	if !success {
		return http.StatusBadRequest, nil
	}

	// Loop through each plugin to get the setting.
	for name, plugin := range c.Storage.Site.Plugins {
		plugin.Enabled = (r.FormValue(name) == "on")
		c.Storage.Site.Plugins[name] = plugin
	}

	// Save to storage.
	err = c.Storage.Save()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/dashboard/plugins", http.StatusFound)
	return
}

// LoadPlugins will load the plugins.
func LoadPlugins(mux *router.Mux, storage *datastorage.Storage) {
	// Define the plugins.
	plugins := []ambsystem.IPlugin{
		stackedit.Activate(),
	}

	// Load the plugins.
	needSave := false
	ps := storage.Site.Plugins
	for _, v := range plugins {
		name := v.PluginName()
		_, found := ps[name]
		if !found {
			fmt.Printf("Load new plugin: %v\n", name)
			ps[name] = ambsystem.PluginSettings{
				Enabled: false,
			}
			needSave = true
		} else {
			fmt.Printf("Plugin already found: %v\n", name)
		}

		v.SetPages(mux)
	}

	if needSave {
		// Save the plugins.
		storage.Site.Plugins = ps
		err := storage.Save()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

// func setLogging() {
// 	//temp := os.Stdout
// 	// Turn off logging so it plugins don't have control over output.
// 	os.Stdout = nil
// 	//os.Stdout = temp   // restore it
// }
