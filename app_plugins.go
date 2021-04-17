package ambient

import (
	"fmt"
)

// SetPlugins sets the plugins.
func (app *App) SetPlugins(Plugins IPluginList) error {
	// Ensure there is at least the storage plugin.
	if len(Plugins) == 0 {
		return fmt.Errorf("ambient: no plugins found")
	}

	// Initialize the plugin system.
	ps, err := NewPluginSystem(app.log, app.storage, Plugins)
	if err != nil {
		return err
	}

	// Get the session manager from the plugins.
	var sess IAppSession
	for _, name := range ps.Names() {
		// Get the plugin.
		p, err := ps.Plugin(name)
		if err != nil {
			// This shouldn't happen because the names are based off the plugin list.
			return fmt.Errorf("ambient: could not find plugin (%v): %v", name, err.Error())
		}

		// Skip if the plugin isn't enabled.
		if !ps.Enabled(name) {
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
		return fmt.Errorf("ambient: no default session manager found")
	}

	// Set up the template injector.
	pi := NewPlugininjector(app.log, app.storage, sess, ps, Plugins)

	// Get the router from the plugins.
	var te IRender
	for _, name := range ps.Names() {
		// Skip if the plugin isn't found.
		plugin, err := ps.Plugin(name)
		if err != nil {
			// This shouldn't happen because the names are based off the plugin list.
			return fmt.Errorf("ambient: could not find plugin (%v): %v", name, err.Error())
		}

		// Skip if the plugin isn't enabled.
		if !ps.Enabled(name) {
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
		return fmt.Errorf("ambient: no default template engine found")
	}

	// Get the router from the plugins.
	var mux IAppRouter
	for _, name := range ps.Names() {
		// Skip if the plugin isn't found.
		plugin, err := ps.Plugin(name)
		if err != nil {
			continue
		}

		// Skip if the plugin isn't enabled.
		if !ps.Enabled(name) {
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
		return fmt.Errorf("ambient: no default template engine found")
	}

	// Create secure site for the core application. This should always be
	// ambient so it gets full permissions.
	securesite := NewSecureSite("ambient", app.log, app.storage, ps, sess, mux, te)

	// Load the plugin pages.
	err = securesite.LoadAllPluginPages()
	if err != nil {
		return err
	}

	// Enable the middleware from the plugins.
	app.handler = securesite.LoadAllPluginMiddleware(mux)

	return nil
}
