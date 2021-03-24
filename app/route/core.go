// Package route provides the handlers for the application.
package route

import (
	"fmt"
	"log"
	"net/http"

	"github.com/josephspurrier/ambient/app/core"
	"github.com/josephspurrier/ambient/app/lib/htmltemplate"
	"github.com/josephspurrier/ambient/app/lib/router"
)

// Register all routes.
func Register(c *core.App) {
	// Register routes.
	registerHomePost(&HomePost{c})
	registerStyles(&Styles{c})
	registerAuthUtil(&AuthUtil{c})
	registerXMLUtil(&XMLUtil{c})
	registerAdminPost(&AdminPost{c})

	// This should be last because it catches all other pages at the root.
	registerPost(&Post{c})
}

// SetupRouter returns a router with the NotFound handler and the default
// handler set.
func SetupRouter(tmpl *htmltemplate.Engine) *router.Mux {
	// Set the handling of all responses.
	customServeHTTP := func(w http.ResponseWriter, r *http.Request, status int, err error) {
		// Handle only errors.
		if status >= 400 {
			vars := make(map[string]interface{})
			vars["title"] = fmt.Sprint(status)

			errTemplate := "400"

			switch status {
			case 403:
				// Already logged on plugin access denials.
				errTemplate = "403"
			case 404:
				// No need to log.
				errTemplate = "404"
			default:
				if err != nil {
					fmt.Println(err.Error())
				}
			}

			status, err = tmpl.ErrorTemplate(w, r, "layout/page", errTemplate, vars)
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
	rr := router.New(customServeHTTP, notFound)

	return rr
}
