package app

import (
	"log"
	"os"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/cmd/myapp/app/draft/hello"
	"github.com/josephspurrier/ambient/cmd/myapp/app/draft/navigation"
	"github.com/josephspurrier/ambient/plugin/generic/author"
	"github.com/josephspurrier/ambient/plugin/generic/bearblog"
	"github.com/josephspurrier/ambient/plugin/generic/bearcss"
	"github.com/josephspurrier/ambient/plugin/generic/charset"
	"github.com/josephspurrier/ambient/plugin/generic/debugpprof"
	"github.com/josephspurrier/ambient/plugin/generic/description"
	"github.com/josephspurrier/ambient/plugin/generic/disqus"
	"github.com/josephspurrier/ambient/plugin/generic/envinfo"
	"github.com/josephspurrier/ambient/plugin/generic/googleanalytics"
	"github.com/josephspurrier/ambient/plugin/generic/pluginmanager"
	"github.com/josephspurrier/ambient/plugin/generic/prism"
	"github.com/josephspurrier/ambient/plugin/generic/robots"
	"github.com/josephspurrier/ambient/plugin/generic/rssfeed"
	"github.com/josephspurrier/ambient/plugin/generic/simplelogin"
	"github.com/josephspurrier/ambient/plugin/generic/sitemap"
	"github.com/josephspurrier/ambient/plugin/generic/stackedit"
	"github.com/josephspurrier/ambient/plugin/generic/styles"
	"github.com/josephspurrier/ambient/plugin/generic/viewport"
	"github.com/josephspurrier/ambient/plugin/middleware/etagcache"
	"github.com/josephspurrier/ambient/plugin/middleware/gzipresponse"
	"github.com/josephspurrier/ambient/plugin/middleware/logrequest"
	"github.com/josephspurrier/ambient/plugin/middleware/notrailingslash"
	"github.com/josephspurrier/ambient/plugin/middleware/redirecttourl"
	"github.com/josephspurrier/ambient/plugin/middleware/securedashboard"
	"github.com/josephspurrier/ambient/plugin/middleware/uptimerobotok"
	"github.com/josephspurrier/ambient/plugin/router/awayrouter"
	"github.com/josephspurrier/ambient/plugin/sessionmanager/scssession"
	"github.com/josephspurrier/ambient/plugin/templateengine/htmlengine"
)

var (
	// StorageSitePath is the location of the site file.
	StorageSitePath = "storage/site.bin"
	// StorageSessionPath is the location of the session file.
	StorageSessionPath = "storage/session.bin"
)

// Plugins defines the plugins.
func Plugins() *ambient.PluginLoader {
	// Get the environment variables.
	secretKey := os.Getenv("AMB_SESSION_KEY")
	if len(secretKey) == 0 {
		log.Fatalf("app: environment variable missing: %v\n", "AMB_SESSION_KEY")
	}

	passwordHash := os.Getenv("AMB_PASSWORD_HASH")
	if len(passwordHash) == 0 {
		log.Fatalf("app: environment variable is missing: %v\n", "AMB_PASSWORD_HASH")
	}

	// Define the session manager so it can be used as a core plugin and
	// middleware.
	sessionManager := scssession.New(secretKey)

	return &ambient.PluginLoader{
		// Core plugins are implicitly trusted.
		Router:         awayrouter.New(nil),
		TemplateEngine: htmlengine.New(),
		SessionManager: sessionManager,
		// Trusted plugins are those that are typically needed to boot so they
		// will be enabled and given full access.
		TrustedPlugins: map[string]bool{
			"pluginmanager": true, // Page to manage plugins.
			"simplelogin":   true, // Simple login page.
			"bearcss":       true, // Bear Blog styling.
		},
		Plugins: []ambient.Plugin{
			// Marketplace plugins.
			charset.New(),                 // Charset to the HTML head.
			simplelogin.New(passwordHash), // Simple login page.
			bearblog.New(passwordHash),    // Bear Blog functionality.
			bearcss.New(),                 // Bear Blog styling.
			debugpprof.New(),              // Go pprof debug endpoints.
			viewport.New(),                // Viewport in the HTML head.
			author.New(),                  // Author in the HTML head.
			description.New(),             // Description the HTML head.
			pluginmanager.New(),           // Page to manage plugins.
			prism.New(),                   // Prism CSS for codeblocks.
			stackedit.New(),               // Stackedit for editing markdown.
			googleanalytics.New(),         // Google Analytics.
			disqus.New(),                  // Disqus for comments for blog posts.
			robots.New(),                  // Robots file.
			sitemap.New(),                 // Sitemap generator.
			rssfeed.New(),                 // RSS feed generator.
			styles.New(),                  // Style editing page.
			envinfo.New(),                 // Show environment variables on the server.

			// App plugins.
			hello.New(),
			navigation.New(),
		},
		Middleware: []ambient.MiddlewarePlugin{
			// Middleware - executes bottom to top.
			notrailingslash.New(), // Redirect all requests with a trailing slash.
			uptimerobotok.New(),   // Provide 200 on HEAD /.
			securedashboard.New(), // Secure all /dashboard routes.
			redirecttourl.New(),   // Redirect to production URL.
			etagcache.New(),       // Cache using Etag.
			gzipresponse.New(),    // Compress all HTTP responses.
			sessionManager,        // Session manager middleware.
			logrequest.New(),      // Log every request as INFO.
		},
	}
}
