package ambientapp

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/internal/config"
	"github.com/ambientkit/ambient/internal/devconsole"
	"github.com/ambientkit/ambient/internal/grpcsystem"
	"github.com/ambientkit/ambient/internal/injector"
	"github.com/ambientkit/ambient/internal/pluginsafe"
	"github.com/ambientkit/ambient/internal/secureconfig"
	"github.com/ambientkit/ambient/pkg/envdetect"
	"github.com/ambientkit/ambient/pkg/requestuuid"
)

// App represents an Ambient app that supports plugins.
type App struct {
	log           ambient.AppLogger
	pluginsystem  ambient.PluginSystem
	grpcsystem    ambient.GRPCSystem
	sessionstorer ambient.SessionStorer
	mux           ambient.AppRouter
	renderer      ambient.Renderer
	sess          ambient.AppSession
	recorder      *pluginsafe.RouteRecorder
	securesite    *secureconfig.SecureSite

	debugTemplates  bool
	escapeTemplates bool
}

// NewAppLogger returns a logger from Ambient without all the other dependencies.
func NewAppLogger(appName string, appVersion string, logPlugin ambient.LoggingPlugin, logLevel ambient.LogLevel) (ambient.AppLogger, error) {
	// Set the time zone. Required for plugins that rely on timzone like MFA.
	tz := os.Getenv("AMB_TIMEZONE")
	if len(tz) > 0 {
		os.Setenv("TZ", tz)
	}

	// Get the logger from the plugin.
	log, err := loadLogger(appName, appVersion, logPlugin)
	if err != nil {
		return nil, err
	}

	// Set the initial log level.
	log.SetLogLevel(logLevel)

	return log, nil
}

// LoadLogger returns the logger.
func loadLogger(appName string, appVersion string, plugin ambient.LoggingPlugin) (ambient.AppLogger, error) {
	// Validate plugin name and version.
	err := ambient.Validate(plugin)
	if err != nil {
		return nil, err
	}

	// Get the logger from the plugins.
	log, err := plugin.Logger(appName, appVersion, nil)
	if err != nil {
		return nil, err
	} else if log == nil {
		return nil, fmt.Errorf("ambient: no logger found")
	} else {
		log.Info("ambient: using logger from plugin: %v", plugin.PluginName())
	}

	return log, nil
}

// NewApp returns a new Ambient app that supports plugins.
func NewApp(appName string, appVersion string, logPlugin ambient.LoggingPlugin,
	storagePluginGroup ambient.StoragePluginGroup, plugins *ambient.PluginLoader) (*App, ambient.AppLogger, error) {
	// Set up the logger first.
	log, err := NewAppLogger(appName, appVersion, logPlugin, ambient.EnvLogLevel())
	if err != nil {
		return nil, nil, err
	}

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
	pluginsystem, err := config.NewPluginSystem(log, storage, plugins)
	if err != nil {
		log.Fatal(err.Error())
	}

	grpcsystem := grpcsystem.New(log, pluginsystem)
	grpcsystem.ConnectAll()

	ambientApp := &App{
		log:             log,
		pluginsystem:    pluginsystem,
		grpcsystem:      grpcsystem,
		sessionstorer:   sessionstorer,
		escapeTemplates: true,
	}

	// Enable the trusted plugins.
	ambientApp.grantAccess()

	return ambientApp, log, nil
}

// PluginSystem returns the plugin system.
func (app *App) PluginSystem() ambient.PluginSystem {
	return app.pluginsystem
}

// LoadStorage returns the storage.
func loadStorage(log ambient.AppLogger, pluginGroup ambient.StoragePluginGroup) (*config.Storage, ambient.SessionStorer, error) {
	// Detect if storage plugin is missing.
	if pluginGroup.Storage == nil {
		return nil, nil, fmt.Errorf("ambient: storage plugin is missing")
	}

	plugin := pluginGroup.Storage

	// Validate plugin name and version.
	err := ambient.Validate(plugin)
	if err != nil {
		return nil, nil, err
	}

	// Define the storage managers.
	var ds ambient.DataStorer
	var ss ambient.SessionStorer

	// Get the storage manager from the plugins.
	pds, pss, err := plugin.Storage(log)
	if err != nil {
		log.Error(err.Error())
	} else if pds != nil && pss != nil {
		log.Info("ambient: using storage from first plugin: %v", plugin.PluginName())
		ds = pds
		ss = pss
	}
	if ds == nil || ss == nil {
		return nil, nil, fmt.Errorf("ambient: no storage manager found")
	}

	// Set up the data storage provider.
	storage, err := config.NewStorage(log, ds, pluginGroup.Encryption)
	if err != nil {
		return nil, nil, err
	}

	return storage, ss, err
}

// StopGRPCClients stops the gRPC plugins.
func (app *App) StopGRPCClients() {
	app.grpcsystem.Disconnect()
}

