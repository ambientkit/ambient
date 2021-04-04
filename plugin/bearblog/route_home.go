package bearblog

import (
	"net/http"
)

func (p *Plugin) edit(w http.ResponseWriter, r *http.Request) (status int, err error) {
	siteContent, err := p.Site.Content()
	if err != nil {
		return p.Site.Error(err)
	}

	siteTitle, err := p.Site.Title()
	if err != nil {
		return p.Site.Error(err)
	}

	siteSubtitle, err := p.Site.Subtitle()
	if err != nil {
		return p.Site.Error(err)
	}

	siteURL, err := p.Site.URL()
	if err != nil {
		return p.Site.Error(err)
	}

	siteScheme, err := p.Site.Scheme()
	if err != nil {
		return p.Site.Error(err)
	}

	loginURL, err := p.Site.PluginField(LoginURL)
	if err != nil {
		return p.Site.Error(err)
	}

	footer, err := p.Site.Footer()
	if err != nil {
		return p.Site.Error(err)
	}

	vars := make(map[string]interface{})
	vars["title"] = "Edit site"
	vars["homeContent"] = siteContent
	vars["ptitle"] = siteTitle
	vars["subtitle"] = siteSubtitle
	vars["token"] = p.Security.SetCSRF(r)

	// Help the user set the domain based off the current URL.
	if siteURL == "" {
		vars["domain"] = r.Host
	} else {
		vars["domain"] = siteURL
	}

	vars["scheme"] = siteScheme
	vars["loginurl"] = loginURL
	vars["footer"] = footer

	return p.Render.PluginDashboard(w, r, assets, "template/home_edit", p.FuncMap(r), vars)
}

func (p *Plugin) update(w http.ResponseWriter, r *http.Request) (status int, err error) {
	r.ParseForm()

	// CSRF protection.
	success := p.Security.CSRF(r)
	if !success {
		return http.StatusBadRequest, nil
	}

	err = p.Site.SetTitle(r.FormValue("title"))
	if err != nil {
		return p.Site.Error(err)
	}

	err = p.Site.SetSubtitle(r.FormValue("subtitle"))
	if err != nil {
		return p.Site.Error(err)
	}

	//c.Storage.Site.Title =
	//c.Storage.Site.Subtitle = r.FormValue("subtitle")
	// c.Storage.Site.URL = r.FormValue("domain")
	// c.Storage.Site.Content = r.FormValue("content")
	// c.Storage.Site.Scheme = r.FormValue("scheme")
	// c.Storage.Site.LoginURL = r.FormValue("loginurl")
	// c.Storage.Site.Footer = r.FormValue("footer")
	// c.Storage.Site.Updated = time.Now()

	// err = c.Storage.Save()
	// if err != nil {
	// 	return http.StatusInternalServerError, err
	// }

	http.Redirect(w, r, "/dashboard", http.StatusFound)
	return
}

func (p *Plugin) reload(w http.ResponseWriter, r *http.Request) (status int, err error) {
	err = p.Site.Load()
	if err != nil {
		p.Site.Error(err)
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)
	return
}
