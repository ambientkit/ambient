package ambient

import (
	"bytes"
	"embed"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/josephspurrier/ambient/lib/cachecontrol"
)

// Plugins returns the plugin list.
func (ss *SecureSite) Plugins() (map[string]PluginData, error) {
	if !ss.Authorized(GrantSitePluginRead) {
		return nil, ErrAccessDenied
	}

	return ss.storage.Site.PluginStorage, nil
}

// PluginNames returns the list of plugin name.
func (ss *SecureSite) PluginNames() ([]string, error) {
	if !ss.Authorized(GrantSitePluginRead) {
		return nil, ErrAccessDenied
	}

	return ss.pluginsystem.names, nil
}

// DeletePlugin deletes a plugin.
func (ss *SecureSite) DeletePlugin(name string) error {
	if !ss.Authorized(GrantSitePluginDelete) {
		return ErrAccessDenied
	}

	delete(ss.storage.Site.PluginStorage, name)

	err := ss.pluginsystem.InitializePlugin(name)
	if err != nil {
		return err
	}

	return ss.storage.Save()
}

// EnablePlugin enables a plugin.
func (ss *SecureSite) EnablePlugin(pluginName string, loadPlugin bool) error {
	if !ss.Authorized(GrantSitePluginEnable) {
		return ErrAccessDenied
	}

	if loadPlugin {
		// Load the plugin and routes.
		err := ss.loadSinglePlugin(pluginName)
		if err != nil {
			return err
		}
	}

	pluginData, ok := ss.storage.Site.PluginStorage[pluginName]
	if !ok {
		return ErrNotFound
	}

	pluginData.Enabled = true
	ss.storage.Site.PluginStorage[pluginName] = pluginData

	return ss.storage.Save()
}

// LoadAllPluginPages loads all of the pages from the plugins.
func (ss *SecureSite) LoadAllPluginPages() error {
	if !ss.Authorized(GrantSitePluginEnable) {
		return ErrAccessDenied
	}

	plugins, err := ss.Plugins()
	if err != nil {
		return err
	}

	for _, name := range ss.pluginsystem.names {
		// Skip plugins that are not enabled.
		if !ss.pluginsystem.Enabled(name) {
			continue
		}

		// Load plugin.
		ss.loadSinglePluginPages(name, plugins)
	}

	return nil
}

func (ss *SecureSite) loadSinglePlugin(name string) error {
	plugins, err := ss.Plugins()
	if err != nil {
		return err
	}

	ss.loadSinglePluginPages(name, plugins)

	return nil
}

func (ss *SecureSite) loadSinglePluginPages(name string, pluginsData map[string]PluginData) {
	if name == "ambient" {
		ss.log.Error("plugin load: preventing loading plugin with reserved name: %v", name)
		return
	}

	v, err := ss.pluginsystem.Plugin(name)
	if err != nil {
		ss.log.Error("plugin load: problem loading plugin %v: %v", name, err.Error())
		return
	}

	recorder := NewRecorder(name, ss.log, ss.storage, ss.mux)

	toolkit := &Toolkit{
		Mux:    recorder,
		Render: NewRenderer(ss.render),
		Site:   NewSecureSite(name, ss.log, ss.storage, ss.pluginsystem, ss.sess, ss.mux, ss.render),
		Log:    NewPluginLogger(ss.log),
	}

	// Enable the plugin and pass in the toolkit.
	err = v.Enable(toolkit)
	if err != nil {
		ss.log.Error("plugin load: problem enabling plugin %v: %v", name, err.Error())
		return
	}

	// Load the routes.
	v.Routes()

	// Load the assets.
	assets, files := v.Assets()
	if files != nil {
		// Handle embedded assets.
		err = embeddedAssets(recorder, ss.sess, name, assets, files)
		if err != nil {
			ss.log.Error("plugin load: problem loading assets for plugin %v: %v", name, err.Error())
		}
	}

	// Save the plugin routes so they can be removed if disabled.
	saveRoutesForPlugin(name, recorder, ss.storage)

}

// DisablePlugin disables a plugin.
func (ss *SecureSite) DisablePlugin(pluginName string, unloadPlugin bool) error {
	if !ss.Authorized(GrantSitePluginDisable) {
		return ErrAccessDenied
	}

	if unloadPlugin {
		// Get the plugin.
		plugin, ok := ss.pluginsystem.plugins[pluginName]
		if !ok {
			return ErrNotFound
		}

		// Disable the plugin.
		err := plugin.Disable()
		if err != nil {
			return err
		}

		// Get the routes for the plugin, not all plugins have routes so don't
		// check if it's ok for not.
		routes := ss.storage.PluginRoutes.Routes[pluginName]

		// Clear each route.
		for _, v := range routes {
			ss.mux.Clear(v.Method, v.Path)
		}
	}

	// Get the plugin data.
	pluginData, ok := ss.storage.Site.PluginStorage[pluginName]
	if !ok {
		return ErrNotFound
	}

	// Disable the plugin.
	pluginData.Enabled = false
	ss.storage.Site.PluginStorage[pluginName] = pluginData

	return ss.storage.Save()
}

func saveRoutesForPlugin(name string, recorder *Recorder, storage *Storage) {
	// Save the routes.
	arr := make([]Route, 0)
	for _, route := range recorder.routes() {
		arr = append(arr, Route{
			Method: route.Method,
			Path:   route.Path,
		})
	}
	storage.PluginRoutes.Routes[name] = arr
}

func embeddedAssets(mux IRouter, sess IAppSession, pluginName string, files []Asset, assets *embed.FS) error {
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

// LoadAllPluginMiddleware returns a handler that is wrapped in conditional
// middlware from the plugins. This only needs to be run once at start up
// and should never be called again.
func (ss *SecureSite) LoadAllPluginMiddleware(h http.Handler) http.Handler {
	for _, pluginName := range ss.pluginsystem.names {
		plugin, ok := ss.pluginsystem.plugins[pluginName]
		if !ok {
			continue
		}

		h = ss.loadSinglePluginMiddleware(h, plugin)
	}

	return h
}

// LoadSinglePluginMiddleware returns a handler that is wrapped in conditional
// middlware from the plugins.
func (ss *SecureSite) loadSinglePluginMiddleware(h http.Handler, plugin IPlugin) http.Handler {
	// Loop through each piece of middleware.
	arrHandlers := plugin.Middleware()
	if len(arrHandlers) > 0 {
		ss.log.Debug("plugin middleware: loading %v middleware for plugin: %v", len(plugin.Middleware()), plugin.PluginName())
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
				safePluginSettings, ok := ss.storage.Site.PluginStorage[safePlugin.PluginName()]
				if !ok {
					ss.log.Debug("plugin middleware: plugin %v not found", safePlugin.PluginName())
					next.ServeHTTP(w, r)
					return
				}

				// If the plugin is enabled, then wrap with the middleware.
				if safePluginSettings.Enabled {
					ss.log.Debug("plugin middleware: running (enabled) middleware %v by plugin: %v", middlewareIndex, safePlugin.PluginName())
					safePluginMiddleware(next).ServeHTTP(w, r)
				} else {
					ss.log.Debug("plugin middleware: skipping (disabled) middleware %v by plugin: %v", middlewareIndex, safePlugin.PluginName())
					next.ServeHTTP(w, r)
				}
			})
		}(h)
	}

	return h
}
