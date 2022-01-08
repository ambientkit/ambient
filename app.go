package ambient

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/josephspurrier/ambient/lib/envdetect"
	"github.com/josephspurrier/ambient/plugin/router/awayrouter/router"
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
	mux           AppRouter
	renderer      Renderer
	sess          AppSession

	debugTemplates  bool
	escapeTemplates bool
}

// NewApp returns a new Ambient app that supports plugins.
func NewApp(appName string, appVersion string, logPlugin LoggingPlugin, storagePluginGroup StoragePluginGroup, plugins *PluginLoader) (*App, AppLogger, error) {
	// Set the time zone. Required for plugins that rely on timzone like MFA.
	tz := os.Getenv("AMB_TIMEZONE")
	if len(tz) > 0 {
		os.Setenv("TZ", tz)
	}

	// Get the logger from the plugin.
	log, err := loadLogger(appName, appVersion, logPlugin)
	if err != nil {
		return nil, nil, err
	}

	// Set the default log level.
	log.SetLogLevel(LogLevelInfo)

	// Get the storage manager.
	storage, sessionstorer, err := loadStorage(log, storagePluginGroup)
	if err != nil {
		return nil, log, err
	}

	// Implicitly trust session manager so the middleware will work properly.
	if plugins.SessionManager != nil {
		plugins.TrustedPlugins[plugins.SessionManager.PluginName()] = true
	}

	// Initialize the plugin system.
	pluginsystem, err := NewPluginSystem(log, storage, plugins)
	if err != nil {
		log.Fatal("", err.Error())
	}

	ambientApp := &App{
		log:             log,
		pluginsystem:    pluginsystem,
		sessionstorer:   sessionstorer,
		escapeTemplates: true,
	}

	// Enable the trusted plugins.
	ambientApp.GrantAccess()

	// Start local dev server for configuration.
	// TODO: Change so amb is a flag instead of a hard-coded name.
	if envdetect.RunningLocalDev() && appName != "amb" {
		// TODO: Make the port dynamic.
		devPort := "8081"
		log.Info("ambient: dev console started on: %v", devPort)

		go func() {
			mux := router.New()
			mux.Post("/storage/encrypt", func(w http.ResponseWriter, r *http.Request) (int, error) {
				log.Info("ambient: dev console - site.bin encrypted")
				err = storage.LoadDecrypted()
				if err != nil {
					return http.StatusInternalServerError, err
				}
				err = storage.Save()
				if err != nil {
					return http.StatusInternalServerError, err
				}

				return http.StatusOK, nil
			})

			mux.Post("/storage/decrypt", func(w http.ResponseWriter, r *http.Request) (int, error) {
				log.Info("ambient: dev console - site.bin decrypted")
				err = storage.SaveDecrypted()
				if err != nil {
					return http.StatusInternalServerError, err
				}

				return http.StatusOK, nil
			})

			err = http.ListenAndServe(":"+devPort, mux)
			if err != nil {
				log.Error("ambient: dev config server cannot start: %v", err.Error())
			}
		}()
	}

	return ambientApp, log, nil
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

// SetLogLevel sets the log level.
func (app *App) SetLogLevel(level LogLevel) {
	app.log.SetLogLevel(level)
}

// SetEscapeTemplates sets the injector to disable (enabled by default) escaping
// templates.
func (app *App) SetEscapeTemplates(enable bool) {
	app.escapeTemplates = enable
}

// ListenAndServe will start the web listener on port 8080 or will pull the
// environment variable from:
// PORT (GCP), _LAMBDA_SERVER_PORT (AWS), or FUNCTIONS_CUSTOMHANDLER_PORT (Azure).
func (app *App) ListenAndServe(h http.Handler) {
	// Start the web server. Google Cloud uses standardized PORT env variable.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get the AWS Lambda port if it's set.
	awsPort, exists := os.LookupEnv("_LAMBDA_SERVER_PORT")
	if exists {
		port = awsPort
	}

	// Get the Microsoft Azure Functions port if it's set.
	azurePort, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if exists {
		port = azurePort
	}

	app.handleExit()

	app.log.Info("ambient: web server listening on port: %v", port)
	app.log.Fatal("", http.ListenAndServe(":"+port, h))
}

// handleExit will handle app shutdown when Ctrl+c is pressed.
func (app *App) handleExit() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		app.cleanup()
		os.Exit(0)
	}()
}

// cleanup runs the final steps to ensure the server shutdown doesn't leave
// the application in a bad state.
func (app *App) cleanup() {
	var err error
	app.log.Info("ambient: shutdown started")

	// Load decrypted just in case the storage was decrypted by AMB.
	app.log.Info("ambient: loading storage")
	err = app.pluginsystem.storage.LoadDecrypted()
	if err != nil {
		app.log.Error("ambient: could not load storage: %v", err.Error())
	}

	app.log.Info("ambient: saving storage")
	err = app.pluginsystem.storage.Save()
	if err != nil {
		app.log.Error("ambient: could not save storage: %v", err.Error())
	}

	app.log.Info("ambient: shutdown done")
}
