package ambient

// PluginLoader contains the plugins for the Ambient app.
type PluginLoader struct {
	Router         RouterPlugin
	TemplateEngine TemplateEnginePlugin
	SessionManager SessionManagerPlugin
	TrustedPlugins map[string]bool
	Plugins        []Plugin
	GRPCPlugins    []GRPCPlugin
	Middleware     []MiddlewarePlugin
}

// GRPCPlugin is a plugin over gRPC.
type GRPCPlugin struct {
	PluginName string
	PluginPath string
}
