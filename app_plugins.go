package ambient

import (
	"fmt"
	"net/http"

	"github.com/ambientkit/ambient/pkg/envdetect"
)

// Handler loads the plugins and returns the handler.
func (app *App) Handler() (http.Handler, error) {
	// Get the session manager from the plugins.
	if app.pluginsystem.sessionManager != nil {
		sm, err := app.pluginsystem.sessionManager.SessionManager(app.log, app.sessionstorer)
		if err != nil {
			app.log.Error(err.Error())
		} else if sm != nil {
			// Only set the session manager once.
			app.log.Info("ambient: using session manager from plugin: %v", app.pluginsystem.sessionManager.PluginName())
			app.sess = sm
		}
	}
	if app.sess == nil {
		return nil, fmt.Errorf("ambient: no session manager found, ensure it is trusted")
	}

	// Set up the template injector.
	pi := NewPlugininjector(app.log, app.pluginsystem, app.sess, app.debugTemplates, app.escapeTemplates)

	// Get the template engine.
	if app.pluginsystem.templateEngine != nil {
		tt, err := app.pluginsystem.templateEngine.TemplateEngine(app.log, pi)
		if err != nil {
			return nil, err
		} else if tt != nil {
			// Only set the router once.
			app.log.Info("ambient: using template engine from plugin: %v", app.pluginsystem.templateEngine.PluginName())
			app.renderer = tt
		}
	}
	if app.renderer == nil {
		return nil, fmt.Errorf("ambient: no template engine found")
	}

	// Get the router.
	if app.pluginsystem.router != nil {
		rm, err := app.pluginsystem.router.Router(app.log, app.renderer)
		if err != nil {
			return nil, err
		} else if rm != nil {
			// Only set the router once.
			app.log.Info("ambient: using router from plugin: %v", app.pluginsystem.router.PluginName())
			app.mux = rm
		}
	}
	if app.mux == nil {
		return nil, fmt.Errorf("ambient: no router found")
	}

	app.recorder = NewRouteRecorder(app.log, app.pluginsystem, app.mux)

	// Create secure site for the core app and use "ambient" so it gets
	// full permissions.
	securesite := NewSecureSite("ambient", app.log, app.pluginsystem, app.sess, app.mux, app.renderer, app.recorder)

	// Load the plugin pages.
	err := securesite.LoadAllPluginPages()
	if err != nil {
		return nil, err
	}

	// Enable the middleware from the plugins.
	handler := securesite.LoadAllPluginMiddleware()

	// Start Dev Console if enabled via environment variable.
	if envdetect.DevConsoleEnabled() {
		// TODO: Should probably store in an object that can be edited by system.
		dc := NewDevConsole(securesite)
		dc.EnableDevConsole()
	}

	return handler, nil
}

// GrantAccess grants access to all trusted plugins.
func (app *App) GrantAccess() {
	// Get the plugin system.
	pluginsystem := app.PluginSystem()

	// Create secure site for the core app and use "ambient" so it gets
	// full permissions.
	securestorage := NewSecureSite("ambient", app.log, pluginsystem, nil, nil, nil, nil)

	// Enable trusted plugins.
	for _, pluginName := range pluginsystem.TrustedPluginNames() {
		// If plugin is not enabled, then enable.
		if !securestorage.pluginsystem.Enabled(pluginName) {
			app.log.Info("ambient: enabling trusted plugin: %v", pluginName)
			err := securestorage.EnablePlugin(pluginName, false)
			if err != nil {
				app.log.Error(err.Error())
			}
		}

		p, err := pluginsystem.Plugin(pluginName)
		if err != nil {
			app.log.Error("error with plugin (%v): %v", pluginName, err.Error())
			return
		}

		for _, request := range p.GrantRequests() {
			// If plugin is not granted permission, then grant.
			if !securestorage.pluginsystem.Granted(pluginName, request.Grant) {
				app.log.Info("ambient: for plugin (%v), adding grant: %v", pluginName, request.Grant)
				err := securestorage.SetNeighborPluginGrant(pluginName, request.Grant, true)
				if err != nil {
					app.log.Error(err.Error())
				}
			}
		}
	}
}
