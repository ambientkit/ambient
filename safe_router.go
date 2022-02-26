package ambient

import (
	"net/http"
	"os"
	"path"
)

// Recorder -
type Recorder struct {
	log          AppLogger
	pluginsystem *PluginSystem
	mux          AppRouter

	pluginName string
	routeList  []Route
}

// NewRecorder is a route recorder for plugins.
func NewRecorder(pluginName string, log AppLogger, pluginsystem *PluginSystem, mux AppRouter) *Recorder {
	return &Recorder{
		log:          log,
		pluginsystem: pluginsystem,
		mux:          mux,

		pluginName: pluginName,
	}
}

// Routes returns list of routes.
func (rec *Recorder) routes() []Route {
	return rec.routeList
}

func (rec *Recorder) handleRoute(rawpath string, fn func(http.ResponseWriter, *http.Request) (status int, err error),
	method string, callable func(path string, fn func(http.ResponseWriter, *http.Request) (int, error))) {
	// Add the URL prefix to each route.
	path := prefixedRoute(rawpath)
	rec.routeList = append(rec.routeList, Route{
		Method: method,
		Path:   path,
	})
	callable(path, rec.protect(fn))
}

func (rec *Recorder) protect(h func(http.ResponseWriter, *http.Request) (status int, err error)) func(
	http.ResponseWriter, *http.Request) (status int, err error) {
	return func(w http.ResponseWriter, r *http.Request) (status int, err error) {
		if !rec.pluginsystem.Authorized(rec.pluginName, GrantRouterRouteWrite) {
			return http.StatusForbidden, nil
		}

		return h(w, r)
	}
}

func prefixedRoute(urlpath string) string {
	return path.Join(os.Getenv("AMB_URL_PREFIX"), urlpath)
}

// Get -
func (rec *Recorder) Get(path string, fn func(http.ResponseWriter, *http.Request) (status int, err error)) {
	if rec.mux == nil {
		return
	}

	rec.handleRoute(path, fn, http.MethodGet, rec.mux.Get)
}

// Post -
func (rec *Recorder) Post(path string, fn func(http.ResponseWriter, *http.Request) (status int, err error)) {
	if rec.mux == nil {
		return
	}

	rec.handleRoute(path, fn, http.MethodPost, rec.mux.Post)
}

// Patch -
func (rec *Recorder) Patch(path string, fn func(http.ResponseWriter, *http.Request) (status int, err error)) {
	if rec.mux == nil {
		return
	}

	rec.handleRoute(path, fn, http.MethodPatch, rec.mux.Patch)
}

// Put -
func (rec *Recorder) Put(path string, fn func(http.ResponseWriter, *http.Request) (status int, err error)) {
	if rec.mux == nil {
		return
	}

	rec.handleRoute(path, fn, http.MethodPut, rec.mux.Put)
}

// Head -
func (rec *Recorder) Head(path string, fn func(http.ResponseWriter, *http.Request) (status int, err error)) {
	if rec.mux == nil {
		return
	}

	rec.handleRoute(path, fn, http.MethodHead, rec.mux.Head)
}

// Options -
func (rec *Recorder) Options(path string, fn func(http.ResponseWriter, *http.Request) (status int, err error)) {
	if rec.mux == nil {
		return
	}

	rec.handleRoute(path, fn, http.MethodOptions, rec.mux.Options)
}

// Delete -
func (rec *Recorder) Delete(path string, fn func(http.ResponseWriter, *http.Request) (status int, err error)) {
	if rec.mux == nil {
		return
	}

	rec.handleRoute(path, fn, http.MethodDelete, rec.mux.Delete)
}

// Param -
func (rec *Recorder) Param(r *http.Request, name string) string {
	if rec.mux == nil {
		return ""
	}

	return rec.mux.Param(r, name)
}

// Error -
func (rec *Recorder) Error(status int, w http.ResponseWriter, r *http.Request) {
	if rec.mux == nil {
		return
	}

	rec.mux.Error(status, w, r)
}

// Wrap -
func (rec *Recorder) Wrap(handler http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) (status int, err error) {
	return rec.mux.Wrap(handler)
}
