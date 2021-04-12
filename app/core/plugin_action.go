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

// LoadAllPluginPages loads all of the pages from the plugins.
func (c *App) LoadAllPluginPages() error {
	// Set up the plugin routes.
	shouldSave := false
	for name := range c.Storage.Site.PluginStorage {
		// Skip plugins that are not enabled.
		if !c.Plugins.Enabled(name) {
			continue
		}

		// Load plugin.
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
	_, ok := c.Storage.Site.PluginStorage[plugin.PluginName()]
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
				safePluginSettings, ok := c.Storage.Site.PluginStorage[safePlugin.PluginName()]
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
	v, err := c.Plugins.Plugin(name)
	if err != nil {
		return err
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
	v, err := c.Plugins.Plugin(name)
	if err != nil {
		c.Log.Error("plugin load: problem loading plugin %v: %v", name, err.Error())
		return false
	}

	recorder := routerrecorder.NewRecorder(c.Router)

	toolkit := &Toolkit{
		Mux:      recorder,
		Render:   c.Render, // FIXME: Should probably remove this and create a new struct so it's more secure. A plugin could use a type conversion.
		Security: c.Sess,
		Site:     NewSecureSite(name, c.Log, c.Storage, c.Sess, c.Router, c.Render, c.Plugins),
		Log:      c.Log,
	}

	// Enable the plugin and pass in the toolkit.
	err = v.Enable(toolkit)
	if err != nil {
		c.Log.Error("plugin load: problem enabling plugin %v: %v", name, err.Error())
		return false
	}

	// Load the routes.
	v.Routes()

	// Load the assets.
	assets, files := v.Assets()
	if files == nil {
		// Save the plugin routes so they can be removed if disabled.
		saveRoutesForPlugin(name, recorder, c.Storage)
		return false
	}

	// Handle embedded assets.
	err = embeddedAssets(recorder, c.Sess, name, assets, files)
	if err != nil {
		c.Log.Error("plugin load: problem loading assets for plugin %v: %v", name, err.Error())
	}

	// Save the plugin routes so they can be removed if disabled.
	saveRoutesForPlugin(name, recorder, c.Storage)

	return false
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
