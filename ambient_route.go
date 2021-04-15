package ambient

// IRouteList represents a list of routes.
type IRouteList interface {
	Routes() []IRoute
}

// IRoute represents a route.
type IRoute interface {
	Method() string
	Path() string
}

// PluginRoutes holds a map of routes.
type PluginRoutes struct {
	Routes map[string][]Route
}

// Route is a route for a router.
type Route struct {
	Method string
	Path   string
}
