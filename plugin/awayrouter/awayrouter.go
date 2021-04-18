// Package awayrouter provides a router using a forked version of way
// for an Ambient application.
package awayrouter

import (
	"fmt"
	"net/http"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/plugin/awayrouter/router"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit
}

// New returns a new awayrouter plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
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
func (p *Plugin) Router(logger ambient.Logger, te ambient.IRender) (ambient.AppRouter, error) {
	// Set up the default router.
	mux := router.New()

	// Set the NotFound and custom ServeHTTP handlers.
	setupRouter(logger, mux, te)

	return mux, nil
}

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *ambient.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

// setupRouter returns a router with the NotFound handler and the default
// handler set.
func setupRouter(logger ambient.Logger, mux ambient.AppRouter, te ambient.IRender) {
	// Set the handling of all responses.
	customServeHTTP := func(w http.ResponseWriter, r *http.Request, status int, err error) {
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

	// Send all 404 to the customer handler.
	notFound := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customServeHTTP(w, r, http.StatusNotFound, nil)
	})

	// Set up the router.
	mux.SetServeHTTP(customServeHTTP)
	mux.SetNotFound(notFound)
}
