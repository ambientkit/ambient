// Package app initializes all the services for the application.
package app

import (
	"fmt"
	"net/http"

	"github.com/josephspurrier/ambient/app/core"
	"github.com/josephspurrier/ambient/plugin/author"
	"github.com/josephspurrier/ambient/plugin/awayrouter"
	"github.com/josephspurrier/ambient/plugin/bearblog"
	"github.com/josephspurrier/ambient/plugin/bearcss"
	"github.com/josephspurrier/ambient/plugin/charset"
	"github.com/josephspurrier/ambient/plugin/description"
	"github.com/josephspurrier/ambient/plugin/disqus"
	"github.com/josephspurrier/ambient/plugin/gcpbucketstorage"
	"github.com/josephspurrier/ambient/plugin/googleanalytics"
	"github.com/josephspurrier/ambient/plugin/gzipresponse"
	"github.com/josephspurrier/ambient/plugin/hello"
	"github.com/josephspurrier/ambient/plugin/htmltemplate"
	"github.com/josephspurrier/ambient/plugin/logrequest"
	"github.com/josephspurrier/ambient/plugin/logruslogger"
	"github.com/josephspurrier/ambient/plugin/navigation"
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

// Plugins defines the plugins - order does matter.
var Plugins = core.IPluginList{
	// Core plugins required to use the system.
	logruslogger.New(),     // Logger must be the first plugin.
	gcpbucketstorage.New(), // GCP and local Storage must be the second plugin.
	htmltemplate.New(),     // HTML template engine.
	awayrouter.New(),       // Request router.

	// Additional plugins.
	charset.New(),
	viewport.New(),
	bearblog.New(),
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
	notrailingslash.New(), // Redirect all request swith trailing slash.
	uptimerobotok.New(),   // Provide 200 on HEAD /.
	securedashboard.New(), // Descure all /dashboard routes.
	redirecttourl.New(),   // Redirect to production URL.
	gzipresponse.New(),    // Compress all HTTP response.
	logrequest.New(),      // Log every request as INFO.
	scssession.New(),      // Session manager.
}

// Boot returns a router with the application ready to be started.
func Boot() (core.IAppLogger, http.Handler, error) {
	// Ensure there is at least the storage plugin.
	if len(Plugins) == 0 {
		return nil, nil, fmt.Errorf("boot: no plugins found")
	}

	// Set up the logger.
	log, err := Logger("ambient", "1.0", Plugins[0])
	if err != nil {
		return nil, nil, err
	}

	// Set up the storage.
	storage, ss, err := Storage(log, Plugins[1])
	if err != nil {
		return log, nil, err
	}

	// Initialize the plugin system.
	ps, err := core.NewPluginSystem(log, Plugins, storage)
	if err != nil {
		return log, nil, err
	}

	// Get the session manager from the plugins.
	var sess core.IAppSession
	for _, name := range ps.Names() {
		// Get the plugin.
		p, err := ps.Plugin(name)
		if err != nil {
			log.Error("boot: could not find plugin (%v): %v", name, err.Error())
			continue
		}

		// Skip if the plugin isn't enabled.
		if !ps.Enabled(name) {
			continue
		}

		// Get the session manager.
		sm, err := p.SessionManager(log, ss)
		if err != nil {
			log.Error("", err.Error())
		} else if sm != nil {
			// Only set the session manager once.
			log.Info("boot: using session manager from plugin: %v", name)
			sess = sm
			break
		}
	}

	// FIXME: Need to fail gracefully?
	if sess == nil {
		log.Fatal("boot: no default session manager found")
	}

	// Set up the template injector.
	pi := core.NewPlugininjector(log, storage, sess, ps, Plugins)

	// Get the router from the plugins.
	var te core.IRender
	for _, name := range ps.Names() {
		// Skip if the plugin isn't found.
		plugin, err := ps.Plugin(name)
		if err != nil {
			continue
		}

		// Skip if the plugin isn't enabled.
		if !ps.Enabled(name) {
			continue
		}

		// Get the router.
		tt, err := plugin.TemplateEngine(log, pi)
		if err != nil {
			log.Error("", err.Error())
		} else if tt != nil {
			// Only set the router once.
			log.Info("boot: using template engine from plugin: %v", name)
			te = tt
			break
		}
	}
	if te == nil {
		log.Fatal("boot: no default template engine found")
	}

	// Get the router from the plugins.
	var mux core.IAppRouter
	for _, name := range ps.Names() {
		// Skip if the plugin isn't found.
		plugin, err := ps.Plugin(name)
		if err != nil {
			continue
		}

		// Skip if the plugin isn't enabled.
		if !ps.Enabled(name) {
			continue
		}

		// Get the router.
		rm, err := plugin.Router(log, te)
		if err != nil {
			log.Error("", err.Error())
		} else if rm != nil {
			// Only set the router once.
			log.Info("boot: using router (mux) from plugin: %v", name)
			mux = rm
			break
		}
	}
	if mux == nil {
		log.Fatal("boot: no default router found")
	}

	// Create secure site for the core application.
	securesite := core.NewSecureSite("ambient", log, storage, ps, sess, mux, te)

	// Load the plugin pages.
	err = securesite.LoadAllPluginPages()
	if err != nil {
		return log, nil, err
	}

	// Enable the middleware from the plugins.
	h := securesite.LoadAllPluginMiddleware(mux)

	return log, h, nil
}
