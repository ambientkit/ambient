package simplelogin

import (
	"html/template"
	"net/http"
	"path"
	"time"

	"github.com/ambientkit/ambient"
)

// funcMap returns a map of template functions that can be used in templates.
func (p *Plugin) funcMap(r *http.Request) template.FuncMap {
	fm := make(template.FuncMap)
	fm["simplelogin_Stamp"] = func(t time.Time) string {
		return t.Format("2006-01-02")
	}
	fm["simplelogin_StampFriendly"] = func(t time.Time) string {
		return t.Format("02 Jan, 2006")
	}
	fm["simplelogin_PublishedPages"] = func() []ambient.Post {
		arr, err := p.Site.PublishedPages()
		if err != nil {
			p.Log.Warn("simplelogin: error getting published pages: %v", err.Error())
		}
		return arr
	}
	fm["simplelogin_SiteSubtitle"] = func() string {
		subtitle, err := p.Site.PluginSettingString(Subtitle)
		if err != nil {
			p.Log.Warn("simplelogin: error getting subtitle: %v", err.Error())
		}
		return subtitle
	}
	fm["simplelogin_Authenticated"] = func() bool {
		// If user is not authenticated, don't allow them to access the page.
		_, err := p.Site.AuthenticatedUser(r)
		return err == nil
	}
	fm["simplelogin_SiteFooter"] = func() string {
		f, err := p.Site.PluginSettingString(Footer)
		if err != nil {
			p.Log.Warn("simplelogin: error getting footer: %v", err.Error())
		}
		return f
	}
	fm["simplelogin_PageURL"] = func() string {
		siteURL, err := p.Site.FullURL()
		if err != nil {
			p.Log.Warn("simplelogin: error getting site URL: %v", err.Error())
		}

		return path.Join(siteURL, r.URL.Path)
	}
	fm["simplelogin_MFAEnabled"] = func() bool {
		mfakey, err := p.Site.PluginSettingString(MFAKey)
		if err != nil {
			p.Log.Warn("simplelogin: error getting MFA key: %v", err.Error())
		}
		return len(mfakey) > 0
	}

	return fm
}
