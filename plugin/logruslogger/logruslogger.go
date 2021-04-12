// Package logruslogger provides logruslogger functionality
// for an Ambient application.
package logruslogger

import (
	"github.com/josephspurrier/ambient/app/core"
	"github.com/sirupsen/logrus"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginBase
	*core.Toolkit
}

// New returns a new logruslogger plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &core.PluginBase{},
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
func (p *Plugin) Logger(appName string, appVersion string) (core.IAppLogger, error) {
	// Create the logger.
	log := NewLogger(appName, appVersion)
	//l.SetLevel(uint32(logrus.DebugLevel))
	log.SetLevel(uint32(logrus.InfoLevel))
	// l.SetLevel(logrus.ErrorLevel)
	// l.SetLevel(logrus.FatalLevel)

	return log, nil
}

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *core.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}
