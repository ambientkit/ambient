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

// Templates is the embedded files.
//go:embed *
var Templates embed.FS

// FuncMap returns a map of template functions that can be used in templates.
func FuncMap(r *http.Request, storage *datastorage.Storage, sess *websession.Session) template.FuncMap {
	fm := make(template.FuncMap)
	fm["Stamp"] = func(t time.Time) string {
		return t.Format("2006-01-02")
	}
	fm["StampFriendly"] = func(t time.Time) string {
		return t.Format("02 Jan, 2006")
	}
	fm["PublishedPages"] = func() []model.Post {
		return storage.Site.PublishedPages()
	}
	fm["SiteURL"] = func() string {
		return storage.Site.SiteURL()
	}
	fm["SiteTitle"] = func() string {
		return storage.Site.SiteTitle()
	}
	fm["SiteSubtitle"] = func() string {
		return storage.Site.SiteSubtitle()
	}
	fm["SiteDescription"] = func() string {
		return storage.Site.Description
	}
	fm["SiteAuthor"] = func() string {
		return storage.Site.Author
	}
	fm["SiteFavicon"] = func() string {
		return storage.Site.Favicon
	}
	fm["Authenticated"] = func() bool {
		// If user is not authenticated, don't allow them to access the page.
		_, loggedIn := sess.User(r)
		return loggedIn
	}
	fm["SiteFooter"] = func() string {
		return storage.Site.Footer
	}
	fm["MFAEnabled"] = func() bool {
		return len(os.Getenv("AMB_MFA_KEY")) > 0
	}
	fm["SiteStyles"] = func() template.CSS {
		return template.CSS(storage.Site.Styles)
	}

	return fm
}
