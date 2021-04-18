package main

import (
	"log"
	"os"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/app"
	"github.com/josephspurrier/ambient/plugin/gcpbucketstorage"
	"github.com/josephspurrier/ambient/plugin/zaplogger"
)

func init() {
	// Set the time zone.
	tz := os.Getenv("AMB_TIMEZONE")
	if len(tz) > 0 {
		os.Setenv("TZ", tz)
	}
}

func main() {
	//logger := logruslogger.New()         // Logger
	logger := zaplogger.New()            // Logger
	gcpstorage := gcpbucketstorage.New() // GCP and local Storage must be the second plugin.

	// Create the ambient app.
	ambientApp, err := ambient.NewApp("ambient", "1.0", logger, gcpstorage, app.Plugins())
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Set up the plugins.
	err = ambientApp.LoadPlugins()
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Start the web listener.
	ambientApp.ListenAndServe()
}
