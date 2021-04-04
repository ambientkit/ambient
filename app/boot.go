// Package app initializes all the services for the application.
package app

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/josephspurrier/ambient/app/core"
	"github.com/josephspurrier/ambient/app/lib/datastorage"
	"github.com/josephspurrier/ambient/app/lib/envdetect"
	"github.com/josephspurrier/ambient/app/lib/htmltemplate"
	"github.com/josephspurrier/ambient/app/lib/logger"
	"github.com/josephspurrier/ambient/app/model"
	"github.com/josephspurrier/ambient/app/route"
	"github.com/josephspurrier/ambient/html"
	"github.com/josephspurrier/ambient/plugin/author"
	"github.com/josephspurrier/ambient/plugin/awayrouter"
	"github.com/josephspurrier/ambient/plugin/bearcss"
	"github.com/josephspurrier/ambient/plugin/charset"
	"github.com/josephspurrier/ambient/plugin/description"
	"github.com/josephspurrier/ambient/plugin/disqus"
	"github.com/josephspurrier/ambient/plugin/googleanalytics"
	"github.com/josephspurrier/ambient/plugin/gzipresponse"
	"github.com/josephspurrier/ambient/plugin/hello"
	"github.com/josephspurrier/ambient/plugin/home"
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

var (
	storageSitePath    = "storage/site.json"
	storageSessionPath = "storage/session.bin"
)

// Boot returns a router with the application ready to be started.
func Boot(l *logger.Logger) (http.Handler, error) {
	// Set the storage and session environment variables.
	sitePath := os.Getenv("AMB_SITE_PATH")
	if len(sitePath) > 0 {
		storageSitePath = sitePath
	}

	// Get the environment variables.
	secretKey := os.Getenv("AMB_SESSION_KEY")
	if len(secretKey) == 0 {
		return nil, fmt.Errorf("environment variable missing: %v", "AMB_SESSION_KEY")
	}

	bucket := os.Getenv("AMB_GCP_BUCKET_NAME")
	if len(bucket) == 0 {
		return nil, fmt.Errorf("environment variable missing: %v", "AMB_GCP_BUCKET_NAME")
	}

	allowHTML, err := strconv.ParseBool(os.Getenv("AMB_ALLOW_HTML"))
	if err != nil {
		return nil, fmt.Errorf("environment variable not able to parse as bool: %v", "AMB_ALLOW_HTML")
	}

	// Create new store object with the defaults.
	site := &model.Site{}

	var ds datastorage.Datastorer
	var ss core.SessionStorer

	if !envdetect.RunningLocalDev() {
		// Use Google when running in GCP.
		ds = datastorage.NewGCPStorage(bucket, storageSitePath)
		ss = datastorage.NewGCPStorage(bucket, storageSessionPath)
	} else {
		// Use local filesytem when developing.
		ds = datastorage.NewLocalStorage(storageSitePath)
		ss = datastorage.NewLocalStorage(storageSessionPath)
	}

	// Set up the data storage provider.
	storage, err := datastorage.New(ds, site)
	if err != nil {
		return nil, err
	}

	// Define the plugins.
	arrPlugins := []core.IPlugin{
		awayrouter.New(), // Router.

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
		home.New(),
		navigation.New(),

		// Middleware - executes bottom to top.
		notrailingslash.New(),
		uptimerobotok.New(),
		securedashboard.New(),
		redirecttourl.New(),
		gzipresponse.New(),
		logrequest.New(),
		scssession.New(), // Session manager.
	}

	pluginNames := make([]string, 0)
	for _, v := range arrPlugins {
		pluginNames = append(pluginNames, v.PluginName())
	}

	// Register the plugins.
	plugs, err := core.RegisterPlugins(arrPlugins, storage)
	if err != nil {
		return nil, err
	}

	// TODO: Need to have a default session handler that just throws messages.
	var sess core.ISession
	var mux core.IAppRouter

	// Get the session manager from the plugins.
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
		sm, err := v.SessionManager(ss, secretKey)
		if err != nil {
			l.Error("", err.Error())
		} else if sm != nil && sess == nil {
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

	// Set the session manager if one doesn't exist.
	// var defaultSessionManager *scssession.Plugin
	// if sess == nil {
	// 	// Set up the default session manager.
	// 	defaultSessionManager = scssession.New()
	// 	sess, err = defaultSessionManager.SessionManager(ss, secretKey)
	// 	if err != nil {
	// 		l.Fatal("boot: default session manager cannot be loaded: %v", err.Error())
	// 	}
	// }

	// Set up the template engine.
	tm := html.NewTemplateManager(storage, sess)
	pi := core.NewPlugininjector(storage, sess, plugs)
	tmpl := htmltemplate.New(allowHTML, tm, pi, pluginNames)

	// Get the router from the plugins.
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
		rm, err := v.Router(tmpl)
		if err != nil {
			l.Error("", err.Error())
		} else if rm != nil {
			// Only set the router once.
			l.Info("boot: using router (mux) from plugin: %v", v.PluginName())
			mux = rm
			break
		}
	}

	// FIXME: Need to fail gracefully.
	if mux == nil {
		l.Fatal("boot: no default router found")
	}

	// Set the router if one doesn't exist.
	// if mux == nil {
	// 	// Set up the default router.
	// 	ar := awayrouter.New()
	// 	ar.Router(tmpl)
	// }

	// Create core app.
	c := core.NewApp(l, plugs, tmpl, mux, sess, storage)

	// Load the plugin pages.
	err = c.LoadAllPluginPages()
	if err != nil {
		return nil, err
	}

	// Setup the routes.
	route.Register(c)

	// Enable the middleware from the plugins.
	var h http.Handler = c.Router
	h = c.LoadAllPluginMiddleware(h, arrPlugins)

	return h, nil
}
