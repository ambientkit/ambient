// Package app represents an Ambient app.
package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/cmd/sveltekitgo/draft/webapi"
	"github.com/josephspurrier/ambient/plugin/awayrouter"
	"github.com/josephspurrier/ambient/plugin/gzipresponse"
	"github.com/josephspurrier/ambient/plugin/htmltemplate"
	"github.com/josephspurrier/ambient/plugin/logrequest"
	"github.com/josephspurrier/ambient/plugin/notrailingslash"
	"github.com/josephspurrier/ambient/plugin/proxyrequest"
	"github.com/josephspurrier/ambient/plugin/scssession"
)

var (
	// StorageSitePath is the location of the site file.
	StorageSitePath = "cmd/sveltekitgo/storage/site.json"
	// StorageSessionPath is the location of the session file.
	StorageSessionPath = "cmd/sveltekitgo/storage/session.bin"
)

// Plugins defines the plugins - order does matter.
var Plugins = func() *ambient.PluginLoader {
	// Get the environment variables.
	secretKey := os.Getenv("AMB_SESSION_KEY")
	if len(secretKey) == 0 {
		log.Fatalf("app: environment variable missing: %v\n", "AMB_SESSION_KEY")
	}

	passwordHash := os.Getenv("AMB_PASSWORD_HASH")
	if len(passwordHash) == 0 {
		log.Fatalf("app: environment variable is missing: %v\n", "AMB_PASSWORD_HASH")
	}

	// Front-end proxy.
	urlUI, err := url.Parse("http://localhost:8080")
	if err != nil {
		log.Fatalf("app: UI proxy target error: %v", err.Error())
	}

	return &ambient.PluginLoader{
		Router:         awayrouter.New(ErrorHandler()),
		TemplateEngine: htmltemplate.New(),
		// Trusted plugins are required to boot the application so they will be
		// given full access.
		TrustedPlugins: map[string]bool{
			"scssession": true,
			"webapi":     true,

			// Middleware.
			"proxyrequest":    true,
			"notrailingslash": true,
			"gzipresponse":    true,
			"logrequest":      true,
		},
		Plugins: []ambient.Plugin{
			// Marketplace plugins.

			// App plugins.
			webapi.New(), // REST API.
		},
		Middleware: []ambient.MiddlewarePlugin{
			// Middleware - executes bottom to top.
			proxyrequest.New(urlUI, "/api"),
			notrailingslash.New(),     // Redirect all requests with a trailing slash.
			gzipresponse.New(),        // Compress all HTTP responses.
			scssession.New(secretKey), // Session manager.
			logrequest.New(),          // Log every request as INFO.
		},
	}
}

// ErrorResponse is an API error response.
type ErrorResponse struct {
	Status          int    `json:"status"`
	Message         string `json:"message"`
	FriendlyMessage string `json:"friendlyMessage"`
}

// ErrorHandler returns the JSON error handler for the router.
func ErrorHandler() awayrouter.LoggerHandler {
	return func(logger ambient.Logger, w http.ResponseWriter, r *http.Request, status int, err error) {
		// Handle only errors.
		if status >= 400 {
			errText := strings.ToLower(http.StatusText(status))

			switch status {
			case 403:
				// Already logged on plugin access denials.
				errText = "a plugin has been denied permission"
			case 404:
				// No need to log.
				errText = "darn, we cannot find the page"
			case 400:
				errText = "darn, something went wrong"
				if err != nil {
					errText = err.Error()
					logger.Info("awayrouter: error (%v): %v", status, err.Error())
				}
			default:
				if err != nil {
					logger.Info("awayrouter: error (%v): %v", status, err.Error())
				}
			}

			b, err := json.Marshal(ErrorResponse{
				Status:          status,
				Message:         errText,
				FriendlyMessage: errText,
			})
			if err != nil {
				if err != nil {
					logger.Info("awayrouter: error in rendering error response (%v): %v", status, err.Error())
				}
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			fmt.Fprint(w, string(b))
		}
	}
}
