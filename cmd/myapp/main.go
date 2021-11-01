package main

import (
	stdlog "log"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/cmd/myapp/app"
	"github.com/josephspurrier/ambient/lib/envdetect"
	"github.com/josephspurrier/ambient/plugin/logger/zaplogger"
	"github.com/josephspurrier/ambient/plugin/storage/awsbucketstorage"
	"github.com/josephspurrier/ambient/plugin/storage/azureblobstorage"
	"github.com/josephspurrier/ambient/plugin/storage/gcpbucketstorage"
	"github.com/josephspurrier/ambient/plugin/storage/localstorage"
)

var (
	appName    = "myapp"
	appVersion = "1.0"
)

func init() {
	// Verbose logging with file name and line number for the standard logger.
	stdlog.SetFlags(stdlog.Lshortfile)
}

func main() {
	// Select the storage engine for site and session information.
	var storage ambient.StoragePlugin
	if envdetect.RunningLocalDev() {
		storage = localstorage.New(app.StorageSitePath, app.StorageSessionPath)
	} else if envdetect.RunningInGoogle() {
		storage = gcpbucketstorage.New(app.StorageSitePath, app.StorageSessionPath)
	} else if envdetect.RunningInAWS() {
		storage = awsbucketstorage.New(app.StorageSitePath, app.StorageSessionPath)
	} else if envdetect.RunningInAzureFunction() {
		storage = azureblobstorage.New(app.StorageSitePath, app.StorageSessionPath)
	} else {
		// Defaulting to local storage.
		storage = localstorage.New(app.StorageSitePath, app.StorageSessionPath)
	}

	// Create the ambient app.
	plugins := app.Plugins()
	ambientApp, log, err := ambient.NewApp(appName, appVersion,
		zaplogger.New(),
		storage,
		plugins)
	if err != nil {
		if log != nil {
			// Use the logger if it's available.
			log.Fatal("", err.Error())
		} else {
			// Else use the standard logger.
			stdlog.Fatalln(err.Error())
		}
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

	// Load the plugins and return the handler.
	mux, err := ambientApp.Handler()
	if err != nil {
		log.Fatal("", err.Error())
	}

	// Start the web listener.
	ambientApp.ListenAndServe(mux)
}
