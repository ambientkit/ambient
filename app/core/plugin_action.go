package core

import (
	"bytes"
	"embed"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/josephspurrier/ambient/app/lib/cachecontrol"
	"github.com/josephspurrier/ambient/app/lib/routerrecorder"
)

// RegisterPlugins into storage.
func RegisterPlugins(arr []IPlugin, storage *Storage) (*PluginSystem, error) {
	// Create the plugin system.
	pluginsys := NewPluginSystem()

	// Load the plugins.
	// Loop through all of the plugins set in the boot.go file.
	needSave := false
	ps := storage.Site.PluginSettings
	for _, v := range arr {
		name := v.PluginName()

		// If there is not an entry in the storage for the plugin, then
		// add a new entry.
		_, found := ps[name]
		if !found {
			ps[name] = PluginSettings{
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
		h = c.LoadSinglePluginMiddleware(h, plugin)
	}

	return h
}

// LoadSinglePluginMiddleware returns a handler that is wrapped in conditional
// middlware from the plugins.
func (c *App) LoadSinglePluginMiddleware(h http.Handler, plugin IPlugin) http.Handler {
	// Skip if the plugin isn't found.
	_, ok := c.Storage.Site.PluginSettings[plugin.PluginName()]
	if !ok {
		c.Log.Debug("plugin middleware: plugin not found: %v\n", plugin.PluginName())
		return h
	}

	// Loop through each piece of middleware.
	arrHandlers := plugin.Middleware()
	if len(arrHandlers) > 0 {
		c.Log.Debug("plugin middleware: loading %v middleware for plugin: %v\n", len(plugin.Middleware()), plugin.PluginName())
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

	return h
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

// InitializePluginStorage -
func InitializePluginStorage(name string, storage *Storage, ps *PluginSystem) (v IPlugin, shouldSave bool, skip bool) {
	// Return if the plug isn't found.
	plugin, ok := storage.Site.PluginSettings[name]
	if !ok {
		return nil, false, true
	}

	// Determine if the plugin that is in stored is found in the system.
	v, found := ps.Plugins[name]

	// If the found setting is different, then update it for saving.
	if found != plugin.Found {
		shouldSave = true
		plugin.Found = found
		storage.Site.PluginSettings[name] = plugin
	}

	// If not found - which means there is data, but the plugin is no longer
	// installed, then save that the plugin is no longer found.
	if !found {
		return nil, true, true
	}

	// If the grants are different, then save the new ones.
	if !grantArrayEqual(v.Grants(), plugin.Grants) {
		shouldSave = true
		plugin.Grants = v.Grants()
		storage.Site.PluginSettings[name] = plugin
	}

	// If the fields are different, then update it for saving.
	// Note: This is highly coupled, need to update this if you add fields.
	if !fieldArrayEqual(plugin, v.Fields()) {
		shouldSave = true
		plugin.Fields = FieldList(v.Fields()).ModelFields()

		// Preserve the order of the fields since maps are not ordered.
		arr := make([]string, 0)
		for _, plug := range v.Fields() {
			arr = append(arr, plug.Name)
		}
		plugin.Order = arr

		storage.Site.PluginSettings[name] = plugin
	}

	// If the plugin is not found or not enabled, then skip over it.
	if !found || !plugin.Enabled {
		return v, shouldSave, true
	}

	return v, shouldSave, false
}

func (c *App) loadSinglePluginPages(name string) bool {
	// Initialize the plugin storage.
	v, shouldSave, skip := InitializePluginStorage(name, c.Storage, c.Plugins)
	if skip {
		return shouldSave
	}

	recorder := routerrecorder.NewRecorder(c.Router)

	toolkit := &Toolkit{
		Mux:          recorder,
		Render:       c.Render, // FIXME: Should probably remove this and create a new struct so it's more secure. A plugin could use a type conversion.
		Security:     c.Sess,
		Site:         NewSecureSite(name, c.Log, c.Storage, c.Sess, c.Router),
		PluginLoader: c,
		Log:          c.Log,
	}

	// Enable the plugin and pass in the toolkit.
	err := v.Enable(toolkit)
	if err != nil {
		c.Log.Error("plugin load: problem enabling plugin %v: %v", name, err.Error())
		return shouldSave
	}

	// Load the routes.
	v.Routes()

	// Load the assets.
	assets, files, _ := v.Assets()
	if files == nil {
		// Save the plugin routes so they can be removed if disabled.
		saveRoutesForPlugin(name, recorder, c.Storage)
		return shouldSave
	}

	// Handle embedded assets.
	err = embeddedAssets(recorder, c.Sess, name, assets, files)
	if err != nil {
		c.Log.Error("plugin load: problem loading assets for plugin %v: %v", name, err.Error())
	}

	// Save the plugin routes so they can be removed if disabled.
	saveRoutesForPlugin(name, recorder, c.Storage)

	return shouldSave
}

func saveRoutesForPlugin(name string, recorder *routerrecorder.Recorder, storage *Storage) {
	// Save the routes.
	arr := make([]Route, 0)
	for _, route := range recorder.Routes() {
		arr = append(arr, Route{
			Method: route.Method,
			Path:   route.Path,
		})
	}
	storage.PluginRoutes.Routes[name] = arr
}

func embeddedAssets(mux IRouter, sess ISession, pluginName string, files []Asset, assets *embed.FS) error {
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
			loggedIn, _ := sess.UserAuthenticated(r)
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
