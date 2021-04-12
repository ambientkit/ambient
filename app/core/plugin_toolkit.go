package core

import (
	"embed"
	"html/template"
	"net/http"
)

// Toolkit provides utilities to plugins.
type Toolkit struct {
	Log      ILogger
	Mux      IRouter
	Render   IRender
	Security ISession
	Site     *SecureSite
}

// IAppRouter represents a router.
type IAppRouter interface {
	IRouter

	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Clear(method string, path string)
	SetNotFound(notFound http.Handler)
	SetServeHTTP(h func(w http.ResponseWriter, r *http.Request, status int, err error))
}

// IRender represents a template renderer.
type IRender interface {
	Page(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
	PageContent(w http.ResponseWriter, r *http.Request, content string, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
	Post(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
	PostContent(w http.ResponseWriter, r *http.Request, content string, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
	Error(w http.ResponseWriter, r *http.Request, content string, statusCode int, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
}

// IRouter represents a router.
type IRouter interface {
	Get(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Post(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Error(status int, w http.ResponseWriter, r *http.Request)
	Param(r *http.Request, name string) string
}

// ISession represents a user session.
type ISession interface {
	SetCSRF(r *http.Request) string
	CSRF(r *http.Request) bool
	UserAuthenticated(r *http.Request) (bool, error)
	SetUser(r *http.Request, username string)
	RememberMe(r *http.Request, remember bool)
	Logout(r *http.Request)
}

// ILogger representer the log service for the application.
type ILogger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
}
