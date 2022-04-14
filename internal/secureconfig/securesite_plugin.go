package secureconfig

import (
	"bytes"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/internal/pluginsafe"
	"github.com/ambientkit/ambient/pkg/amberror"
	"go.opentelemetry.io/otel/attribute"
)

// Plugins returns the plugin list.
func (ss *SecureSite) Plugins() (map[string]ambient.PluginData, error) {
	if !ss.Authorized(ambient.GrantSitePluginRead) {
		return nil, amberror.ErrAccessDenied
	}

	return ss.pluginsystem.PluginsData(), nil
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
		ss.LoadSinglePluginPages(name)
	}

	return nil
}

func (ss *SecureSite) loadSinglePlugin(name string) error {
	ss.LoadSinglePluginPages(name)

	return nil
}

// LoadSinglePluginPages loads the plugin.
func (ss *SecureSite) LoadSinglePluginPages(name string) {
	// TODO: Should we do name checking here since we have gRPC dynamic plugin loading
	// now? We should use the ambient.Validate package if so.
	// if name == "ambient" {
	// 	ss.log.Error("plugin load: preventing loading plugin with reserved name: %v", name)
	// 	return
	// }

	v, err := ss.pluginsystem.Plugin(name)
	if err != nil {
		ss.log.Error("plugin load: problem loading plugin (%v): %v", name, err.Error())
		return
	}

	recorder := ss.recorder.WithPlugin(name)

	pss, _, err := NewSecureSite(name, ss.log.Named(name), ss.pluginsystem, ss.sess, ss.mux, ss.render, ss.recorder, false)
	if err != nil {
		ss.log.Error("plugin load: problem creating securesite for (%v): %v", name, err.Error())
		return
	}

	toolkit := &ambient.Toolkit{
		Mux:    recorder,
		Render: pluginsafe.NewRenderer(ss.render),
		Site:   pss,
		Log:    pluginsafe.NewPluginLogger(ss.log.Named(name)),
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
		err = ss.embeddedAssets(recorder, ss.sess, name, assets, files)
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

func (ss *SecureSite) embeddedAssets(mux ambient.Router, sess ambient.AppSession, pluginName string, files []ambient.Asset, assets ambient.FileSystemReader) error {
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
// and should never be called again. The middleware should work dynamically.
func (ss *SecureSite) loadAllPluginMiddleware() http.Handler {
	var h http.Handler = ss.mux
	h = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hi := next
			names := ss.pluginsystem.MiddlewareNames()
			// Iterate in reverse since the nature of middleware is recursive and
			// we want top middleware to execute first consistently.
			for i := len(names) - 1; i >= 0; i-- {
				pluginName := names[i]
				pluginRaw, err := ss.pluginsystem.Plugin(pluginName)
				if err != nil {
					continue
				}

				plugin := pluginRaw.(ambient.MiddlewarePlugin)
				hi = ss.loadSinglePluginMiddleware(hi, plugin)
			}
			hi.ServeHTTP(w, r)
		})
	}(h)

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

	// Iterate in reverse since the nature of middleware is recursive and
	// we want top middleware to execute first consistently.
	for i := len(arrHandlers) - 1; i >= 0; i-- {
		pluginMiddleware := arrHandlers[i]
		// Wrap each middleware with a conditional to only use it if the
		// plugin is enabled.
		h = func(next http.Handler) http.Handler {
			// Get plugin name outside of the closure because closures in
			// Go capture variables by reference.
			safePlugin := plugin
			safePluginMiddleware := pluginMiddleware
			middlewareIndex := i
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx, span := ss.log.Trace(r.Context(), "middleware: "+safePlugin.PluginName())
				defer span.End()

				//fmt.Println("middleware called")
				// If the plugin is not found in the storage, then skip it.
				safePluginSettings, err := ss.pluginsystem.PluginData(safePlugin.PluginName())
				if err != nil {
					span.SetAttributes(attribute.Bool("middleware.found", false))
					ss.log.For(ctx).Debug("plugin middleware: plugin (%v) not found", safePlugin.PluginName())
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}

				span.SetAttributes(attribute.Bool("middleware.enabled", safePluginSettings.Enabled))
				span.SetAttributes(attribute.Bool("middleware.authorized", ss.pluginsystem.Authorized(plugin.PluginName(), ambient.GrantRouterMiddlewareWrite)))

				// If the plugin is enabled, then wrap with the middleware.
				if safePluginSettings.Enabled {
					if !ss.pluginsystem.Authorized(plugin.PluginName(), ambient.GrantRouterMiddlewareWrite) {
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}

					ss.log.For(ctx).Debug("plugin middleware: running (enabled) middleware (%v) by plugin: %v", middlewareIndex, safePlugin.PluginName())
					safePluginMiddleware(next).ServeHTTP(w, r.WithContext(ctx))
				} else {
					ss.log.For(ctx).Debug("plugin middleware: skipping (disabled) middleware (%v) by plugin: %v", middlewareIndex, safePlugin.PluginName())
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			})
		}(h)
	}

	return h
}
