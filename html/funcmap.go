// Package html is the base templates and functions for rendering the application.
package html

import (
	"embed"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/josephspurrier/ambient/app/lib/datastorage"
	"github.com/josephspurrier/ambient/app/lib/websession"
	"github.com/josephspurrier/ambient/app/model"
)

//go:embed *
var templates embed.FS

// TemplateManager represents an object that returns templates and a FuncMap.
type TemplateManager struct {
	storage *datastorage.Storage
	sess    *websession.Session
}

// NewTemplateManager returns a TemplateManager.
func NewTemplateManager(storage *datastorage.Storage, sess *websession.Session) *TemplateManager {
	return &TemplateManager{
		storage: storage,
		sess:    sess,
	}
}

// Templates returns the embedded templates.
func (f *TemplateManager) Templates() embed.FS {
	return templates
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
	fm["SiteDescription"] = func() string {
		return f.storage.Site.Description
	}
	fm["SiteAuthor"] = func() string {
		return f.storage.Site.Author
	}
	fm["SiteFavicon"] = func() string {
		return f.storage.Site.Favicon
	}
	fm["Authenticated"] = func() bool {
		// If user is not authenticated, don't allow them to access the page.
		_, loggedIn := f.sess.User(r)
		return loggedIn
	}
	fm["SiteFooter"] = func() string {
		return f.storage.Site.Footer
	}
	fm["MFAEnabled"] = func() bool {
		return len(os.Getenv("AMB_MFA_KEY")) > 0
	}
	fm["SiteStyles"] = func() template.CSS {
		return template.CSS(f.storage.Site.Styles)
	}

	return fm
}
