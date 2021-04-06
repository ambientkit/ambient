package core

// ClearRoute clears out an old route.
func (ss *SecureSite) ClearRoute(method string, path string) error {
	grant := "router.route:clear"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	ss.mux.Clear(method, path)

	return nil
}

// ClearAllRoutesForPlugin clears out all the routes for a plugin.
func (ss *SecureSite) ClearAllRoutesForPlugin(pluginName string) error {
	grant := "router.neighborroute:clear"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	routes, ok := ss.storage.PluginRoutes.Routes[pluginName]
	if !ok {
		return ErrNotFound
	}

	for _, v := range routes {
		ss.mux.Clear(v.Method, v.Path)
	}

	return nil
}
