package app

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/plugin/logrequest"
)

// Proxy -
type Proxy struct {
	handlerUI  http.Handler
	handlerAPI http.Handler
	lrp        *logrequest.Plugin
}

// NewProxy -
func NewProxy(app *ambient.App, handlerUI, handlerAPI http.Handler) *Proxy {
	lrp := logrequest.New()
	lrp.Enable(app.Toolkit(lrp.PluginName()))

	return &Proxy{
		handlerUI:  handlerUI,
		handlerAPI: handlerAPI,
		lrp:        lrp,
	}
}

// ServeHTTP handles the requests.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// If the path starts with /api, serve the API.
	if strings.HasPrefix(r.URL.Path, "/api") {
		p.handlerAPI.ServeHTTP(w, r)
		return
	}

	p.lrp.LogRequest(p.handlerUI).ServeHTTP(w, r)
}

// LoadProxy returns a proxy for the UI and API.
func LoadProxy(log ambient.AppLogger, app *ambient.App, muxAPI http.Handler) *Proxy {
	// Create a proxy to serve the front-end and backend.
	urlUI, err := url.Parse("http://localhost:8080")
	if err != nil {
		log.Fatal("ui target error:", err)
	}

	uiProxy := httputil.NewSingleHostReverseProxy(urlUI)
	uiProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		w.WriteHeader(http.StatusBadGateway)
	}

	return NewProxy(app, uiProxy, muxAPI)
}
