// Package awayrouter provides a router using a forked version of way
// for an Ambient application.
package awayrouter

import (
	"fmt"
	"net/http"

	"github.com/josephspurrier/ambient/app/core"
	"github.com/josephspurrier/ambient/plugin/awayrouter/router"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginMeta
	*core.Toolkit
}

// New returns a new awayrouter plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "awayrouter",
			Version:    "1.0.0",
			AppVersion: "1.0.0",
		},
	}
}

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *core.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

// Router returns a router.
func (p *Plugin) Router(logger core.ILogger, te core.IRender) (core.IAppRouter, error) {
	// Set up the default router.
	mux := router.New()

	// Set the NotFound and custom ServeHTTP handlers.
	setupRouter(logger, mux, te)

	return mux, nil
}

// setupRouter returns a router with the NotFound handler and the default
// handler set.
func setupRouter(logger core.ILogger, mux core.IAppRouter, te core.IRender) {
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
