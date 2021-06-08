package main

import (
	pkglog "log"
	"os"
	"os/exec"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/cmd/sveltekitgo/app"
	"github.com/josephspurrier/ambient/plugin/gcpbucketstorage"
	"github.com/josephspurrier/ambient/plugin/zaplogger"
)

var (
	appName    = "sveltekitgo"
	appVersion = "1.0"
)

func init() {
	// Set the time zone.
	tz := os.Getenv("AMB_TIMEZONE")
	if len(tz) > 0 {
		os.Setenv("TZ", tz)
	}
}

func main() {
	// Create the ambient app.
	plugins := app.Plugins()
	ambientApp, err := ambient.NewApp(appName, appVersion,
		zaplogger.New(),
		gcpbucketstorage.New(app.StorageSitePath, app.StorageSessionPath),
		plugins)
	if err != nil {
		pkglog.Fatalln(err.Error())
	}

	// Set the log level.
	// ambientApp.SetLogLevel(ambient.LogLevelDebug)
	// ambientApp.SetLogLevel(ambient.LogLevelInfo)
	// ambientApp.SetLogLevel(ambient.LogLevelError)
	// ambientApp.SetLogLevel(ambient.LogLevelFatal)

	// Add template debug information.
	//ambientApp.SetDebugTemplates(true)

	// Get the logger.
	log := ambientApp.Logger()

	// Enable the site plugins.
	app.GrantAccess(log, plugins, ambientApp)

	// Load the plugins and return the handler.
	mux, err := ambientApp.Handler()
	if err != nil {
		log.Fatal("", err.Error())
	}

	// Start node for the front-end.
	log.Info("ambient: web UI running on port: 8080")
	cmd := exec.Command("node", "../svelte-go/dist")
	cmd.Env = []string{"PORT=8080"}
	err = cmd.Start()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Start the web listener for the UI and API.
	ambientApp.ListenAndServe(mux)
}
