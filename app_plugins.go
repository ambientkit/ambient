package ambient

import (
	"fmt"
	"net/http"
)

// Handler loads the plugins and returns the handler.
func (app *App) Handler() (http.Handler, error) {
	// Get the session manager from the plugins.
	for _, name := range app.pluginsystem.Names() {
		// Get the plugin.
		p, err := app.pluginsystem.Plugin(name)
		if err != nil {
			// This shouldn't happen because the names are based off the plugin list.
			return nil, fmt.Errorf("ambient: could not find plugin (%v): %v", name, err.Error())
		}

		// Skip if the plugin isn't enabled.
		if !app.pluginsystem.Enabled(name) {
			continue
		}

		// Get the session manager.
		sm, err := p.SessionManager(app.log, app.sessionstorer)
		if err != nil {
			app.log.Error("", err.Error())
		} else if sm != nil {
			// Only set the session manager once.
			app.log.Info("ambient: using session manager from plugin: %v", name)
			app.sess = sm
			break
		}
	}
	if app.sess == nil {
		return nil, fmt.Errorf("ambient: no session manager found")
	}

	// Set up the template injector.
	pi := NewPlugininjector(app.log, app.pluginsystem, app.sess, app.debugTemplates)

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
			app.log.Info("ambient: using router (mux) from plugin: %v", app.pluginsystem.router.PluginName())
			app.mux = rm
		}
	}
	if app.mux == nil {
		return nil, fmt.Errorf("ambient: no router found")
	}

	// Create secure site for the core application and use "ambient" so it gets
	// full permissions.
	securesite := NewSecureSite("ambient", app.log, app.pluginsystem, app.sess, app.mux, app.renderer)

	// Load the plugin pages.
	err := securesite.LoadAllPluginPages()
	if err != nil {
		return nil, err
	}

	// Enable the middleware from the plugins.
	handler := securesite.LoadAllPluginMiddleware()

	return handler, nil
}

// Toolkit returns a toolkit for use with plugins externally.
func (app *App) Toolkit(pluginName string) *Toolkit {
	toolkit := &Toolkit{
		Mux:    NewRecorder(pluginName, app.log, app.pluginsystem, app.mux),
		Render: NewRenderer(app.renderer),
		Site:   NewSecureSite(pluginName, app.log, app.pluginsystem, app.sess, app.mux, app.renderer),
		Log:    NewPluginLogger(app.log),
	}

	return toolkit
}
