package core

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/josephspurrier/ambient/app/lib/cachecontrol"
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
		bl := c.loadSinglePluginPages(name)
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

// LoadAllPluginMiddleware returns a handler that is wrapped in conditional
// middlware from the plugins.
func (c *App) LoadAllPluginMiddleware(h http.Handler, plugins []IPlugin) http.Handler {
	for _, plugin := range plugins {
		// Skip if the plugin isn't found.
		_, ok := c.Storage.Site.PluginSettings[plugin.PluginName()]
		if !ok {
			c.Log.Debug("plugin middleware: plugin not found: %v \n", plugin.PluginName())
			continue
		}

		// Loop through each piece of middleware.
		arrHandlers := plugin.Middleware()
		if len(arrHandlers) > 0 {
			c.Log.Debug("plugin middleware: loading %v middleware for plugin: %v \n", len(plugin.Middleware()), plugin.PluginName())
		}

		for i, pluginMiddleware := range arrHandlers {
			// Wrap each middleware with a conditional to only use it if the
			// plugin is enabled.
			h = func(next http.Handler) http.Handler {
				// Get plugin name outside of the closure because closures in
				// Go capture variables by reference.
				safePlugin := plugin
				safePluginMiddleware := pluginMiddleware
				middlewareIndex := i
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// If the plugin is not found in the settings, then skip it.
					safePluginSettings, ok := c.Storage.Site.PluginSettings[safePlugin.PluginName()]
					if !ok {
						c.Log.Debug("plugin middleware: plugin %v not found\n", safePlugin.PluginName())
						next.ServeHTTP(w, r)
						return
					}

					// If the plugin is enabled, then wrap with the middleware.
					if safePluginSettings.Enabled {
						c.Log.Debug("plugin middleware: running (enabled) middleware %v by plugin: %v\n", middlewareIndex, safePlugin.PluginName())
						safePluginMiddleware(next).ServeHTTP(w, r)
					} else {
						c.Log.Debug("plugin middleware: skipping (disabled) middleware %v by plugin: %v\n", middlewareIndex, safePlugin.PluginName())
						next.ServeHTTP(w, r)
					}
				})
			}(h)
		}
	}

	return h
}

// fieldArrayEqual tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func fieldArrayEqual(a []model.Field, b []Field) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v.Name != b[i].Name {
			return false
		}
		if string(v.Type) != string(b[i].Type) {
			return false
		}
		if v.Description.Text != b[i].Description.Text {
			return false
		}
		if v.Description.URL != b[i].Description.URL {
			return false
		}
	}
	return true
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

// DisableSinglePlugin will disable a plugin and return an error if one occured.
func (c *App) DisableSinglePlugin(name string) error {
	// Determine if the plugin that is in stored is found in the system.
	v, found := c.Plugins.Plugins[name]
	if !found {
		return nil
	}

	return v.Disable()
}

// LoadSinglePlugin -
// FIXME: Need to add security to this so not any plugin can call it.
func (c *App) LoadSinglePlugin(name string) error {
	save := c.loadSinglePluginPages(name)
	if save {
		err := c.Storage.Save()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *App) loadSinglePluginPages(name string) bool {
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
	// FIXME: This is highly coupled, may just want to save on each load.
	if !fieldArrayEqual(plugin.Fields, v.Fields()) {
		shouldSave = true
		plugin.Fields = FieldList(v.Fields()).ModelFields()
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
	grants["site.url:read"] = true
	grants["site.updated:read"] = true
	grants["site.postsandpages:read"] = true
	grants["site.tags:read"] = true
	grants["site.description:read"] = true
	grants["site.content:read"] = true

	recorder := router.NewRecorder(c.Router)

	toolkit := &Toolkit{
		Router:       recorder,
		Render:       c.Render,
		Security:     c.Sess,
		Site:         NewSecureSite(name, c.Storage, c.Router, grants),
		PluginLoader: c,
		Log:          c.Log,
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

		// Skip files that are external, inline, or generic,
		if !file.Routable() {
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

			// Get the requested file name.
			fname := strings.TrimPrefix(r.URL.Path, path.Join("/plugins", pluginName)+"/")

			// Get the file contents.
			ff, status, err := file.Contents(assets)
			if status != http.StatusOK {
				return status, err
			}

			// Set the etag for cache control.
			handled := cachecontrol.Handle(w, r, ff)
			if handled {
				return
			}

			// Assets all have the same time so it's pointless to use the FS
			// ModTime.
			now := time.Now()

			http.ServeContent(w, r, fname, now, bytes.NewReader(ff))
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
