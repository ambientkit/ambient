// Package core provides plugin functionality for the application.
package core

import (
	"net/http"

	"github.com/josephspurrier/ambient/app/lib/datastorage"
	"github.com/josephspurrier/ambient/app/lib/htmltemplate"
)

// App represents the app dependencies.
type App struct {
	Log     IAppLogger
	Plugins *PluginSystem
	Render  *htmltemplate.Engine
	Router  IAppRouter
	Sess    ISession
	Storage *datastorage.Storage
}

// NewApp returns a new application.
func NewApp(logger IAppLogger,
	plugins *PluginSystem,
	render *htmltemplate.Engine,
	mux IAppRouter,
	sess ISession,
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

// IAppRouter represents a router.
type IAppRouter interface {
	IRouter

	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Clear(method string, path string)
}

// IAppLogger represents a logger.
type IAppLogger interface {
	ILogger

	SetLevel(level uint32)
}