// Handler loads the plugins and returns the handler.
func (app *App) Handler() (http.Handler, error) {
	// Get the session manager from the plugins.
	if app.pluginsystem.SessionManager() != nil {
		sm, err := app.pluginsystem.SessionManager().SessionManager(app.log, app.sessionstorer)
		if err != nil {
			app.log.Error(err.Error())
		} else if sm != nil {
			// Only set the session manager once.
			app.log.Info("ambient: using session manager from plugin: %v", app.pluginsystem.SessionManager().PluginName())
			app.sess = sm
		}
	}
	if app.sess == nil {
		return nil, fmt.Errorf("ambient: no session manager found, ensure it is trusted")
	}

	// Set up the template injector.
	pi := injector.NewPlugininjector(app.log, app.pluginsystem, app.sess, app.debugTemplates, app.escapeTemplates)

	// Get the template engine.
	if app.pluginsystem.TemplateEngine() != nil {
		tt, err := app.pluginsystem.TemplateEngine().TemplateEngine(app.log, pi)
		if err != nil {
			return nil, err
		} else if tt != nil {
			// Only set the router once.
			app.log.Info("ambient: using template engine from plugin: %v", app.pluginsystem.TemplateEngine().PluginName())
			app.renderer = tt
		}
	}
	if app.renderer == nil {
		return nil, fmt.Errorf("ambient: no template engine found")
	}

	// Get the router.
	if app.pluginsystem.Router() != nil {
		rm, err := app.pluginsystem.Router().Router(app.log, app.renderer)
		if err != nil {
			return nil, err
		} else if rm != nil {
			// Only set the router once.
			app.log.Info("ambient: using router from plugin: %v", app.pluginsystem.Router().PluginName())
			app.mux = rm
		}
	}
	if app.mux == nil {
		return nil, fmt.Errorf("ambient: no router found")
	}

	app.recorder = pluginsafe.NewRouteRecorder(app.log, app.pluginsystem, app.mux)

	// Create secure site for the core app and use "ambient" so it gets
	// full permissions.
	var err error
	var handler http.Handler
	app.securesite, handler, err = secureconfig.NewSecureSite("ambient", app.log, app.pluginsystem, app.sess, app.mux, app.renderer, app.recorder, true)
	if err != nil {
		return nil, err
	}

	// Start monitoring with the ability to restart/reload plugin.
	app.grpcsystem.Monitor(app.securesite)

	// Start Dev Console if enabled via environment variable.
	if envdetect.DevConsoleEnabled() {
		// TODO: Should probably store in an object that can be edited by system.
		dc := devconsole.NewDevConsole(app.log, app.pluginsystem, app.pluginsystem.StorageManager(), app.securesite)
		dc.EnableDevConsole()
	}

	// Add a request UUID around all routes.
	return requestuuid.Middleware(handler), nil
}

// GrantAccess grants access to all trusted plugins.
func (app *App) grantAccess() {
	pluginsData := app.pluginsystem.PluginsData()

	// Enable trusted plugins.
	for _, pluginName := range app.pluginsystem.TrustedPluginNames() {
		// If plugin is not enabled, then enable.
		pluginInfo, found := pluginsData[pluginName]
		if !found {
			continue
		}

		if !pluginInfo.Enabled {
			app.log.Info("ambient: enabling trusted plugin: %v", pluginName)
			err := app.pluginsystem.SetEnabled(pluginName, true)
			if err != nil {
				app.log.Error(err.Error())
			}
		}

		p, err := app.pluginsystem.Plugin(pluginName)
		if err != nil {
			app.log.Error("error with plugin (%v): %v", pluginName, err.Error())
			return
		}

		for _, request := range p.GrantRequests() {
			// If plugin is not granted permission, then grant.
			if !app.pluginsystem.Granted(pluginName, request.Grant) {
				app.log.Info("ambient: for plugin (%v), adding grant: %v", pluginName, request.Grant)
				err = app.pluginsystem.SetGrant(pluginName, request.Grant)
				if err != nil {
					app.log.Error(err.Error())
				}
			}
		}
	}
}

// SetDebugTemplates sets the injector to enable verbose debug output in
// templates.
func (app *App) SetDebugTemplates(enable bool) {
	app.debugTemplates = enable
}

// SetLogLevel sets the log level.
func (app *App) SetLogLevel(level ambient.LogLevel) {
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
func (app *App) ListenAndServe(h http.Handler) error {
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
	return http.ListenAndServe(":"+port, h)
}

// handleExit will handle app shutdown when Ctrl+c is pressed.
func (app *App) handleExit() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		app.CleanUp()
		os.Exit(0)
	}()
}

// SecureSite returns the secure site configuration.
func (app *App) SecureSite() *secureconfig.SecureSite {
	return app.securesite
}

// CleanUp runs the final steps to ensure the server shutdown doesn't leave
// the app in a bad state.
func (app *App) CleanUp() {
	var err error
	app.log.Info("ambient: shutdown started")

	app.log.Info("ambient: stopping gRPC plugins")
	app.StopGRPCClients()

	// Load decrypted just in case the storage was decrypted by AMB.
	app.log.Info("ambient: loading storage")
	err = app.pluginsystem.StorageManager().LoadDecrypted()
	if err != nil {
		app.log.Error("ambient: could not load storage: %v", err.Error())
	}

	app.log.Info("ambient: saving storage")
	err = app.pluginsystem.StorageManager().Save()
	if err != nil {
		app.log.Error("ambient: could not save storage: %v", err.Error())
	}

	app.log.Info("ambient: shutdown done")
}
