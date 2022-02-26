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
