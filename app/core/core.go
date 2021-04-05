// Package core provides plugin functionality for the application.
package core

import (
	"net/http"
)

// App represents the app dependencies.
type App struct {
	Log     IAppLogger
	Plugins *PluginSystem
	Render  IRender
	Router  IAppRouter
	Sess    ISession
	Storage *Storage
}

// NewApp returns a new application.
func NewApp(logger IAppLogger,
	plugins *PluginSystem,
	render IRender,
	mux IAppRouter,
	sess ISession,
	storage *Storage) *App {
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
	SetNotFound(notFound http.Handler)
	SetServeHTTP(h func(w http.ResponseWriter, r *http.Request, status int, err error))
}

// IAppLogger represents a logger.
type IAppLogger interface {
	ILogger

	SetLevel(level uint32)
}
