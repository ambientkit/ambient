// Package core provides plugin functionality for the application.
package core

import (
	"net/http"

	"github.com/josephspurrier/ambient/app/lib/datastorage"
	"github.com/josephspurrier/ambient/app/lib/htmltemplate"
	"github.com/josephspurrier/ambient/app/lib/logger"
)

// App -
type App struct {
	Log     *logger.Logger
	Plugins *PluginSystem
	Render  *htmltemplate.Engine
	Router  IAppRouter
	Sess    ISession
	Storage *datastorage.Storage
}

// NewApp returns a new application.
func NewApp(logger *logger.Logger,
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

// IAppRouter -
type IAppRouter interface {
	Get(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Post(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Error(status int, w http.ResponseWriter, r *http.Request)
	Param(r *http.Request, param string) string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Clear(method string, path string)
}
