package secureconfig

import (
	"bytes"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/internal/amberror"
	"github.com/ambientkit/ambient/internal/pluginsafe"
)

// Plugins returns the plugin list.
func (ss *SecureSite) Plugins() (map[string]ambient.PluginData, error) {
	if !ss.Authorized(ambient.GrantSitePluginRead) {
		return nil, amberror.ErrAccessDenied
	}

	return ss.pluginsystem.Plugins(), nil
}

// PluginNames returns the list of plugin name.
func (ss *SecureSite) PluginNames() ([]string, error) {
	if !ss.Authorized(ambient.GrantSitePluginRead) {
		return nil, amberror.ErrAccessDenied
	}
	return ss.pluginsystem.Names(), nil
}

// DeletePlugin deletes a plugin.
func (ss *SecureSite) DeletePlugin(name string) error {
	if !ss.Authorized(ambient.GrantSitePluginDelete) {
		return amberror.ErrAccessDenied
	}

	err := ss.pluginsystem.RemovePlugin(name)
	if err != nil {
		return err
	}

	p, err := ss.pluginsystem.Plugin(name)
	if err != nil {
		return err
	}

	return ss.pluginsystem.InitializePlugin(name, p.PluginVersion())
}

// EnablePlugin enables a plugin.
func (ss *SecureSite) EnablePlugin(pluginName string, loadPlugin bool) error {
	if !ss.Authorized(ambient.GrantSitePluginEnable) {
		return amberror.ErrAccessDenied
	}

	if loadPlugin {
		// Load the plugin and routes.
		err := ss.loadSinglePlugin(pluginName)
		if err != nil {
			return err
		}
	}

	return ss.pluginsystem.SetEnabled(pluginName, true)
}

// loadAllPluginPages loads all of the pages from the plugins.
func (ss *SecureSite) loadAllPluginPages() error {
	if !ss.Authorized(ambient.GrantSitePluginEnable) {
		return amberror.ErrAccessDenied
	}

	for _, name := range ss.pluginsystem.Names() {
		// Skip plugins that are not enabled.
		if !ss.pluginsystem.Enabled(name) {
			continue
		}

		// Load plugin.
		ss.loadSinglePluginPages(name)
	}

	return nil
}

func (ss *SecureSite) loadSinglePlugin(name string) error {
	ss.loadSinglePluginPages(name)

	return nil
}

func (ss *SecureSite) loadSinglePluginPages(name string) {
	if name == "ambient" {
		ss.log.Error("plugin load: preventing loading plugin with reserved name: %v", name)
		return
	}

	v, err := ss.pluginsystem.Plugin(name)
	if err != nil {
		ss.log.Error("plugin load: problem loading plugin (%v): %v", name, err.Error())
		return
	}

	recorder := ss.recorder.WithPlugin(name)

	pss, _, _ := NewSecureSite(name, ss.log, ss.pluginsystem, ss.sess, ss.mux, ss.render, ss.recorder, false)

	toolkit := &ambient.Toolkit{
		Mux:    recorder,
		Render: pluginsafe.NewRenderer(ss.render),
		Site:   pss,
		Log:    pluginsafe.NewPluginLogger(ss.log),
	}

	// Enable the plugin and pass in the toolkit.
	err = v.Enable(toolkit)
	if err != nil {
		ss.log.Error("plugin load: problem enabling plugin (%v): %v", name, err.Error())
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
			ss.log.Error("plugin load: problem loading assets for plugin (%v): %v", name, err.Error())
		}
	}

	// Save the plugin routes so they can be removed if disabled.
	SaveRoutesForPlugin(name, recorder, ss.pluginsystem)
}

// DisablePlugin disables a plugin.
func (ss *SecureSite) DisablePlugin(pluginName string, unloadPlugin bool) error {
	if !ss.Authorized(ambient.GrantSitePluginDisable) {
		return amberror.ErrAccessDenied
	}

	if unloadPlugin {
		// Get the plugin.
		plugin, err := ss.pluginsystem.Plugin(pluginName)
		if err != nil {
			return amberror.ErrNotFound
		}

		// Disable the plugin.
		err = plugin.Disable()
		if err != nil {
			return err
		}
	}

	// Disable plugin.
	return ss.pluginsystem.SetEnabled(pluginName, false)
}

