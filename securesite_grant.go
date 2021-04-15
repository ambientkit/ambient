package ambient

// NeighborPluginGrantList gets the grants requests for a neighbor plugin.
func (ss *SecureSite) NeighborPluginGrantList(pluginName string) ([]GrantRequest, error) {
	if !ss.Authorized(GrantPluginNeighborGrantRead) {
		return nil, ErrAccessDenied
	}

	plugin, err := ss.pluginsystem.Plugin(pluginName)
	if err != nil {
		return nil, ErrNotFound
	}

	return plugin.GrantRequests(), nil
}

// NeighborPluginGrants gets the map of granted permissions.
func (ss *SecureSite) NeighborPluginGrants(pluginName string) (map[Grant]bool, error) {
	if !ss.Authorized(GrantPluginNeighborGrantRead) {
		return nil, ErrAccessDenied
	}

	plugin, err := ss.pluginsystem.Plugin(pluginName)
	if err != nil {
		return nil, ErrNotFound
	}

	grants := make(map[Grant]bool)
	for _, grant := range plugin.GrantRequests() {
		grants[grant.Grant] = Granted(ss.log, ss.storage, pluginName, grant.Grant)
	}

	return grants, nil
}

// SetNeighborPluginGrant sets a grant for a neighbor plugin.
func (ss *SecureSite) SetNeighborPluginGrant(pluginName string, grantName Grant, granted bool) error {
	if !ss.Authorized(GrantPluginNeighborGrantWrite) {
		return ErrAccessDenied
	}

	var err error
	if granted {
		// Get the list of grants and ensure the grant is requested by the
		// plugin or else deny it.
		var grants []GrantRequest
		grants, err = ss.NeighborPluginGrantList(pluginName)
		if err != nil {
			return err
		}

		found := false
		for _, grant := range grants {
			if grant.Grant == grantName {
				found = true
				break
			}
		}

		if !found {
			ss.log.Debug("grant to enable on plugin %v was not requested by the plugin: %v", pluginName, grantName)
			return ErrGrantNotRequested
		}

		err = ss.pluginsystem.SetGrant(pluginName, grantName)
	} else {
		err = ss.pluginsystem.RemoveGrant(pluginName, grantName)
	}
	if err != nil {
		return err
	}

	return ss.storage.Save()
}
