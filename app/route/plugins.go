package route

import (
	"fmt"
	"log"
	"net/http"

	"github.com/josephspurrier/ambient/app/lib/ambsystem"
	"github.com/josephspurrier/ambient/app/lib/datastorage"
	"github.com/josephspurrier/ambient/app/lib/router"
	"github.com/josephspurrier/ambient/plugin/prism"
	"github.com/matryer/way"
)

// PluginPage -
type PluginPage struct {
	*Core
}

func registerPluginPage(c *PluginPage) {
	c.Router.Get("/dashboard/plugins", c.edit)
	c.Router.Post("/dashboard/plugins", c.update)
	c.Router.Get("/dashboard/plugins/:id/delete", c.destroy)
}

// edit -
func (c *PluginPage) edit(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "Plugins"
	vars["token"] = c.Sess.SetCSRF(r)
	vars["plugins"] = c.Storage.Site.Plugins

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

func (c *PluginPage) destroy(w http.ResponseWriter, r *http.Request) (status int, err error) {
	ID := way.Param(r.Context(), "id")

	var ok bool
	if _, ok = c.Storage.Site.Plugins[ID]; !ok {
		return http.StatusNotFound, nil
	}

	delete(c.Storage.Site.Plugins, ID)

	err = c.Storage.Save()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/dashboard/plugins", http.StatusFound)
	return
}

// LoadPlugins will load the plugins.
func LoadPlugins(storage *datastorage.Storage) *ambsystem.PluginSystem {
	// Create the plugin system.
	pluginsys := ambsystem.NewPluginSystem()

	// Define the plugins.
	plugins := []ambsystem.IPlugin{
		prism.New(),
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

		// Add the system.
		pluginsys.Plugins[name] = v

		//v.SetPages(mux)
	}

	if needSave {
		// Save the plugins.
		storage.Site.Plugins = ps
		err := storage.Save()
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return pluginsys
}

func LoadPluginPages(storage *datastorage.Storage, mux *router.Mux, plugins *ambsystem.PluginSystem) error {
	// Set up the plugin routes.
	shouldSave := false
	ps := storage.Site.Plugins
	for name, plugin := range ps {
		if !plugin.Enabled {
			continue
		}

		// Determine if the plugin that is in stored is found in the system.
		v, found := plugins.Plugins[name]

		// If the found setting is different, then update it for saving.
		if found != plugin.Found {
			shouldSave = true
			plugin.Found = found
			ps[name] = plugin
		}

		// If the plugin is not found, then skip over trying to read from it.
		if !found {
			continue
		}

		// Load the pages.
		err := v.SetPages(mux)
		if err != nil {
			log.Printf("problem loading pages from plugin %v: %v", name, err.Error())
		}
	}

	if shouldSave {
		// Save the plugin state if something changed.
		storage.Site.Plugins = ps
		err := storage.Save()
		if err != nil {
			return err
		}
	}

	return nil
}

// func setLogging() {
// 	//temp := os.Stdout
// 	// Turn off logging so it plugins don't have control over output.
// 	os.Stdout = nil
// 	//os.Stdout = temp   // restore it
// }
