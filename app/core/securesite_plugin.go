package core

// Plugins returns the plugin list.
func (ss *SecureSite) Plugins() (map[string]PluginData, error) {
	if !ss.Authorized(GrantSitePluginRead) {
		return nil, ErrAccessDenied
	}

	return ss.storage.Site.PluginStorage, nil
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
func (ss *SecureSite) EnablePlugin(name string) error {
	if !ss.Authorized(GrantSitePluginEnable) {
		return ErrAccessDenied
	}

	plugin, ok := ss.storage.Site.PluginStorage[name]
	if !ok {
		return ErrNotFound
	}

	plugin.Enabled = true
	ss.storage.Site.PluginStorage[name] = plugin

	return ss.storage.Save()
}

// DisablePlugin disables a plugin.
func (ss *SecureSite) DisablePlugin(name string) error {
	if !ss.Authorized(GrantSitePluginDisable) {
		return ErrAccessDenied
	}

	plugin, ok := ss.storage.Site.PluginStorage[name]
	if !ok {
		return ErrNotFound
	}

	plugin.Enabled = false
	ss.storage.Site.PluginStorage[name] = plugin

	return ss.storage.Save()
}
