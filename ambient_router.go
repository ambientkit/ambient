package ambient

import (
	"net/http"
)

// AppRouter represents a router.
type AppRouter interface {
	Router

	ServeHTTP(w http.ResponseWriter, r *http.Request)
	SetNotFound(notFound http.Handler)
	SetServeHTTP(h func(w http.ResponseWriter, r *http.Request, status int, err error))
}

// Router represents a router.
type Router interface {
	Handle(method string, path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Get(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Post(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Patch(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Put(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Delete(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Head(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Options(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Error(status int, w http.ResponseWriter, r *http.Request)
	Param(r *http.Request, name string) string
	Wrap(handler http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) (status int, err error)
}

// Route is a route for a router.
type Route struct {
	Method string
	Path   string
}

// CustomServeHTTP allows customization of error handling by the router.
type CustomServeHTTP func(log Logger, renderer Renderer,
	w http.ResponseWriter, r *http.Request, status int, err error)

// SetupRouter sets the router with the NotFound handler and the default handler.
func SetupRouter(logger Logger, mux AppRouter, te Renderer, customServeHTTP CustomServeHTTP) {
	// Set the default handler.
	defaultServeHTTP := func(w http.ResponseWriter, r *http.Request, status int, err error) {
		w.WriteHeader(status)
	}

	// Use the custom handler if it's set.
	serveHTTP := defaultServeHTTP
	if customServeHTTP != nil {
		serveHTTP = func(w http.ResponseWriter, r *http.Request, status int, err error) {
			customServeHTTP(logger, te, w, r, status, err)
		}
	}

	// Send all 404 to the handler.
	notFound := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.Error(http.StatusNotFound, w, r)
	})

	// Set up the router.
	mux.SetServeHTTP(serveHTTP)
	mux.SetNotFound(notFound)
}
