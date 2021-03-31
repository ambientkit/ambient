// Package hello provides a hello page for an Ambient application.
package hello

import (
	"embed"
	"fmt"
	"time"

	"github.com/josephspurrier/ambient/app/core"
)

//go:embed *
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginMeta
	*core.Toolkit
}

// New returns a new hello plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "hello",
			Version:    "1.0.0",
			AppVersion: "1.0.0",
		},
	}
}

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *core.Toolkit) error {
	p.Toolkit = toolkit
	startBackgroundTask()
	return nil
}

// Disable the plugin background tasks.
func (p *Plugin) Disable() error {
	stopBackgroundTask()
	return nil
}

// Routes gets routes for the plugin.
func (p *Plugin) Routes() {
	p.Router.Get("/dashboard/hello", p.index)
}

var (
	done   chan bool
	ticker *time.Ticker
)

func stopBackgroundTask() {
	done <- true
	ticker.Stop()
}

func startBackgroundTask() {
	ticker = time.NewTicker(500 * time.Millisecond)
	done = make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("Background task stopped")
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			}
		}
	}()

	fmt.Println("Background task started")
}
