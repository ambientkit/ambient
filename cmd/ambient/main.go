package main

import (
	"log"
	"os"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/app"
	"github.com/josephspurrier/ambient/plugin/gcpbucketstorage"
	"github.com/josephspurrier/ambient/plugin/zaplogger"
)

var (
	appName    = "myapp"
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
	ambientApp, err := ambient.NewApp(appName, appVersion, zaplogger.New(), gcpbucketstorage.New(), app.Plugins())
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Load the plugins.
	err = ambientApp.LoadPlugins()
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Start the web listener.
	ambientApp.ListenAndServe()
}
