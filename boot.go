package ambient

import (
	"fmt"
	"net/http"
)

// Boot returns a router with the application ready to be started.
func Boot(Plugins IPluginList) (IAppLogger, http.Handler, error) {
	// Ensure there is at least the storage plugin.
	if len(Plugins) == 0 {
		return nil, nil, fmt.Errorf("boot: no plugins found")
	}

	// Set up the logger.
	log, err := LoadLogger("ambient", "1.0", Plugins[0])
	if err != nil {
		return nil, nil, err
	}

	// Set the log level.
	//log.SetLogLevel(LogLevelDebug)
	log.SetLogLevel(LogLevelInfo)
	//log.SetLogLevel(LogLevelError)
	//log.SetLogLevel(LogLevelFatal)

	// Set up the storage.
	storage, ss, err := LoadStorage(log, Plugins[1])
	if err != nil {
		return log, nil, err
	}

	// Initialize the plugin system.
	ps, err := NewPluginSystem(log, Plugins, storage)
	if err != nil {
		return log, nil, err
	}

	// Get the session manager from the plugins.
	var sess IAppSession
	for _, name := range ps.Names() {
		// Get the plugin.
		p, err := ps.Plugin(name)
		if err != nil {
			log.Error("boot: could not find plugin (%v): %v", name, err.Error())
			continue
		}

		// Skip if the plugin isn't enabled.
		if !ps.Enabled(name) {
			continue
		}

		// Get the session manager.
		sm, err := p.SessionManager(log, ss)
		if err != nil {
			log.Error("", err.Error())
		} else if sm != nil {
			// Only set the session manager once.
			log.Info("boot: using session manager from plugin: %v", name)
			sess = sm
			break
		}
	}

	// FIXME: Need to fail gracefully?
	if sess == nil {
		log.Fatal("boot: no default session manager found")
	}

	// Set up the template injector.
	pi := NewPlugininjector(log, storage, sess, ps, Plugins)

	// Get the router from the plugins.
	var te IRender
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
		tt, err := plugin.TemplateEngine(log, pi)
		if err != nil {
			log.Error("", err.Error())
		} else if tt != nil {
			// Only set the router once.
			log.Info("boot: using template engine from plugin: %v", name)
			te = tt
			break
		}
	}
	if te == nil {
		log.Fatal("boot: no default template engine found")
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
		rm, err := plugin.Router(log, te)
		if err != nil {
			log.Error("", err.Error())
		} else if rm != nil {
			// Only set the router once.
			log.Info("boot: using router (mux) from plugin: %v", name)
			mux = rm
			break
		}
	}
	if mux == nil {
		log.Fatal("boot: no default router found")
	}

	// Create secure site for the core application.
	securesite := NewSecureSite("ambient", log, storage, ps, sess, mux, te)

	// Load the plugin pages.
	err = securesite.LoadAllPluginPages()
	if err != nil {
		return log, nil, err
	}

	// Enable the middleware from the plugins.
	h := securesite.LoadAllPluginMiddleware(mux)

	return log, h, nil
}
