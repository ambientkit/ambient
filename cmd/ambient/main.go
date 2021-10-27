package main

import (
	stdlog "log"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/app"
	"github.com/josephspurrier/ambient/plugin/gcpbucketstorage"
	"github.com/josephspurrier/ambient/plugin/zaplogger"
)

var (
	appName    = "myapp"
	appVersion = "1.0"
)

func main() {
	// Create the ambient app.
	plugins := app.Plugins()
	ambientApp, err := ambient.NewApp(appName, appVersion,
		zaplogger.New(),
		gcpbucketstorage.New(app.StorageSitePath, app.StorageSessionPath),
		plugins)
	if err != nil {
		stdlog.Fatalln(err.Error())
	}

	// Set the log level.
	// ambientApp.SetLogLevel(ambient.LogLevelDebug)
	// ambientApp.SetLogLevel(ambient.LogLevelInfo)
	// ambientApp.SetLogLevel(ambient.LogLevelError)
	// ambientApp.SetLogLevel(ambient.LogLevelFatal)

	// Add template debug information.
	// ambientApp.SetDebugTemplates(true)

	// Enable the trusted plugins.
	ambientApp.GrantAccess(plugins)

	// Get the logger.
	log := ambientApp.Logger()

	// Load the plugins and return the handler.
	mux, err := ambientApp.Handler()
	if err != nil {
		log.Fatal("", err.Error())
	}

	// Start the web listener.
	ambientApp.ListenAndServe(mux)
}
