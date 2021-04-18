package ambient

import (
	"fmt"
	"net/http"
)

// Handler loads the plugins and returns the handler.
func (app *App) Handler() (http.Handler, error) {
	// Get the session manager from the plugins.
	var sess AppSession
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
			sess = sm
			break
		}
	}
	if sess == nil {
		return nil, fmt.Errorf("ambient: no session manager found")
	}

	// Set up the template injector.
	pi := NewPlugininjector(app.log, app.pluginsystem, sess, app.debugTemplates)

	// Get the template engine.
	var te Renderer
	if app.pluginsystem.templateEngine != nil {
		tt, err := app.pluginsystem.templateEngine.TemplateEngine(app.log, pi)
		if err != nil {
			return nil, err
		} else if tt != nil {
			// Only set the router once.
			app.log.Info("ambient: using template engine from plugin: %v", app.pluginsystem.templateEngine.PluginName())
			te = tt
		}
	}
	if te == nil {
		return nil, fmt.Errorf("ambient: no template engine found")
	}

	// Get the router.
	var mux AppRouter
	if app.pluginsystem.router != nil {
		rm, err := app.pluginsystem.router.Router(app.log, te)
		if err != nil {
			return nil, err
		} else if rm != nil {
			// Only set the router once.
			app.log.Info("ambient: using router (mux) from plugin: %v", app.pluginsystem.router.PluginName())
			mux = rm
		}
	}
	if mux == nil {
		return nil, fmt.Errorf("ambient: no router found")
	}

	// Create secure site for the core application and use "ambient" so it gets
	// full permissions.
	securesite := NewSecureSite("ambient", app.log, app.pluginsystem, sess, mux, te)

	// Load the plugin pages.
	err := securesite.LoadAllPluginPages()
	if err != nil {
		return nil, err
	}

	// Enable the middleware from the plugins.
	handler := securesite.LoadAllPluginMiddleware(mux)

	return handler, nil
}
