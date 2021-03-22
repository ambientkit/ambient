package model

// PluginSettings -
type PluginSettings struct {
	Enabled bool `json:"enabled"`
	Found   bool `json:"found"`
}

// PluginRoutes -
type PluginRoutes struct {
	Routes map[string][]Route
}

// Route -
type Route struct {
	Method string
	Path   string
}
