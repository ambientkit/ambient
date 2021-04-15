// Package app represents an app.
package app

import (
	"log"
	"os"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/app/draft/hello"
	"github.com/josephspurrier/ambient/app/draft/navigation"
	"github.com/josephspurrier/ambient/plugin/author"
	"github.com/josephspurrier/ambient/plugin/awayrouter"
	"github.com/josephspurrier/ambient/plugin/bearblog"
	"github.com/josephspurrier/ambient/plugin/bearcss"
	"github.com/josephspurrier/ambient/plugin/charset"
	"github.com/josephspurrier/ambient/plugin/debugpprof"
	"github.com/josephspurrier/ambient/plugin/description"
	"github.com/josephspurrier/ambient/plugin/disqus"
	"github.com/josephspurrier/ambient/plugin/gcpbucketstorage"
	"github.com/josephspurrier/ambient/plugin/googleanalytics"
	"github.com/josephspurrier/ambient/plugin/gzipresponse"
	"github.com/josephspurrier/ambient/plugin/htmltemplate"
	"github.com/josephspurrier/ambient/plugin/logrequest"
	"github.com/josephspurrier/ambient/plugin/notrailingslash"
	"github.com/josephspurrier/ambient/plugin/plugins"
	"github.com/josephspurrier/ambient/plugin/prism"
	"github.com/josephspurrier/ambient/plugin/redirecttourl"
	"github.com/josephspurrier/ambient/plugin/robots"
	"github.com/josephspurrier/ambient/plugin/rssfeed"
	"github.com/josephspurrier/ambient/plugin/scssession"
	"github.com/josephspurrier/ambient/plugin/securedashboard"
	"github.com/josephspurrier/ambient/plugin/sitemap"
	"github.com/josephspurrier/ambient/plugin/stackedit"
	"github.com/josephspurrier/ambient/plugin/styles"
	"github.com/josephspurrier/ambient/plugin/uptimerobotok"
	"github.com/josephspurrier/ambient/plugin/viewport"
	"github.com/josephspurrier/ambient/plugin/zaplogger"
)

// Plugins defines the plugins - order does matter.
var Plugins = func() ambient.IPluginList {
	// Get the environment variables.
	secretKey := os.Getenv("AMB_SESSION_KEY")
	if len(secretKey) == 0 {
		log.Fatalf("app: environment variable missing: %v\n", "AMB_SESSION_KEY")
	}

	passwordHash := os.Getenv("AMB_PASSWORD_HASH")
	if len(passwordHash) == 0 {
		log.Fatalf("app: environment variable is missing: %v\n", "AMB_PASSWORD_HASH")
	}

	return ambient.IPluginList{
		// Core plugins required to use the system.
		//logruslogger.New(), // Logger must be the first plugin.
		zaplogger.New(),        // Logger must be the first plugin.
		gcpbucketstorage.New(), // GCP and local Storage must be the second plugin.
		htmltemplate.New(),     // HTML template engine.
		awayrouter.New(),       // Request router.

		// Additional plugins.
		debugpprof.New(),
		charset.New(),
		viewport.New(),
		bearblog.New(passwordHash),
		author.New(),
		description.New(),
		bearcss.New(),
		plugins.New(),
		prism.New(),
		stackedit.New(),
		googleanalytics.New(),
		disqus.New(),
		hello.New(),
		robots.New(),
		sitemap.New(),
		rssfeed.New(),
		styles.New(),
		navigation.New(),

		// Middleware - executes bottom to top.
		notrailingslash.New(),     // Redirect all request swith trailing slash.
		uptimerobotok.New(),       // Provide 200 on HEAD /.
		securedashboard.New(),     // Descure all /dashboard routes.
		redirecttourl.New(),       // Redirect to production URL.
		gzipresponse.New(),        // Compress all HTTP response.
		logrequest.New(),          // Log every request as INFO.
		scssession.New(secretKey), // Session manager.
	}
}
