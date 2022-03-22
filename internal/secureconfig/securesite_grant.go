package secureconfig

import (
	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/internal/amberror"
)

// NeighborPluginGrantList gets the grants requests for a neighbor plugin.
func (ss *SecureSite) NeighborPluginGrantList(pluginName string) ([]ambient.GrantRequest, error) {
	if !ss.Authorized(ambient.GrantPluginNeighborGrantRead) {
		return nil, amberror.ErrAccessDenied
	}

	plugin, err := ss.pluginsystem.Plugin(pluginName)
	if err != nil {
		return nil, amberror.ErrNotFound
	}

	return plugin.GrantRequests(), nil
}

// NeighborPluginGrants gets the map of granted permissions.
func (ss *SecureSite) NeighborPluginGrants(pluginName string) (map[ambient.Grant]bool, error) {
	if !ss.Authorized(ambient.GrantPluginNeighborGrantRead) {
		return nil, amberror.ErrAccessDenied
	}

	plugin, err := ss.pluginsystem.Plugin(pluginName)
	if err != nil {
		return nil, amberror.ErrNotFound
	}

	grants := make(map[ambient.Grant]bool)
	for _, grant := range plugin.GrantRequests() {
		grants[grant.Grant] = ss.pluginsystem.Granted(pluginName, grant.Grant)
	}

	return grants, nil
}

// NeighborPluginGranted returns true if the plugin has the grant.
func (ss *SecureSite) NeighborPluginGranted(pluginName string, grantName ambient.Grant) (bool, error) {
	if !ss.Authorized(ambient.GrantPluginNeighborGrantRead) {
		return false, amberror.ErrAccessDenied
	}

	return ss.pluginsystem.Authorized(pluginName, grantName), nil
}

// NeighborPluginRequestedGrant returns true if the plugin requests the grant.
// This shouldn't be used to determine if a plugin has been approved the grant.
func (ss *SecureSite) NeighborPluginRequestedGrant(pluginName string, grantName ambient.Grant) (bool, error) {
	if !ss.Authorized(ambient.GrantPluginNeighborGrantRead) {
		return false, amberror.ErrAccessDenied
	}

	plugin, err := ss.pluginsystem.Plugin(pluginName)
	if err != nil {
		return false, amberror.ErrNotFound
	}

	for _, grant := range plugin.GrantRequests() {
		if grant.Grant == grantName {
			return true, nil
		}
	}

	return false, nil
}

// SetNeighborPluginGrant sets a grant for a neighbor plugin.
func (ss *SecureSite) SetNeighborPluginGrant(pluginName string, grantName ambient.Grant, granted bool) error {
	if !ss.Authorized(ambient.GrantPluginNeighborGrantWrite) {
		return amberror.ErrAccessDenied
	}

	var err error
	if granted {
		// Get the list of grants and ensure the grant is requested by the
		// plugin or else deny it.
		var grants []ambient.GrantRequest
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
			return amberror.ErrGrantNotRequested
		}

		err = ss.pluginsystem.SetGrant(pluginName, grantName)
	} else {
		err = ss.pluginsystem.RemoveGrant(pluginName, grantName)
	}
	if err != nil {
		return err
	}

	return ss.pluginsystem.Save()
}
