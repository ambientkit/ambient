package app

import (
	"log"
	"os"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/plugin/author"
	"github.com/josephspurrier/ambient/plugin/awayrouter"
	"github.com/josephspurrier/ambient/plugin/bearblog"
	"github.com/josephspurrier/ambient/plugin/bearcss"
	"github.com/josephspurrier/ambient/plugin/charset"
	"github.com/josephspurrier/ambient/plugin/debugpprof"
	"github.com/josephspurrier/ambient/plugin/description"
	"github.com/josephspurrier/ambient/plugin/disqus"
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
)

var (
	// StorageSitePath is the location of the site file.
	StorageSitePath = "storage/site.json"
	// StorageSessionPath is the location of the session file.
	StorageSessionPath = "storage/session.bin"
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
		Router:         awayrouter.New(nil),
		TemplateEngine: htmltemplate.New(),
		// Trusted plugins are required to boot the application so they will be
		// given full access.
		TrustedPlugins: map[string]bool{
			"scssession": true, // Session manager.
			"plugins":    true, // Page to manage plugins.
			"bearblog":   true, // Bear Blog functionality.
			"bearcss":    true, // Bear Blog styling.
		},
		Plugins: []ambient.Plugin{
			// Marketplace plugins.
			debugpprof.New(),           // Go pprof debug endpoints.
			charset.New(),              // Charset to the HTML head.
			viewport.New(),             // Viewport in the HTML head.
			bearblog.New(passwordHash), // Bear Blog functionality.
			author.New(),               // Author in the HTML head.
			description.New(),          // Description the HTML head.
			bearcss.New(),              // Bear Blog styling.
			plugins.New(),              // Page to manage plugins.
			prism.New(),                // Prism CSS for codeblocks.
			stackedit.New(),            // Stackedit for editing markdown.
			googleanalytics.New(),      // Google Analytics.
			disqus.New(),               // Disqus for comments for blog posts.
			robots.New(),               // Robots file.
			sitemap.New(),              // Sitemap generator.
			rssfeed.New(),              // RSS feed generator.
			styles.New(),               // Style editing page.
		},
		Middleware: []ambient.MiddlewarePlugin{
			// Middleware - executes bottom to top.
			notrailingslash.New(),     // Redirect all requests with a trailing slash.
			uptimerobotok.New(),       // Provide 200 on HEAD /.
			securedashboard.New(),     // Secure all /dashboard routes.
			redirecttourl.New(),       // Redirect to production URL.
			gzipresponse.New(),        // Compress all HTTP responses.
			scssession.New(secretKey), // Session manager.
			logrequest.New(),          // Log every request as INFO.
		},
	}
}
