// Package hello provides a hello page for an Ambient app.
package hello

import (
	"embed"
	"time"

	"github.com/ambientkit/ambient"
)

//go:embed template/*.tmpl
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit
}

// New returns a new hello plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "hello"
}

// PluginVersion returns the plugin version.
func (p *Plugin) PluginVersion() string {
	return "1.0.0"
}

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *ambient.Toolkit) error {
	p.Toolkit = toolkit
	p.startBackgroundTask()
	return nil
}

// Disable the plugin background tasks.
func (p *Plugin) Disable() error {
	stopBackgroundTask()
	return nil
}

// Routes gets routes for the plugin.
func (p *Plugin) Routes() {
	p.Mux.Get("/dashboard/hello", p.index)
}

var (
	done   chan bool
	ticker *time.Ticker
)

func stopBackgroundTask() {
	done <- true
	ticker.Stop()
}

func (p *Plugin) startBackgroundTask() {
	ticker = time.NewTicker(500 * time.Millisecond)
	done = make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				p.Log.Info("", "Background task stopped")
				return
			case t := <-ticker.C:
				p.Log.Info("Tick at %v", t)
			}
		}
	}()

	p.Log.Info("", "Background task started")
}
