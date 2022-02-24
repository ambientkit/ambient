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

// ClearRoute clears out an old route.
func (ss *SecureSite) ClearRoute(method string, path string) error {
	if !ss.Authorized(GrantRouterRouteClear) {
		return ErrAccessDenied
	}

	ss.mux.Clear(method, path)

	return nil
}

// ClearAllRoutesForPlugin clears out all the routes for a plugin.
func (ss *SecureSite) ClearAllRoutesForPlugin(pluginName string) error {
	if !ss.Authorized(GrantRouterNeighborRouteClear) {
		return ErrAccessDenied
	}

	routes, ok := ss.pluginsystem.routes[pluginName]
	if !ok {
		return ErrNotFound
	}

	for _, v := range routes {
		ss.mux.Clear(v.Method, v.Path)
	}

	return nil
}
