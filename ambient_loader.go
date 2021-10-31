package ambient

import "sort"

// PluginLoader contains the plugins for the Ambient app.
type PluginLoader struct {
	Router         RouterPlugin
	TemplateEngine TemplateEnginePlugin
	TrustedPlugins map[string]bool
	Plugins        []Plugin
	Middleware     []MiddlewarePlugin
}

// TrustedPluginNames returns a list of sorted trusted names.
func (p *PluginLoader) TrustedPluginNames() []string {
	names := make([]string, 0)
	for name := range p.TrustedPlugins {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}