// SaveRoutesForPlugin will save the routes in the plugin system.
func SaveRoutesForPlugin(name string, recorder *pluginsafe.PluginRouteRecorder, pluginsystem ambient.PluginSystem) {
	// Save the routes.
	arr := make([]ambient.Route, 0)
	for _, route := range recorder.Routes() {
		arr = append(arr, ambient.Route{
			Method: route.Method,
			Path:   route.Path,
		})
	}
	pluginsystem.SetRoute(name, arr)
}

func embeddedAssets(mux ambient.Router, sess ambient.AppSession, pluginName string, files []ambient.Asset, assets ambient.FileSystemReader) error {
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

		if !unsafeFile.SkipExistCheck {
			exists := fileExists(assets, file.SanitizedPath())
			if !exists {
				return fmt.Errorf("plugin (%v) has missing file, please check 'SetAssets()': %v", pluginName, file)
			}
		}

		mux.Get(fileurl, func(w http.ResponseWriter, r *http.Request) (err error) {
			// Don't allow folder browsing.
			if strings.HasSuffix(r.URL.Path, "/") {
				return mux.StatusError(http.StatusNotFound, nil)
			}

			// Handle authentication on resources without changing resources.
			_, err = sess.AuthenticatedUser(r)
			if !ambient.AuthAssetAllowed(err == nil, file) {
				return mux.StatusError(http.StatusNotFound, nil)
			}

			// Get the requested file name.
			fname := strings.TrimPrefix(r.URL.Path, path.Join("/plugins", pluginName)+"/")

			// Get the file contents.
			ff, status, err := file.Contents(assets)
			if status != http.StatusOK {
				return mux.StatusError(status, err)
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

// loadAllPluginMiddleware returns a handler that is wrapped in conditional
// middleware from the plugins. This only needs to be run once at start up
// and should never be called again.
func (ss *SecureSite) loadAllPluginMiddleware() http.Handler {
	var h http.Handler = ss.mux
	for _, pluginName := range ss.pluginsystem.MiddlewareNames() {
		plugin, err := ss.pluginsystem.Plugin(pluginName)
		if err != nil {
			continue
		}

		h = ss.loadSinglePluginMiddleware(h, plugin.(ambient.MiddlewarePlugin))
	}

	return h
}

// LoadSinglePluginMiddleware returns a handler that is wrapped in conditional
// middleware from the plugins.
func (ss *SecureSite) loadSinglePluginMiddleware(h http.Handler, plugin ambient.MiddlewarePlugin) http.Handler {
	// Loop through each piece of middleware.
	arrHandlers := plugin.Middleware()
	if len(arrHandlers) > 0 {
		ss.log.Debug("plugin middleware: loading (%v) middleware for plugin: %v", len(plugin.Middleware()), plugin.PluginName())
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
				// If the plugin is not found in the storage, then skip it.
				safePluginSettings, err := ss.pluginsystem.PluginData(safePlugin.PluginName())
				if err != nil {
					ss.log.Debug("plugin middleware: plugin (%v) not found", safePlugin.PluginName())
					next.ServeHTTP(w, r)
					return
				}

				// If the plugin is enabled, then wrap with the middleware.
				if safePluginSettings.Enabled {
					if !ss.pluginsystem.Authorized(plugin.PluginName(), ambient.GrantRouterMiddlewareWrite) {
						next.ServeHTTP(w, r)
						return
					}

					ss.log.Debug("plugin middleware: running (enabled) middleware (%v) by plugin: %v", middlewareIndex, safePlugin.PluginName())
					safePluginMiddleware(next).ServeHTTP(w, r)
				} else {
					ss.log.Debug("plugin middleware: skipping (disabled) middleware (%v) by plugin: %v", middlewareIndex, safePlugin.PluginName())
					next.ServeHTTP(w, r)
				}
			})
		}(h)
	}

	return h
}
