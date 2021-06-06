// Package app represents an Ambient app.
package app

import (
	"log"
	"os"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/cmd/sveltekitgo/draft/webapi"
	"github.com/josephspurrier/ambient/plugin/awayrouter"
	"github.com/josephspurrier/ambient/plugin/gzipresponse"
	"github.com/josephspurrier/ambient/plugin/htmltemplate"
	"github.com/josephspurrier/ambient/plugin/logrequest"
	"github.com/josephspurrier/ambient/plugin/notrailingslash"
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

	return &ambient.PluginLoader{
		Router:         awayrouter.New(),
		TemplateEngine: htmltemplate.New(),
		// Trusted plugins are required to boot the application so they will be
		// given full access.
		TrustedPlugins: map[string]bool{
			"scssession": true,
			"webapi":     true,

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
			notrailingslash.New(),     // Redirect all requests with a trailing slash.
			gzipresponse.New(),        // Compress all HTTP responses.
			scssession.New(secretKey), // Session manager.
			logrequest.New(),          // Log every request as INFO.
		},
	}
}
