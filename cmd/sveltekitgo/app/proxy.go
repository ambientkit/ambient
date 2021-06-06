package app

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Proxy -
type Proxy struct {
	handlerUI  http.Handler
	handlerAPI http.Handler
}

// NewProxy -
func NewProxy(handlerUI, handlerAPI http.Handler) *Proxy {
	return &Proxy{
		handlerUI:  handlerUI,
		handlerAPI: handlerAPI,
	}
}

// ServeHTTP handles the requests.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// If the path starts with /api, serve the API.
	if strings.HasPrefix(r.URL.Path, "/api") {
		p.handlerAPI.ServeHTTP(w, r)
		return
	}

	p.handlerUI.ServeHTTP(w, r)
}

// LoadProxy returns a proxy for the UI and API.
func LoadProxy(muxAPI http.Handler) *Proxy {
	// Create a proxy to serve the front-end and backend.
	urlUI, err := url.Parse("http://localhost:8080")
	if err != nil {
		log.Fatal("ui target error:", err)
	}

	uiProxy := httputil.NewSingleHostReverseProxy(urlUI)
	uiProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		w.WriteHeader(http.StatusBadGateway)
	}

	return NewProxy(uiProxy, muxAPI)
}
