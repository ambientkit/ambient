package main

import (
	stdlog "log"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/cmd/polarbearblog/app"
	"github.com/josephspurrier/ambient/plugin/gcpbucketstorage"
	"github.com/josephspurrier/ambient/plugin/zaplogger"
)

var (
	appName    = "polarbearblog"
	appVersion = "1.0"
)

const (
	// LoginURL allows user to set the login URL.
	LoginURL = "Login URL"
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

	// Get the logger.
	log := ambientApp.Logger()

	// Enable the trusted site plugins.
	ambientApp.GrantAccess(plugins)

	// Load the plugins and return the handler.
	mux, err := ambientApp.Handler()
	if err != nil {
		log.Fatal("", err.Error())
	}

	// Start the web listener.
	ambientApp.ListenAndServe(mux)
}
