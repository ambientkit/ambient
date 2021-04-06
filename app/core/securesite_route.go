package core

// ClearRoute clears out an old route.
func (ss *SecureSite) ClearRoute(method string, path string) error {
	grant := "router:clear"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	ss.mux.Clear(method, path)

	return nil
}

// ClearRoutePlugin clears out all the routes for a plugin.
func (ss *SecureSite) ClearRoutePlugin(pluginName string) error {
	grant := "router:clear"

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
