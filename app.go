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
	log           AppLogger
	pluginsystem  *PluginSystem
	sessionstorer SessionStorer

	debugTemplates bool
}

// NewApp returns a new Ambient app that supports plugins.
func NewApp(appName string, appVersion string, logPlugin LoggingPlugin, storagePlugin StoragePlugin, plugins *PluginLoader) (*App, error) {
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
	storage, sessionstorer, err := loadStorage(log, storagePlugin)
	if err != nil {
		return nil, err
	}

	// Initialize the plugin system.
	pluginsystem, err := NewPluginSystem(log, storage, plugins)
	if err != nil {
		log.Fatal("", err.Error())
	}

	return &App{
		log:           log,
		pluginsystem:  pluginsystem,
		sessionstorer: sessionstorer,
	}, nil
}

// Logger returns the logger.
func (app *App) Logger() AppLogger {
	return app.log
}

// PluginSystem returns the plugin system.
func (app *App) PluginSystem() *PluginSystem {
	return app.pluginsystem
}

// SetDebugTemplates sets the injector to enable verbose debug output in
// templates.
func (app *App) SetDebugTemplates(enable bool) {
	app.debugTemplates = enable
}

// ListenAndServe will start the web listener on port 8080 or will pull the
// environment variable from: PORT.
func (app *App) ListenAndServe(h http.Handler) {
	// Start the web server.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app.log.Info("ambient: web server listening on port: %v", port)
	app.log.Fatal("", http.ListenAndServe(":"+port, h))
}
