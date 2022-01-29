// Package awayrouter is an Ambient plugin for a router using a variation of the way router.
package awayrouter

import (
	"fmt"
	"net/http"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/plugin/router/awayrouter/router"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	serveHTTP CustomServeHTTP
}

// New returns an Ambient plugin for a router using a variation of the way router.
// A custom CustomServeHTTP can be passed in to override how errors are handled.
func New(serveHTTP CustomServeHTTP) *Plugin {
	return &Plugin{
		serveHTTP: serveHTTP,
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "awayrouter"
}

// PluginVersion returns the plugin version.
func (p *Plugin) PluginVersion() string {
	return "1.0.0"
}

// Router returns a router.
func (p *Plugin) Router(logger ambient.Logger, te ambient.Renderer) (ambient.AppRouter, error) {
	// Set up the default router.
	mux := router.New()

	// Set the NotFound and custom ServeHTTP handlers.
	p.setupRouter(logger, mux, te)

	return mux, nil
}

// CustomServeHTTP allows customization of error handling by the router.
type CustomServeHTTP func(log ambient.Logger, renderer ambient.Renderer,
	w http.ResponseWriter, r *http.Request, status int, err error)

// setupRouter returns a router with the NotFound handler and the default
// handler set.
func (p *Plugin) setupRouter(logger ambient.Logger, mux ambient.AppRouter, te ambient.Renderer) {
	// Set the handling of all responses.
	defaultServeHTTP := func(w http.ResponseWriter, r *http.Request, status int, err error) {
		// Handle only errors.
		if status >= 400 {
			errText := http.StatusText(status)

			switch status {
			case 403:
				// Already logged on plugin access denials.
				errText = "A plugin has been denied permission."
			case 404:
				// No need to log.
				errText = "Darn, we cannot find the page."
			case 400:
				errText = "Darn, something went wrong."
				if err != nil {
					logger.Info("awayrouter: error (%v): %v", status, err.Error())
				}
			default:
				if err != nil {
					logger.Info("awayrouter: error (%v): %v", status, err.Error())
				}
			}

			status, err = te.Error(w, r, fmt.Sprintf("<h1>%v</h1>%v", status, errText), status, nil, nil)
			if err != nil {
				if err != nil {
					logger.Info("awayrouter: error in rendering error template (%v): %v", status, err.Error())
				}
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}
		}
	}

	serveHTTP := defaultServeHTTP
	if p.serveHTTP != nil {
		serveHTTP = func(w http.ResponseWriter, r *http.Request, status int, err error) {
			p.serveHTTP(logger, te, w, r, status, err)
		}
	}

	// Send all 404 to the handler.
	notFound := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveHTTP(w, r, http.StatusNotFound, nil)
	})

	// Set up the router.
	mux.SetServeHTTP(serveHTTP)
	mux.SetNotFound(notFound)
}
