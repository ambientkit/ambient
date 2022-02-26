package ambient

// MockPlugin represents an Ambient plugin.
type MockPlugin struct {
	*PluginBase

	pluginName    string
	pluginVersion string

	MockGrants []GrantRequest
	MockRoutes func(p *PluginBase)
}

// NewMockPlugin returns a new mock plugin.
func NewMockPlugin(name string, version string) *MockPlugin {
	return &MockPlugin{
		PluginBase: &PluginBase{},

		pluginName:    name,
		pluginVersion: version,
	}
}

// PluginName returns the plugin name.
func (p *MockPlugin) PluginName() string {
	return p.pluginName
}

// PluginVersion returns the plugin version.
func (p *MockPlugin) PluginVersion() string {
	return p.pluginVersion
}

// GrantRequests returns a list of grants requested by the plugin.
func (p *MockPlugin) GrantRequests() []GrantRequest {
	return p.MockGrants
}

// Routes gets routes for the plugin.
func (p *MockPlugin) Routes() {
	p.MockRoutes(p.PluginBase)
}
