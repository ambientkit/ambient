package router

import (
	"net/http"
)

type handler struct {
	HandlerFunc     func(w http.ResponseWriter, r *http.Request) (status int, err error)
	CustomServeHTTP func(w http.ResponseWriter, r *http.Request, status int, err error)
}

// ServeHTTP handles all the errors from the HTTP handlers.
func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := fn.HandlerFunc(w, r)
	fn.CustomServeHTTP(w, r, status, err)
}
