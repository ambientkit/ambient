package router

import "net/http"

// Delete is a shortcut for router.Handle("DELETE", path, handle)
func (m *Mux) Delete(path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	m.router.Handle("DELETE", path, handler{
		handlerFunc:     fn,
		customServeHTTP: m.customServeHTTP,
	})
}

// Get is a shortcut for router.Handle("GET", path, handle)
func (m *Mux) Get(path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	m.router.Handle("GET", path, handler{
		handlerFunc:     fn,
		customServeHTTP: m.customServeHTTP,
	})
}

// Head is a shortcut for router.Handle("HEAD", path, handle)
func (m *Mux) Head(path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	m.router.Handle("HEAD", path, handler{
		handlerFunc:     fn,
		customServeHTTP: m.customServeHTTP,
	})
}

// Options is a shortcut for router.Handle("OPTIONS", path, handle)
func (m *Mux) Options(path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	m.router.Handle("OPTIONS", path, handler{
		handlerFunc:     fn,
		customServeHTTP: m.customServeHTTP,
	})
}

// Patch is a shortcut for router.Handle("PATCH", path, handle)
func (m *Mux) Patch(path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	m.router.Handle("PATCH", path, handler{
		handlerFunc:     fn,
		customServeHTTP: m.customServeHTTP,
	})
}

// Post is a shortcut for router.Handle("POST", path, handle)
func (m *Mux) Post(path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	m.router.Handle("POST", path, handler{
		handlerFunc:     fn,
		customServeHTTP: m.customServeHTTP,
	})
}

// Put is a shortcut for router.Handle("PUT", path, handle)
func (m *Mux) Put(path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	m.router.Handle("PUT", path, handler{
		handlerFunc:     fn,
		customServeHTTP: m.customServeHTTP,
	})
}
