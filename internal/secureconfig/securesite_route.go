package secureconfig

import (
	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/internal/amberror"
)

// PluginNeighborRoutesList gets the routes for a neighbor plugin.
func (ss *SecureSite) PluginNeighborRoutesList(pluginName string) ([]ambient.Route, error) {
	if !ss.Authorized(ambient.GrantPluginNeighborRouteRead) {
		return nil, amberror.ErrAccessDenied
	}

	return ss.pluginsystem.Routes(pluginName), nil
}
