package ambient

import (
	"errors"
	"fmt"
	"net/http"
)

// RouteRecorder handles routing for plugins.
type RouteRecorder struct {
	log          AppLogger
	pluginsystem *PluginSystem
	mux          AppRouter
	plugins      []*PluginRouteRecorder

	routeMap map[string][]PluginFn
	// TODO: Add a mutex.
}

// PluginFn maps a plugin to a function.
type PluginFn struct {
	PluginName string
	Fn         func(http.ResponseWriter, *http.Request) (int, error)
}

// NewRouteRecorder returns a route recorder for use in plugins.
func NewRouteRecorder(log AppLogger, pluginsystem *PluginSystem, mux AppRouter) *RouteRecorder {
	return &RouteRecorder{
		log:          log,
		pluginsystem: pluginsystem,
		mux:          mux,
		routeMap:     make(map[string][]PluginFn),
	}
}

// PluginRouteRecorder is a route recorder for a plugin.
type PluginRouteRecorder struct {
	rr *RouteRecorder

	pluginName string
	routeList  []RouteDef
}

// RouteDef is a route for a router.
type RouteDef struct {
	Method string
	Path   string
	Fn     func(http.ResponseWriter, *http.Request) (int, error)
}

// Get request handler.
func (rec *RouteRecorder) withPlugin(pluginName string) *PluginRouteRecorder {
	pr := &PluginRouteRecorder{
		rr:         rec,
		pluginName: pluginName,
	}

	rec.plugins = append(rec.plugins, pr)

	return pr
}

func (rec *PluginRouteRecorder) handleRoute(method string, rawpath string, fn func(http.ResponseWriter, *http.Request) (status int, err error)) {
	if rec.rr.mux == nil {
		return
	}

	// Add the URL prefix to each route.
	path := prefixedRoute(rawpath)
	rec.routeList = append(rec.routeList, RouteDef{
		Method: method,
		Path:   path,
		Fn:     rec.protect(fn),
	})

	rs := fmt.Sprintf("%v %v", method, path)

	_, ok := rec.rr.routeMap[rs]
	if !ok {
		// If the route does not exist, then initialize the map entry.
		rec.rr.routeMap[rs] = make([]PluginFn, 0)
		rec.rr.routeMap[rs] = append(rec.rr.routeMap[rs], PluginFn{
			PluginName: rec.pluginName,
			Fn:         fn,
		})

		rec.rr.mux.Handle(method, path, func(w http.ResponseWriter, r *http.Request) (status int, err error) {
			routes, ok := rec.rr.routeMap[rs]
			if !ok {
				return http.StatusNotFound, nil
			}
			for _, plugin := range routes {
				// Skip plugins that aren't enabled.
				if !rec.rr.pluginsystem.Enabled(plugin.PluginName) {
					continue
				}

				return plugin.Fn(w, r)
			}

			return http.StatusNotFound, nil
		})
		return
	}

	// Add the function to the map.
	rec.rr.routeMap[rs] = append(rec.rr.routeMap[rs], PluginFn{
		PluginName: rec.pluginName,
		Fn:         fn,
	})
}

var errPluginNotEnabled = errors.New("plugin not enabled")

func (rec *PluginRouteRecorder) protect(h func(http.ResponseWriter, *http.Request) (status int, err error)) func(
	http.ResponseWriter, *http.Request) (status int, err error) {
	return func(w http.ResponseWriter, r *http.Request) (status int, err error) {
		if !rec.rr.pluginsystem.Authorized(rec.pluginName, GrantRouterRouteWrite) {
			return http.StatusForbidden, nil
		}

		return h(w, r)
	}
}

////////////////////////////////////////////////////////////////////////////////

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

// Error -
func (rec *PluginRouteRecorder) Error(status int, w http.ResponseWriter, r *http.Request) {
	rec.rr.mux.Error(status, w, r)
}

// Wrap -
func (rec *PluginRouteRecorder) Wrap(handler http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) (status int, err error) {
	return rec.rr.mux.Wrap(handler)
}
