package ambient

// PluginLoader contains the plugins for the Ambient app.
type PluginLoader struct {
	Router         RouterPlugin
	TemplateEngine TemplateEnginePlugin
	SessionManager SessionManagerPlugin
	TrustedPlugins map[string]bool
	Plugins        []Plugin
	Middleware     []MiddlewarePlugin
}
