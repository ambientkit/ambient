// Package app initializes all the services for the application.
package app

import (
	"net/http"

	"github.com/josephspurrier/ambient/app/core"
	"github.com/josephspurrier/ambient/app/lib/logger"
	"github.com/josephspurrier/ambient/html"
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

// Boot returns a router with the application ready to be started.
func Boot(l *logger.Logger) (http.Handler, error) {
	// Define the plugins - order does matter.
	arrPlugins := []core.IPlugin{
		// Core plugins required to use the system.
		gcpbucketstorage.New(), // GCP and local Storage - storage plugin must always come first.
		htmltemplate.New(),     // HTML template engine.
		awayrouter.New(),       // Request router.

		// Additional plugins.
		bearblog.New(),
		charset.New(),
		viewport.New(),
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

	// Create a list of the plugin names.
	pluginNames := make([]string, 0)
	for _, v := range arrPlugins {
		pluginNames = append(pluginNames, v.PluginName())
	}

	// Define the storage managers.
	var ds core.DataStorer
	var ss core.SessionStorer

	// Get the storage manager from the plugins.
	// This must be the first plugin or else it fails.
	if len(arrPlugins) > 0 {
		firstPlugin := arrPlugins[0]
		// Get the storage system.
		pds, pss, err := firstPlugin.Storage()
		if err != nil {
			l.Error("", err.Error())
		} else if pds != nil && pss != nil {
			l.Info("boot: using storage from first plugin: %v", firstPlugin.PluginName())
			ds = pds
			ss = pss
		}
	}
	if ds == nil || ss == nil {
		l.Fatal("boot: no default storage found")
	}

	// Create new store object with the defaults.
	site := &core.Site{}

	// Set up the data storage provider.
	storage, err := core.NewDatastore(ds, site)
	if err != nil {
		return nil, err
	}

	// Register the plugins.
	plugs, err := core.RegisterPlugins(arrPlugins, storage)
	if err != nil {
		return nil, err
	}

	// Get the session manager from the plugins.
	var sess core.ISession
	for _, v := range arrPlugins {
		// Skip if the plugin isn't found.
		ps, ok := storage.Site.PluginSettings[v.PluginName()]
		if !ok {
			continue
		}

		// Skip if the plugin isn't enabled.
		if !ps.Enabled {
			continue
		}

		// Get the session manager.
		sm, err := v.SessionManager(ss)
		if err != nil {
			l.Error("", err.Error())
		} else if sm != nil {
			// Only set the session manager once.
			l.Info("boot: using session manager from plugin: %v", v.PluginName())
			sess = sm
			break
		}
	}

	// FIXME: Need to fail gracefully.
	if sess == nil {
		l.Fatal("boot: no default session manager found")
	}

	// Set up the template engine.
	pi := core.NewPlugininjector(storage, sess, plugs)
	templateManager := html.NewTemplateManager(storage, sess)

	// Get the router from the plugins.
	var te core.IRender
	for _, v := range arrPlugins {
		// Skip if the plugin isn't found.
		ps, ok := storage.Site.PluginSettings[v.PluginName()]
		if !ok {
			continue
		}

		// Skip if the plugin isn't enabled.
		if !ps.Enabled {
			continue
		}

		// Get the router.
		tt, err := v.TemplateEngine(templateManager, pi, pluginNames)
		if err != nil {
			l.Error("", err.Error())
		} else if tt != nil {
			// Only set the router once.
			l.Info("boot: using template engine from plugin: %v", v.PluginName())
			te = tt
			break
		}
	}
	if te == nil {
		l.Fatal("boot: no default template engine found")
	}

	// Get the router from the plugins.
	var mux core.IAppRouter
	for _, v := range arrPlugins {
		// Skip if the plugin isn't found.
		ps, ok := storage.Site.PluginSettings[v.PluginName()]
		if !ok {
			continue
		}

		// Skip if the plugin isn't enabled.
		if !ps.Enabled {
			continue
		}

		// Get the router.
		rm, err := v.Router(te)
		if err != nil {
			l.Error("", err.Error())
		} else if rm != nil {
			// Only set the router once.
			l.Info("boot: using router (mux) from plugin: %v", v.PluginName())
			mux = rm
			break
		}
	}
	if mux == nil {
		l.Fatal("boot: no default router found")
	}

	// Create core app.
	c := core.NewApp(l, plugs, te, mux, sess, storage)

	// Load the plugin pages.
	err = c.LoadAllPluginPages()
	if err != nil {
		return nil, err
	}

	// Enable the middleware from the plugins.
	var h http.Handler = c.Router
	h = c.LoadAllPluginMiddleware(h, arrPlugins)

	return h, nil
}
