package html

import (
	"crypto/md5"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/josephspurrier/ambient/app/lib/datastorage"
	"github.com/josephspurrier/ambient/app/lib/envdetect"
	"github.com/josephspurrier/ambient/app/lib/websession"
	"github.com/josephspurrier/ambient/app/model"
	"github.com/josephspurrier/ambient/assets"
)

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
	fm["GoogleAnalyticsID"] = func() string {
		if envdetect.RunningLocalDev() {
			return ""
		}
		return storage.Site.GoogleAnalyticsID
	}
	fm["DisqusID"] = func() string {
		if envdetect.RunningLocalDev() {
			return ""
		}
		return storage.Site.DisqusID
	}
	fm["SiteFooter"] = func() string {
		return storage.Site.Footer
	}
	fm["MFAEnabled"] = func() bool {
		return len(os.Getenv("AMB_MFA_KEY")) > 0
	}
	fm["AssetStamp"] = func(f string) string {
		return assetTimePath(f)
	}
	fm["SiteStyles"] = func() template.CSS {
		return template.CSS(storage.Site.Styles)
	}
	fm["StylesAppend"] = func() bool {
		if len(storage.Site.Styles) == 0 {
			// If there are no style, then always append.
			return true
		} else if storage.Site.StylesAppend {
			// Else if there are style and it's append, then append.
			return true
		}
		return false
	}

	return fm
}

// assetTimePath returns a URL with a MD5 hash appended.
func assetTimePath(s string) string {
	// Use the root directory.
	fsys, err := fs.Sub(assets.CSS, ".")
	if err != nil {
		return s
	}

	// Get the requested file name.
	fname := strings.TrimPrefix(s, "/assets/")

	// Open the file.
	f, err := fsys.Open(fname)
	if err != nil {
		return s
	}
	defer f.Close()

	// Get all the content.s
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return s
	}

	// Create a hash.
	hsh := md5.New()
	hsh.Write(b)

	return fmt.Sprintf("%v?%x", s, hsh.Sum(nil))
}
