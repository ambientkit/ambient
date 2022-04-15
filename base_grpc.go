package ambient

import "context"

// GRPCPluginBase represents a base gRPC plugin that works with Ambient.
type GRPCPluginBase struct {
	PluginBase

	pluginName string
	pluginPath string
}

// NewGRPCPlugin returns gRPC plugin base.
func NewGRPCPlugin(pluginName string, pluginPath string) *GRPCPluginBase {
	return &GRPCPluginBase{
		pluginName: pluginName,
		pluginPath: pluginPath,
	}
}

// PluginName returns the gRPC plugin name.
func (p *GRPCPluginBase) PluginName(context.Context) string {
	return p.pluginName
}

// PluginVersion returns the gRPc text.
func (p *GRPCPluginBase) PluginVersion(context.Context) string {
	return "gRPC"
}

// PluginPath returns the gRPC plugin path.
func (p *GRPCPluginBase) PluginPath() string {
	return p.pluginPath
}
