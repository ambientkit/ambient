package router

import (
	"net/http"
)

// IRouter represents a router.
type IRouter interface {
	Get(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Post(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
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

// Get -
func (rec *Recorder) Get(path string, fn func(http.ResponseWriter, *http.Request) (status int, err error)) {
	rec.routes = append(rec.routes, Route{
		Method: http.MethodGet,
		Path:   path,
	})
	rec.mux.Get(path, fn)
}

// Post -
func (rec *Recorder) Post(path string, fn func(http.ResponseWriter, *http.Request) (status int, err error)) {
	rec.routes = append(rec.routes, Route{
		Method: http.MethodPost,
		Path:   path,
	})
	rec.mux.Post(path, fn)
}

// Param -
func (rec *Recorder) Param(r *http.Request, name string) string {
	return rec.mux.Param(r, name)
}

// Error -
func (rec *Recorder) Error(status int, w http.ResponseWriter, r *http.Request) {
	rec.mux.Error(status, w, r)
}
