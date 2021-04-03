// Package app initializes all the services for the application.
package app

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/josephspurrier/ambient/app/core"
	"github.com/josephspurrier/ambient/app/lib/datastorage"
	"github.com/josephspurrier/ambient/app/lib/envdetect"
	"github.com/josephspurrier/ambient/app/lib/htmltemplate"
	"github.com/josephspurrier/ambient/app/lib/logger"
	"github.com/josephspurrier/ambient/app/lib/websession"
	"github.com/josephspurrier/ambient/app/middleware"
	"github.com/josephspurrier/ambient/app/model"
	"github.com/josephspurrier/ambient/app/route"
	"github.com/josephspurrier/ambient/html"
	"github.com/josephspurrier/ambient/plugin/author"
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
	"github.com/josephspurrier/ambient/plugin/plugins"
	"github.com/josephspurrier/ambient/plugin/prism"
	"github.com/josephspurrier/ambient/plugin/robots"
	"github.com/josephspurrier/ambient/plugin/rssfeed"
	"github.com/josephspurrier/ambient/plugin/sitemap"
	"github.com/josephspurrier/ambient/plugin/stackedit"
	"github.com/josephspurrier/ambient/plugin/styles"
	"github.com/josephspurrier/ambient/plugin/viewport"
)

var (
	storageSitePath    = "storage/site.json"
	storageSessionPath = "storage/session.bin"
	sessionName        = "session"
)

// Boot returns a router with the application ready to be started.
func Boot(l *logger.Logger) (http.Handler, error) {
	// Set the storage and session environment variables.
	sitePath := os.Getenv("AMB_SITE_PATH")
	if len(sitePath) > 0 {
		storageSitePath = sitePath
	}

	sessionPath := os.Getenv("AMB_SESSION_PATH")
	if len(sessionPath) > 0 {
		storageSessionPath = sessionPath
	}

	sname := os.Getenv("AMB_SESSION_NAME")
	if len(sname) > 0 {
		sessionName = sname
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
	var ss websession.Sessionstorer

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

	// Set up the session storage provider.
	en := websession.NewEncryptedStorage(secretKey)
	store, err := websession.NewJSONSession(ss, en)
	if err != nil {
		return nil, err
	}

	// Initialize a new session manager and configure the session lifetime.
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = false
	sessionManager.Store = store
	sess := websession.New(sessionName, sessionManager)

	// Define the plugins.
	arrPlugins := []core.IPlugin{
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
		gzipresponse.New(),
		logrequest.New(),
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

	// Set up the template engine.
	tm := html.NewTemplateManager(storage, sess)
	pi := core.NewPlugininjector(storage, sess, plugs)
	tmpl := htmltemplate.New(allowHTML, tm, pi, pluginNames)

	// Set up the router.
	mux := route.SetupRouter(tmpl)

	// Create core app.
	c := core.NewApp(l, plugs, tmpl, mux, sess, storage)

	// Load the plugin pages.
	err = c.LoadAllPluginPages()
	if err != nil {
		return nil, err
	}

	// Setup the routes.
	route.Register(c)

	// Setup the middleware.
	mw := middleware.NewHandler(c.Render, c.Sess, c.Router, c.Storage.Site.URL, c.Storage.Site.Scheme)
	var h http.Handler = c.Router
	arrMiddleware := []func(next http.Handler) http.Handler{
		mw.Redirect,
		middleware.Head,
		mw.DisallowAnon,
		sessionManager.LoadAndSave,
	}
	for _, v := range arrMiddleware {
		h = v(h)
	}

	// Enable the middleware from the plugins.
	h = c.LoadAllPluginMiddleware(h, arrPlugins)

	return h, nil
}
