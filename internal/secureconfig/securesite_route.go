package secureconfig

import (
	"context"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/amberror"
)

// PluginNeighborRoutesList gets the routes for a neighbor plugin.
func (ss *SecureSite) PluginNeighborRoutesList(ctx context.Context, pluginName string) ([]ambient.Route, error) {
	if !ss.Authorized(ctx, ambient.GrantPluginNeighborRouteRead) {
		return nil, amberror.ErrAccessDenied
	}

	return ss.pluginsystem.Routes(pluginName), nil
}
