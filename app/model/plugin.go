package model

// PluginSettings -
type PluginSettings struct {
	Enabled bool     `json:"enabled"`
	Found   bool     `json:"found"`
	Fields  []string `json:"fields"`
}

// PluginFields -
type PluginFields struct {
	Fields map[string]string `json:"fields"`
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
