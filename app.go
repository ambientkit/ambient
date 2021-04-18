package ambient

import (
	"net/http"
	"os"
)

const (
	// Version is the Ambient version.
	Version = "1.0"
)

// App represents an Ambient app that supports plugins.
type App struct {
	log          IAppLogger
	storage      *Storage
	sess         SessionStorer
	handler      http.Handler
	pluginsystem *PluginSystem
}

// NewApp returns a new Ambient app that supports plugins.
func NewApp(appName string, appVersion string, logPlugin IPlugin, storagePlugin IPlugin, plugins []IPlugin) (*App, error) {
	// Get the logger from the plugin.
	log, err := loadLogger(appName, appVersion, logPlugin)
	if err != nil {
		return nil, err
	}

	// Set the log level.
	//log.SetLogLevel(LogLevelDebug)
	log.SetLogLevel(LogLevelInfo)
	//log.SetLogLevel(LogLevelError)
	//log.SetLogLevel(LogLevelFatal)

	// Get the storage manager.
	storage, sess, err := loadStorage(log, storagePlugin)
	if err != nil {
		return nil, err
	}

	// Initialize the plugin system.
	pluginsystem, err := NewPluginSystem(log, storage, plugins)
	if err != nil {
		log.Fatal("", err.Error())
	}

	return &App{
		log:          log,
		storage:      storage,
		sess:         sess,
		pluginsystem: pluginsystem,
	}, nil
}

// Logger returns the logger.
func (app *App) Logger() IAppLogger {
	return app.log
}

// Storage returns the storage manager.
// func (app *App) Storage() *Storage {
// 	return app.storage
// }

// PluginSystem returns the plugin system.
func (app *App) PluginSystem() *PluginSystem {
	return app.pluginsystem
}

// Mux returns the HTTP request multiplexer.
func (app *App) Mux() http.Handler {
	return app.handler
}

// ListenAndServe will start the web listener.
func (app *App) ListenAndServe() {
	// Start the web server.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app.log.Info("ambient: web server listening on port: %v", port)
	app.log.Fatal("", http.ListenAndServe(":"+port, app.Mux()))
}
