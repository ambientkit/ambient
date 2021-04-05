// Package awayrouter provides a router using a forked version of way
// for an Ambient application.
package awayrouter

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/josephspurrier/ambient/app/core"
	"github.com/josephspurrier/ambient/plugin/awayrouter/router"
)

//go:embed *
var assets embed.FS

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
func (p *Plugin) Router(te core.IRender) (core.IAppRouter, error) {
	// Set up the default router.
	mux := router.New()

	// Set the NotFound and custom ServeHTTP handlers.
	setupRouter(mux, te)

	return mux, nil
}

// setupRouter returns a router with the NotFound handler and the default
// handler set.
func setupRouter(mux core.IAppRouter, te core.IRender) {
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
					fmt.Println(err.Error())
				}
			default:
				if err != nil {
					fmt.Println(err.Error())
				}
			}

			status, err = te.PageContent(w, r, fmt.Sprintf("# %v\n%v", status, errText), nil, nil)
			if err != nil {
				if err != nil {
					log.Println(err.Error())
				}
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}
		}

		// Display server errors.
		if status >= 500 {
			if err != nil {
				log.Println(err.Error())
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
