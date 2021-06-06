package main

import (
	pkglog "log"
	"os"

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
	grantAccess(log, plugins, ambientApp)

	// Load the plugins and return the handler.
	mux, err := ambientApp.Handler()
	if err != nil {
		log.Fatal("", err.Error())
	}

	// Start the web listener.
	ambientApp.ListenAndServe(mux)
}

func grantAccess(log ambient.AppLogger, plugins *ambient.PluginLoader, ambientApp *ambient.App) {
	// Get the plugin system.
	pluginsystem := ambientApp.PluginSystem()

	// Create secure site for the core application and use "ambient" so it gets
	// full permissions.
	securestorage := ambient.NewSecureSite("ambient", log, pluginsystem, nil, nil, nil)

	// Enable plugins.
	for pluginName, trusted := range plugins.TrustedPlugins {
		if trusted {
			log.Info("enabling plugin: %v", pluginName)
			err := securestorage.EnablePlugin(pluginName, false)
			if err != nil {
				log.Error("", err.Error())
			}

			p, err := pluginsystem.Plugin(pluginName)
			if err != nil {
				log.Error("error with plugin (%v): %v", pluginName, err.Error())
				return
			}

			for _, request := range p.GrantRequests() {
				log.Info("%v - add grant: %v", pluginName, request.Grant)
				err := securestorage.SetNeighborPluginGrant(pluginName, request.Grant, true)
				if err != nil {
					log.Error("", err.Error())
				}
			}
		}
	}
}
