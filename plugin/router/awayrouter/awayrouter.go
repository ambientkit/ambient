// Package awayrouter is an Ambient plugin for a router using a variation of the way router.
package awayrouter

import (
	"fmt"
	"net/http"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/plugin/awayrouter/router"
)

// LoggerHandler -
type LoggerHandler func(log ambient.Logger, w http.ResponseWriter, r *http.Request, status int, err error)

// RouterHandler -
type RouterHandler func(w http.ResponseWriter, r *http.Request, status int, err error)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit

	serveHTTP LoggerHandler
}

// New returns an Ambient plugin for a router using a variation of the way router.
func New(serveHTTP LoggerHandler) *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},

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

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *ambient.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

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
			p.serveHTTP(logger, w, r, status, err)
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
