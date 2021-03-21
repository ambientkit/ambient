package router

import (
	"net/http"

	"github.com/josephspurrier/ambient/app/lib/away"
)

// Mux contains the router.
type Mux struct {
	router *away.Router

	// customServeHTTP is the serve function.
	customServeHTTP func(w http.ResponseWriter, r *http.Request, status int, err error)
}

// New returns an instance of the router.
func New(csh func(w http.ResponseWriter, r *http.Request, status int, err error), notFound http.Handler) *Mux {
	r := away.NewRouter()
	if notFound != nil {
		r.NotFound = notFound
	}

	return &Mux{
		router:          r,
		customServeHTTP: csh,
	}
}

// Clear will remove a path from the router.
func (m *Mux) Clear(path string) {
	m.router.Remove(path)
}

// ServeHTTP routes the incoming http.Request based on method and path
// extracting path parameters as it goes.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.router.ServeHTTP(w, r)
}

// NotFound shows the 404 page.
func (m *Mux) NotFound(w http.ResponseWriter, r *http.Request) {
	m.customServeHTTP(w, r, http.StatusNotFound, nil)
}

// BadRequest shows the 400 page.
func (m *Mux) BadRequest(w http.ResponseWriter, r *http.Request) {
	m.customServeHTTP(w, r, http.StatusBadRequest, nil)
}

// Param returns a URL parameter.
func (m *Mux) Param(r *http.Request, param string) string {
	return away.Param(r.Context(), param)
}
