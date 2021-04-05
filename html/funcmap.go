// Package html is the base templates and functions for rendering the application.
package html

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/josephspurrier/ambient/app/core"
)

//go:embed layout/page.tmpl
var templates embed.FS

// TemplateManager represents an object that returns templates and a FuncMap.
type TemplateManager struct {
	storage *core.Storage
	sess    ISession
}

// ISession represents a user session.
type ISession interface {
	UserAuthenticated(r *http.Request) (bool, error)
}

// NewTemplateManager returns a TemplateManager.
func NewTemplateManager(storage *core.Storage, sess ISession) *TemplateManager {
	return &TemplateManager{
		storage: storage,
		sess:    sess,
	}
}

// Templates returns the embedded templates.
func (f *TemplateManager) Templates() *embed.FS {
	return &templates
}

// FuncMap returns a map of template functions that can be used in templates.
func (f *TemplateManager) FuncMap(r *http.Request) template.FuncMap {
	fm := make(template.FuncMap)
	fm["Authenticated"] = func() bool {
		// If user is not authenticated, don't allow them to access the page.
		loggedIn, err := f.sess.UserAuthenticated(r)
		if err != nil {
			// TODO: Need to switch over to the logger.
			log.Println(err)
		}
		return loggedIn
	}

	return fm
}
