package app

import "github.com/josephspurrier/ambient"

// GrantAccess grants access to all trusted plugins.
func GrantAccess(log ambient.AppLogger, plugins *ambient.PluginLoader, ambientApp *ambient.App) {
	// Get the plugin system.
	pluginsystem := ambientApp.PluginSystem()

	// Create secure site for the core application and use "ambient" so it gets
	// full permissions.
	securestorage := ambient.NewSecureSite("ambient", log, pluginsystem, nil, nil, nil)

	// Enable plugins.
	for _, pluginName := range plugins.TrustedPluginNames() {
		trusted := plugins.TrustedPlugins[pluginName]
		if trusted {
			log.Info("enabling plugin: %v", pluginName)
			err := securestorage.EnablePlugin(pluginName, false)
			if err != nil {
				log.Error("", err.Error())
			}

			p, err := pluginsystem.Plugin(pluginName)
			if err != nil {
				log.Error("error with plugin (%v): %v", pluginName, err.Error())
				return
			}

			for _, request := range p.GrantRequests() {
				log.Info("%v - add grant: %v", pluginName, request.Grant)
				err := securestorage.SetNeighborPluginGrant(pluginName, request.Grant, true)
				if err != nil {
					log.Error("", err.Error())
				}
			}
		}
	}
}
