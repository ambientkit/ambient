// Package zaplogger provides log functionality
// for an Ambient application.
package zaplogger

import (
	"github.com/josephspurrier/ambient/app/core"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginBase
	*core.Toolkit
}

// New returns a new zaplogger plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &core.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "zaplogger"
}

// PluginVersion returns the plugin version.
func (p *Plugin) PluginVersion() string {
	return "1.0.0"
}

// Logger returns a logger.
func (p *Plugin) Logger(appName string, appVersion string) (core.IAppLogger, error) {
	// Create the logger.
	log := NewLogger(appName, appVersion)

	return log, nil
}

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *core.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}
