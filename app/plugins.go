package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/josephspurrier/ambient/app/lib/ambsystem"
	"github.com/josephspurrier/ambient/app/lib/datastorage"
	"github.com/josephspurrier/ambient/app/lib/envdetect"
	"github.com/josephspurrier/ambient/app/lib/htmltemplate"
	"github.com/josephspurrier/ambient/app/lib/router"
	"github.com/josephspurrier/ambient/app/lib/websession"
	"github.com/josephspurrier/ambient/app/model"
	"github.com/josephspurrier/ambient/app/route"
	"github.com/josephspurrier/ambient/html"
	"github.com/josephspurrier/ambient/plugin/stackedit"
)

// LoadPlugins will load the plugins.
func LoadPlugins() *router.Mux {
	c, err := boot()
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Define the plugins.
	plugins := []ambsystem.IPlugin{
		stackedit.Activate(),
	}

	// Load the plugins.
	for _, v := range plugins {
		fmt.Printf("Load plugin: %v\n", v.PluginName())
		v.SetPages(c.Router)
	}

	return c.Router
}

// func setLogging() {
// 	//temp := os.Stdout
// 	// Turn off logging so it plugins don't have control over output.
// 	os.Stdout = nil
// 	//os.Stdout = temp   // restore it
// }

func boot() (*route.Core, error) {
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

	// Set up the template engine.
	tm := html.NewTemplateManager(storage, sess)
	tmpl := htmltemplate.New(tm, allowHTML)

	// Create core app.
	c := &route.Core{
		Router:  setupRouter(tmpl),
		Storage: storage,
		Render:  tmpl,
		Sess:    sess,
	}

	return c, nil
}

func setupRouter(tmpl *htmltemplate.Engine) *router.Mux {
	// Set the handling of all responses.
	customServeHTTP := func(w http.ResponseWriter, r *http.Request, status int, err error) {
		// Handle only errors.
		if status >= 400 {
			// vars := make(map[string]interface{})
			// vars["title"] = fmt.Sprint(status)
			// errTemplate := "400"
			// if status == 404 {
			// 	errTemplate = "404"
			// }
			// status, err = tmpl.ErrorTemplate(w, r, "base", errTemplate, vars)
			// if err != nil {
			// 	log.Println(err.Error())
			// 	http.Error(w, "500 internal server error", http.StatusInternalServerError)
			// 	return
			// }
			http.Error(w, http.StatusText(status), status)
		}

		// Display server errors.
		if status >= 500 {
			if err != nil {
				log.Println(err.Error())
			}
			http.Error(w, http.StatusText(status), status)
		}
	}

	// Send all 404 to the customer handler.
	notFound := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customServeHTTP(w, r, http.StatusNotFound, nil)
	})

	// Set up the router.
	return router.New(customServeHTTP, notFound)
}
