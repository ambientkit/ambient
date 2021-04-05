package bearblog

import (
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/josephspurrier/ambient/app/core"
)

// FuncMap returns a map of template functions that can be used in templates.
func (p *Plugin) FuncMap(r *http.Request) template.FuncMap {
	fm := make(template.FuncMap)
	fm["Stamp"] = func(t time.Time) string {
		return t.Format("2006-01-02")
	}
	fm["StampFriendly"] = func(t time.Time) string {
		return t.Format("02 Jan, 2006")
	}
	fm["PublishedPages"] = func() []core.Post {
		arr, err := p.Site.PublishedPages()
		if err != nil {
			p.Log.Error("bearblog: error getting published pages: %v", err.Error())
		}
		return arr
	}
	fm["SiteSubtitle"] = func() string {
		subtitle, err := p.Site.PluginField(Subtitle)
		if err != nil {
			p.Log.Error("bearblog: error getting subtitle: %v", err.Error())
		}
		return subtitle
	}
	fm["Authenticated"] = func() bool {
		// If user is not authenticated, don't allow them to access the page.
		loggedIn, err := p.Site.UserAuthenticated(r)
		if err != nil {
			p.Log.Error("bearblog: error getting if user is authenticated: %v", err.Error())
		}
		return loggedIn
	}
	fm["SiteFooter"] = func() string {
		f, err := p.Site.PluginField(Footer)
		if err != nil {
			p.Log.Error("bearblog: error getting footer: %v", err.Error())
		}
		return f
	}
	fm["MFAEnabled"] = func() bool {
		return len(os.Getenv("AMB_MFA_KEY")) > 0
	}

	return fm
}
