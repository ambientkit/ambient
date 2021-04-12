// Package routerrecorder provides recording functionality.
package routerrecorder

import (
	"net/http"
)

// IRouter represents a router.
type IRouter interface {
	Get(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Post(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Patch(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Put(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Head(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Options(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Delete(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Error(status int, w http.ResponseWriter, r *http.Request)
	Param(r *http.Request, param string) string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// Recorder -
type Recorder struct {
	mux IRouter

	routes []Route
}

// Route -
type Route struct {
	Method string
	Path   string
}

// NewRecorder is a route recorder for plugins.
func NewRecorder(mux IRouter) *Recorder {
	return &Recorder{
		mux: mux,
	}
}

// Routes -
func (rec *Recorder) Routes() []Route {
	return rec.routes
}

func (rec *Recorder) handleRoute(path string, fn func(http.ResponseWriter, *http.Request) (status int, err error), method string, callable func(path string, fn func(http.ResponseWriter, *http.Request) (int, error))) {
	rec.routes = append(rec.routes, Route{
		Method: method,
		Path:   path,
	})
	callable(path, fn)
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
