package pluginsafe

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/ambientkit/ambient"
)

// RouteRecorder handles routing for plugins.
type RouteRecorder struct {
	log           ambient.AppLogger
	pluginsystem  ambient.PluginSystem
	mux           ambient.AppRouter
	routeMap      map[string][]PluginFn
	routeMapMutex sync.RWMutex
}

// PluginFn maps a plugin to a function.
type PluginFn struct {
	PluginName string
	Fn         func(http.ResponseWriter, *http.Request) (int, error)
}

// NewRouteRecorder returns a route recorder for use in plugins.
func NewRouteRecorder(log ambient.AppLogger, pluginsystem ambient.PluginSystem, mux ambient.AppRouter) *RouteRecorder {
	return &RouteRecorder{
		log:          log,
		pluginsystem: pluginsystem,
		mux:          mux,
		routeMap:     make(map[string][]PluginFn),
	}
}

// PluginRouteRecorder is a route recorder for a plugin.
type PluginRouteRecorder struct {
	rr         *RouteRecorder
	pluginName string
	routeList  []ambient.Route
}

// WithPlugin sets up recorder for a plugin.
func (rec *RouteRecorder) WithPlugin(pluginName string) *PluginRouteRecorder {
	return &PluginRouteRecorder{
		rr:         rec,
		pluginName: pluginName,
		routeList:  make([]ambient.Route, 0),
	}
}

// Routes returns list of routes.
func (rec *PluginRouteRecorder) Routes() []ambient.Route {
	return rec.routeList
}

func pathKey(method string, path string) string {
	return fmt.Sprintf("%v %v", method, path)
}

func prefixedRoute(urlpath string) string {
	return path.Join(os.Getenv("AMB_URL_PREFIX"), urlpath)
}

func (rec *PluginRouteRecorder) handleRoute(method string, rawpath string, fn func(http.ResponseWriter, *http.Request) (status int, err error)) {
	if rec.rr.mux == nil {
		return
	}

	// Add the URL prefix to each route.
	path := prefixedRoute(rawpath)

	// Store the routes to they can be used later.
	rec.routeList = append(rec.routeList, ambient.Route{
		Method: method,
		Path:   path,
	})

	rs := pathKey(method, path)

	rec.rr.routeMapMutex.Lock()
	_, ok := rec.rr.routeMap[rs]
	if !ok {
		// If the route does not exist, then initialize the map entry.
		rec.rr.routeMap[rs] = make([]PluginFn, 0)
		rec.rr.routeMap[rs] = append(rec.rr.routeMap[rs], PluginFn{
			PluginName: rec.pluginName,
			Fn:         rec.protect(fn),
		})
		rec.rr.routeMapMutex.Unlock()

		rec.rr.mux.Handle(method, path, func(w http.ResponseWriter, r *http.Request) (status int, err error) {
			pathKey := pathKey(method, path)

			// Determine if there are any plugins with routes.
			// This protects against if the route list if modified.
			rec.rr.routeMapMutex.RLock()
			routes, ok := rec.rr.routeMap[pathKey]
			rec.rr.routeMapMutex.RUnlock()
			if !ok {
				return http.StatusNotFound, nil
			}

			for _, plugin := range routes {
				// Skip plugins that aren't enabled.
				if !rec.rr.pluginsystem.Enabled(plugin.PluginName) {
					continue
				}

				// Render the first enable plugin.
				return plugin.Fn(w, r)
			}

			return http.StatusNotFound, nil
		})
		return
	}

	// Ensure the plugin is not already added.
	for _, v := range rec.rr.routeMap[rs] {
		if v.PluginName == rec.pluginName {
			rec.rr.log.Debug("routerecorder: plugin (%v) route already registered: %v", v.PluginName, rs)
			rec.rr.routeMapMutex.Unlock()
			return
		}
	}

	rec.rr.log.Debug("routerecorder: plugin (%v) route added: %v", rec.pluginName, rs)

	// Add the function to the map.
	rec.rr.routeMap[rs] = append(rec.rr.routeMap[rs], PluginFn{
		PluginName: rec.pluginName,
		Fn:         rec.protect(fn),
	})
	rec.rr.routeMapMutex.Unlock()
}

func (rec *PluginRouteRecorder) protect(h func(http.ResponseWriter, *http.Request) (status int, err error)) func(
	http.ResponseWriter, *http.Request) (status int, err error) {
	return func(w http.ResponseWriter, r *http.Request) (status int, err error) {
		if !rec.rr.pluginsystem.Authorized(rec.pluginName, ambient.GrantRouterRouteWrite) {
			return http.StatusForbidden, nil
		}

		return h(w, r)
	}
}

// Get request handler.
func (rec *PluginRouteRecorder) Get(path string, fn func(http.ResponseWriter, *http.Request) (status int, err error)) {
	rec.handleRoute(http.MethodGet, path, fn)
}

// Post request handler.
func (rec *PluginRouteRecorder) Post(path string, fn func(http.ResponseWriter, *http.Request) (status int, err error)) {
	rec.handleRoute(http.MethodPost, path, fn)
}

// Patch request handler.
func (rec *PluginRouteRecorder) Patch(path string, fn func(http.ResponseWriter, *http.Request) (status int, err error)) {
	rec.handleRoute(http.MethodPatch, path, fn)
}

// Put request handler.
func (rec *PluginRouteRecorder) Put(path string, fn func(http.ResponseWriter, *http.Request) (status int, err error)) {
	rec.handleRoute(http.MethodPut, path, fn)
}

// Handle request handler.
func (rec *PluginRouteRecorder) Handle(method string, path string, fn func(http.ResponseWriter, *http.Request) (status int, err error)) {
	rec.handleRoute(method, path, fn)
}

// Head request handler.
func (rec *PluginRouteRecorder) Head(path string, fn func(http.ResponseWriter, *http.Request) (status int, err error)) {
	rec.handleRoute(http.MethodHead, path, fn)
}

// Options request handler.
func (rec *PluginRouteRecorder) Options(path string, fn func(http.ResponseWriter, *http.Request) (status int, err error)) {
	rec.handleRoute(http.MethodOptions, path, fn)
}

// Delete request handler.
func (rec *PluginRouteRecorder) Delete(path string, fn func(http.ResponseWriter, *http.Request) (status int, err error)) {
	rec.handleRoute(http.MethodDelete, path, fn)
}

// Param request handler.
func (rec *PluginRouteRecorder) Param(r *http.Request, name string) string {
	if rec.rr.mux == nil {
		return ""
	}

	return rec.rr.mux.Param(r, name)
}

// Error handler.
func (rec *PluginRouteRecorder) Error(status int, w http.ResponseWriter, r *http.Request) {
	rec.rr.mux.Error(status, w, r)
}

// Wrap for http.HandlerFunc.
func (rec *PluginRouteRecorder) Wrap(handler http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) (status int, err error) {
	return rec.rr.mux.Wrap(handler)
}
