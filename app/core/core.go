// Package core provides plugin functionality for the application.
package core

import (
	"github.com/josephspurrier/ambient/app/lib/datastorage"
	"github.com/josephspurrier/ambient/app/lib/htmltemplate"
	"github.com/josephspurrier/ambient/app/lib/logger"
	"github.com/josephspurrier/ambient/app/lib/router"
	"github.com/josephspurrier/ambient/app/lib/websession"
)

// NewApp returns a new application.
func NewApp(logger *logger.Logger,
	plugins *PluginSystem,
	render *htmltemplate.Engine,
	mux *router.Mux,
	sess *websession.Session,
	storage *datastorage.Storage) *App {
	return &App{
		Log:     logger,
		Plugins: plugins,
		Render:  render,
		Router:  mux,
		Sess:    sess,
		Storage: storage,
	}
}

// App -
type App struct {
	Log     *logger.Logger
	Plugins *PluginSystem
	Render  *htmltemplate.Engine
	Router  *router.Mux
	Sess    *websession.Session
	Storage *datastorage.Storage
}
