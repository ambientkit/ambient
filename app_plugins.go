package ambient

import (
	"fmt"
)

// LoadPlugins loads the plugins.
func (app *App) LoadPlugins() error {
	// Get the session manager from the plugins.
	var sess IAppSession
	for _, name := range app.pluginsystem.Names() {
		// Get the plugin.
		p, err := app.pluginsystem.Plugin(name)
		if err != nil {
			// This shouldn't happen because the names are based off the plugin list.
			return fmt.Errorf("ambient: could not find plugin (%v): %v", name, err.Error())
		}

		// Skip if the plugin isn't enabled.
		if !app.pluginsystem.Enabled(name) {
			continue
		}

		// Get the session manager.
		sm, err := p.SessionManager(app.log, app.sess)
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
		return fmt.Errorf("ambient: no session manager found")
	}

	// Set up the template injector.
	pi := NewPlugininjector(app.log, app.pluginsystem, sess)

	// Get the router from the plugins.
	var te IRender
	for _, name := range app.pluginsystem.Names() {
		// Skip if the plugin isn't found.
		plugin, err := app.pluginsystem.Plugin(name)
		if err != nil {
			// This shouldn't happen because the names are based off the plugin list.
			return fmt.Errorf("ambient: could not find plugin (%v): %v", name, err.Error())
		}

		// Skip if the plugin isn't enabled.
		if !app.pluginsystem.Enabled(name) {
			continue
		}

		// Get the router.
		tt, err := plugin.TemplateEngine(app.log, pi)
		if err != nil {
			return err
		} else if tt != nil {
			// Only set the router once.
			app.log.Info("ambient: using template engine from plugin: %v", name)
			te = tt
			break
		}
	}
	if te == nil {
		return fmt.Errorf("ambient: no template engine found")
	}

	// Get the router from the plugins.
	var mux IAppRouter
	for _, name := range app.pluginsystem.Names() {
		// Skip if the plugin isn't found.
		plugin, err := app.pluginsystem.Plugin(name)
		if err != nil {
			continue
		}

		// Skip if the plugin isn't enabled.
		if !app.pluginsystem.Enabled(name) {
			continue
		}

		// Get the router.
		rm, err := plugin.Router(app.log, te)
		if err != nil {
			return err
		} else if rm != nil {
			// Only set the router once.
			app.log.Info("ambient: using router (mux) from plugin: %v", name)
			mux = rm
			break
		}
	}
	if mux == nil {
		return fmt.Errorf("ambient: no template engine found")
	}

	// Create secure site for the core application and use "ambient" so it gets
	// full permissions.
	securesite := NewSecureSite("ambient", app.log, app.pluginsystem, sess, mux, te)

	// Load the plugin pages.
	err := securesite.LoadAllPluginPages()
	if err != nil {
		return err
	}

	// Enable the middleware from the plugins.
	app.handler = securesite.LoadAllPluginMiddleware(mux)

	return nil
}
