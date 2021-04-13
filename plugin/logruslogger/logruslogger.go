// Package logruslogger provides log functionality
// for an Ambient application.
package logruslogger

import "github.com/josephspurrier/ambient"

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit
}

// New returns a new logruslogger plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "logruslogger"
}

// PluginVersion returns the plugin version.
func (p *Plugin) PluginVersion() string {
	return "1.0.0"
}

// Logger returns a logger.
func (p *Plugin) Logger(appName string, appVersion string) (ambient.IAppLogger, error) {
	// Create the logger.
	log := NewLogger(appName, appVersion)

	return log, nil
}

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *ambient.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}
