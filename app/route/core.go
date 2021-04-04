// Package route provides the handlers for the application.
package route

import (
	"net/http"

	"github.com/josephspurrier/ambient/app/core"
)

// IRouter represents a router.
type IRouter interface {
	SetServeHTTP(csh func(w http.ResponseWriter, r *http.Request, status int, err error))
	SetNotFound(notFound http.Handler)
}

// ITemplateEngine represents a template engine.
type ITemplateEngine interface {
	Error(w http.ResponseWriter, r *http.Request, partialTemplate string, vars map[string]interface{}) (status int, err error)
}

// Register all routes.
func Register(c *core.App) {
	// Register routes.
	registerHomePost(&HomePost{c})
	registerAdminPost(&AdminPost{c})
}
