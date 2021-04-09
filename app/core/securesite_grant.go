package core

// NeighborPluginGrantList gets the grants requests for a neighbor plugin.
func (ss *SecureSite) NeighborPluginGrantList(pluginName string) ([]Grant, error) {
	if !ss.Authorized(GrantPluginNeighborgrantRead) {
		return nil, ErrAccessDenied
	}

	plugin, err := ss.pluginsystem.Plugin(ss.pluginName)
	if err != nil {
		return nil, ErrNotFound
	}

	return plugin.Grants(), nil
}

// NeighborPluginGrants gets the map of granted permissions.
func (ss *SecureSite) NeighborPluginGrants(pluginName string) (map[Grant]bool, error) {
	if !ss.Authorized(GrantPluginNeighborgrantRead) {
		return nil, ErrAccessDenied
	}

	plugin, err := ss.pluginsystem.Plugin(ss.pluginName)
	if err != nil {
		return nil, ErrNotFound
	}

	grants := make(map[Grant]bool)
	for _, grant := range plugin.Grants() {
		grants[grant] = ss.pluginsystem.Granted(ss.pluginName, grant)
	}

	return grants, nil
}

// SetNeighborPluginGrant sets a grant for a neighbor plugin.
func (ss *SecureSite) SetNeighborPluginGrant(pluginName string, grantName Grant, granted bool) error {
	if !ss.Authorized(GrantPluginNeighborgrantWrite) {
		return ErrAccessDenied
	}

	var err error
	if granted {
		err = ss.pluginsystem.SetGrant(pluginName, grantName)
	} else {
		err = ss.pluginsystem.RemoveGrant(pluginName, grantName)
	}
	if err != nil {
		return err
	}

	return ss.storage.Save()
}
