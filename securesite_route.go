package ambient

// PluginNeighborRoutesList gets the routes for a neighbor plugin.
func (ss *SecureSite) PluginNeighborRoutesList(pluginName string) ([]Route, error) {
	if !ss.Authorized(GrantPluginNeighborRouteRead) {
		return nil, ErrAccessDenied
	}

	routes, ok := ss.pluginsystem.routes[pluginName]
	if !ok {
		return nil, ErrNotFound
	}

	return routes, nil
}
