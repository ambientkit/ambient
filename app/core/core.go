// Package core provides plugin functionality for the application.
package core

import (
	"net/http"

	"github.com/josephspurrier/ambient/app/lib/datastorage"
)

// App represents the app dependencies.
type App struct {
	Log     IAppLogger
	Plugins *PluginSystem
	Render  IAppRender
	Router  IAppRouter
	Sess    ISession
	Storage *datastorage.Storage
}

// NewApp returns a new application.
func NewApp(logger IAppLogger,
	plugins *PluginSystem,
	render IAppRender,
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

// IAppRender represents a renderer.
type IAppRender interface {
	IRender

	Dashboard(w http.ResponseWriter, r *http.Request, partialTemplate string, vars map[string]interface{}) (status int, err error)
	Page(w http.ResponseWriter, r *http.Request, partialTemplate string, vars map[string]interface{}) (status int, err error)
	Post(w http.ResponseWriter, r *http.Request, postContent string, vars map[string]interface{}) (status int, err error)
	Bloglist(w http.ResponseWriter, r *http.Request, partialTemplate string, vars map[string]interface{}) (status int, err error)
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
