package core

import (
	"github.com/josephspurrier/ambient/app/lib/routerrecorder"
)

// Plugins returns the plugin list.
func (ss *SecureSite) Plugins() (map[string]PluginData, error) {
	if !ss.Authorized(GrantSitePluginRead) {
		return nil, ErrAccessDenied
	}

	return ss.storage.Site.PluginStorage, nil
}

// PluginNames returns the list of plugin name.
func (ss *SecureSite) PluginNames() ([]string, error) {
	if !ss.Authorized(GrantSitePluginRead) {
		return nil, ErrAccessDenied
	}

	return ss.pluginsystem.names, nil
}

// DeletePlugin deletes a plugin.
func (ss *SecureSite) DeletePlugin(name string) error {
	if !ss.Authorized(GrantSitePluginDelete) {
		return ErrAccessDenied
	}

	delete(ss.storage.Site.PluginStorage, name)

	return ss.storage.Save()
}

// EnablePlugin enables a plugin.
func (ss *SecureSite) EnablePlugin(pluginName string) error {
	if !ss.Authorized(GrantSitePluginEnable) {
		return ErrAccessDenied
	}

	// Load the plugin and routes.
	err := ss.loadSinglePlugin(pluginName)
	if err != nil {
		return err
	}

	pluginData, ok := ss.storage.Site.PluginStorage[pluginName]
	if !ok {
		return ErrNotFound
	}

	pluginData.Enabled = true
	ss.storage.Site.PluginStorage[pluginName] = pluginData

	return ss.storage.Save()
}

func (ss *SecureSite) loadSinglePlugin(name string) error {
	plugins, err := ss.Plugins()
	if err != nil {
		return err
	}

	save := ss.loadSinglePluginPages(name, plugins)
	if save {
		err := ss.storage.Save()
		if err != nil {
			return err
		}
	}

	return nil
}

func (ss *SecureSite) loadSinglePluginPages(name string, pluginsData map[string]PluginData) bool {
	v, err := ss.pluginsystem.Plugin(name)
	if err != nil {
		ss.log.Error("plugin load: problem loading plugin %v: %v", name, err.Error())
		return false
	}

	recorder := routerrecorder.NewRecorder(ss.mux)

	toolkit := &Toolkit{
		Mux:      recorder,
		Render:   ss.render, // FIXME: Should probably remove this and create a new struct so it's more secure. A plugin could use a type conversion.
		Security: ss.sess,
		Site:     NewSecureSite(name, ss.log, ss.storage, ss.sess, ss.mux, ss.render, ss.pluginsystem),
		Log:      ss.log,
	}

	// Enable the plugin and pass in the toolkit.
	err = v.Enable(toolkit)
	if err != nil {
		ss.log.Error("plugin load: problem enabling plugin %v: %v", name, err.Error())
		return false
	}

	// Load the routes.
	v.Routes()

	// Load the assets.
	assets, files := v.Assets()
	if files == nil {
		// Save the plugin routes so they can be removed if disabled.
		saveRoutesForPlugin(name, recorder, ss.storage)
		return false
	}

	// Handle embedded assets.
	err = embeddedAssets(recorder, ss.sess, name, assets, files)
	if err != nil {
		ss.log.Error("plugin load: problem loading assets for plugin %v: %v", name, err.Error())
	}

	// Save the plugin routes so they can be removed if disabled.
	saveRoutesForPlugin(name, recorder, ss.storage)

	return false
}

// DisablePlugin disables a plugin.
func (ss *SecureSite) DisablePlugin(pluginName string) error {
	if !ss.Authorized(GrantSitePluginDisable) {
		return ErrAccessDenied
	}

	// Get the plugin.
	plugin, ok := ss.pluginsystem.plugins[pluginName]
	if !ok {
		return ErrNotFound
	}

	// Disable the plugin.
	err := plugin.Disable()
	if err != nil {
		return err
	}

	// Get the routes for the plugin.
	routes, ok := ss.storage.PluginRoutes.Routes[pluginName]
	if !ok {
		return ErrNotFound
	}

	// Clear each route.
	for _, v := range routes {
		ss.mux.Clear(v.Method, v.Path)
	}

	// Get the plugin data.
	pluginData, ok := ss.storage.Site.PluginStorage[pluginName]
	if !ok {
		return ErrNotFound
	}

	// Disable the plugin.
	pluginData.Enabled = false
	ss.storage.Site.PluginStorage[pluginName] = pluginData

	return ss.storage.Save()
}
