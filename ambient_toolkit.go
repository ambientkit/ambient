package ambient

import (
	"embed"
	"html/template"
	"net/http"
)

// Toolkit provides utilities to plugins.
type Toolkit struct {
	Log    ILogger
	Mux    IRouter
	Render IRender
	Site   *SecureSite
}

// IRender represents a template renderer.
type IRender interface {
	Page(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
	PageContent(w http.ResponseWriter, r *http.Request, content string, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
	Post(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
	PostContent(w http.ResponseWriter, r *http.Request, content string, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
	Error(w http.ResponseWriter, r *http.Request, content string, statusCode int, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
}

// IAppRouter represents a router.
type IAppRouter interface {
	IRouter

	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Clear(method string, path string)
	SetNotFound(notFound http.Handler)
	SetServeHTTP(h func(w http.ResponseWriter, r *http.Request, status int, err error))
}

// IRouter represents a router.
type IRouter interface {
	Get(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Post(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Patch(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Put(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Delete(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Head(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Options(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Error(status int, w http.ResponseWriter, r *http.Request)
	Param(r *http.Request, name string) string
	Wrap(handler http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) (status int, err error)
}

// IAppSession represents a user session.
type IAppSession interface {
	UserAuthenticated(r *http.Request) (bool, error)
	Login(r *http.Request, username string)
	Logout(r *http.Request)
	Persist(r *http.Request, persist bool)
	SetCSRF(r *http.Request) string
	CSRF(r *http.Request) bool
}
