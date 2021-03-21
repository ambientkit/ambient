package route

import (
	"net/http"

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

// func setLogging() {
// 	//temp := os.Stdout
// 	// Turn off logging so it plugins don't have control over output.
// 	os.Stdout = nil
// 	//os.Stdout = temp   // restore it
// }
