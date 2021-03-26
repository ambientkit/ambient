package core

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/josephspurrier/ambient/app/lib/datastorage"
	"github.com/josephspurrier/ambient/app/lib/router"
	"github.com/josephspurrier/ambient/app/lib/websession"
	"github.com/josephspurrier/ambient/app/model"
)

// RegisterPlugins into storage.
func RegisterPlugins(arr []IPlugin, storage *datastorage.Storage) (*PluginSystem, error) {
	// Create the plugin system.
	pluginsys := NewPluginSystem()

	// Load the plugins.
	needSave := false
	ps := storage.Site.PluginSettings
	for _, v := range arr {
		name := v.PluginName()
		_, found := ps[name]
		if !found {
			ps[name] = model.PluginSettings{
				Enabled: false,
			}
			needSave = true
		}

		// Add to the system.
		pluginsys.Plugins[name] = v
	}

	if needSave {
		// Save the plugins.
		storage.Site.PluginSettings = ps
		err := storage.Save()
		if err != nil {
			return nil, err
		}
	}

	return pluginsys, nil
}

// LoadAllPluginPages loads all of the pages from the plugins.
func (c *App) LoadAllPluginPages() error {
	// Set up the plugin routes.
	shouldSave := false
	for name := range c.Storage.Site.PluginSettings {
		bl := c.LoadSinglePluginPages(name)
		if bl {
			shouldSave = true
		}
	}

	if shouldSave {
		// Save the plugin state if something changed.
		err := c.Storage.Save()
		if err != nil {
			return err
		}
	}

	return nil
}

// stringArrayEqual tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func stringArrayEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// LoadSinglePlugin -
// FIXME: Need to add security to this so not any plugin can call it.
func (c *App) LoadSinglePlugin(name string) error {
	save := c.LoadSinglePluginPages(name)
	if save {
		err := c.Storage.Save()
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadSinglePluginPages -
func (c *App) LoadSinglePluginPages(name string) bool {
	shouldSave := false

	// Return if the plug isn't found.
	plugin, ok := c.Storage.Site.PluginSettings[name]
	if !ok {
		return shouldSave
	}

	// Determine if the plugin that is in stored is found in the system.
	v, found := c.Plugins.Plugins[name]

	// If the found setting is different, then update it for saving.
	if found != plugin.Found {
		shouldSave = true
		plugin.Found = found
		c.Storage.Site.PluginSettings[name] = plugin
	}

	// If the fields are different, then update it for saving.
	if !stringArrayEqual(plugin.Fields, v.Fields()) {
		shouldSave = true
		plugin.Fields = v.Fields()
		c.Storage.Site.PluginSettings[name] = plugin
	}

	// If the plugin is not found or not enabled, then skip over it.
	if !found || !plugin.Enabled {
		return shouldSave
	}

	// FIXME: Need to allows users to grant permissions.
	grants := make(map[string]bool)
	grants["site.title:read"] = true
	grants["site.plugins:read"] = true
	grants["site.plugins:enable"] = true
	grants["site.plugins:disable"] = true
	grants["site.plugins:deleteone"] = true
	grants["router:clear"] = true
	grants["plugin:getfield"] = true
	grants["plugin:setfield"] = true
	grants["plugin:setneighborfield"] = true
	grants["plugin:getneighborfield"] = true

	recorder := router.NewRecorder(c.Router)

	toolkit := &Toolkit{
		Router:       recorder,
		Render:       c.Render,
		Security:     c.Sess,
		Site:         NewSecureSite(name, c.Storage, c.Router, grants),
		PluginLoader: c,
	}

	// Enable the plugin and pass in the toolkit.
	err := v.Enable(toolkit)
	if err != nil {
		log.Printf("problem enabling plugin %v: %v", name, err.Error())
	}

	// Load the routes.
	v.Routes()

	// Load the assets.
	assets, files := v.Assets()
	if files == nil {
		// Save the plugin routes so they can be removed if disabled.
		saveRoutesForPlugin(name, recorder, c.Storage)
		return shouldSave
	}

	// Handle embedded assets.
	err = embeddedAssets(recorder, c.Sess, name, assets, files)
	if err != nil {
		log.Println(err.Error())
	}

	// Save the plugin routes so they can be removed if disabled.
	saveRoutesForPlugin(name, recorder, c.Storage)

	return shouldSave
}

func saveRoutesForPlugin(name string, recorder *router.Recorder, storage *datastorage.Storage) {
	// Save the routes.
	arr := make([]model.Route, 0)
	for _, route := range recorder.Routes() {
		arr = append(arr, model.Route{
			Method: route.Method,
			Path:   route.Path,
		})
	}
	storage.PluginRoutes.Routes[name] = arr
}

func embeddedAssets(mux IRouter, sess *websession.Session, pluginName string, files []Asset, assets *embed.FS) error {
	for _, unsafeFile := range files {
		// Recreate the variable when using closures:
		// https://golang.org/doc/faq#closures_and_goroutines
		file := unsafeFile

		// Skip files that are not embedded.
		if !file.Embedded {
			continue
		}

		fileurl := path.Join("/plugins", pluginName, file.SanitizedPath())

		// TODO: Need to check for missing locations and types.

		exists := fileExists(assets, file.SanitizedPath())
		if !exists {
			return fmt.Errorf("plugin (%v) has missing file, please check 'SetAssets()': %v", pluginName, file)
		}

		mux.Get(fileurl, func(w http.ResponseWriter, r *http.Request) (statusCode int, err error) {
			// Don't allow directory browsing.
			if strings.HasSuffix(r.URL.Path, "/") {
				return http.StatusNotFound, nil
			}

			// Handle authentication on resources without changing resources.
			_, loggedIn := sess.User(r)
			if !authAssetAllowed(loggedIn, file) {
				return http.StatusNotFound, nil
			}

			// Use the root directory.
			fsys, err := fs.Sub(assets, ".")
			if err != nil {
				return http.StatusInternalServerError, err
			}

			// Get the requested file name.
			fname := strings.TrimPrefix(r.URL.Path, path.Join("/plugins", pluginName)+"/")

			// Open the file.
			f, err := fsys.Open(fname)
			if err != nil {
				return http.StatusNotFound, nil
			}
			defer f.Close()

			// Get the file time.
			st, err := f.Stat()
			if err != nil {
				return http.StatusInternalServerError, err
			}

			// Get the contents.
			ff, err := ioutil.ReadAll(f)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			// Loop over the items to replace.
			for _, rep := range file.Replace {
				ff = bytes.ReplaceAll(ff, []byte(rep.Find), []byte(rep.Replace))
			}

			http.ServeContent(w, r, fname, st.ModTime(), bytes.NewReader(ff))
			return
		})
	}

	return nil
}

// fileExists determines if an embedded file exists.
func fileExists(assets *embed.FS, filename string) bool {
	// Use the root directory.
	fsys, err := fs.Sub(assets, ".")
	if err != nil {
		return false
	}

	// Open the file.
	f, err := fsys.Open(filename)
	if err != nil {
		return false
	}
	defer f.Close()

	return true
}
