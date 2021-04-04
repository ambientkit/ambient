// Package html is the base templates and functions for rendering the application.
package html

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/josephspurrier/ambient/app/core"
	"github.com/josephspurrier/ambient/app/model"
)

//go:embed *
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
	fm["Stamp"] = func(t time.Time) string {
		return t.Format("2006-01-02")
	}
	fm["StampFriendly"] = func(t time.Time) string {
		return t.Format("02 Jan, 2006")
	}
	fm["PublishedPages"] = func() []model.Post {
		return f.storage.Site.PublishedPages()
	}
	fm["SiteURL"] = func() string {
		return f.storage.Site.SiteURL()
	}
	fm["SiteTitle"] = func() string {
		return f.storage.Site.SiteTitle()
	}
	fm["SiteSubtitle"] = func() string {
		return f.storage.Site.SiteSubtitle()
	}
	fm["Authenticated"] = func() bool {
		// If user is not authenticated, don't allow them to access the page.
		loggedIn, err := f.sess.UserAuthenticated(r)
		if err != nil {
			// TODO: Need to switch over to the logger.
			log.Println(err)
		}
		return loggedIn
	}
	fm["SiteFooter"] = func() string {
		return f.storage.Site.Footer
	}
	fm["MFAEnabled"] = func() bool {
		return len(os.Getenv("AMB_MFA_KEY")) > 0
	}

	return fm
}
