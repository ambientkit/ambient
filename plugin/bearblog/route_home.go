package bearblog

import (
	"net/http"
)

func (p *Plugin) index(w http.ResponseWriter, r *http.Request) (status int, err error) {
	content, err := p.Site.Content()
	if err != nil {
		return p.Site.Error(err)
	}

	if content == "" {
		content = "*No content yet.*"
	}

	vars := make(map[string]interface{})
	vars["postcontent"] = sanitized(content)
	return p.Render.PluginPage(w, r, assets, "template/content/home", p.FuncMap(r), vars)
}

func (p *Plugin) edit(w http.ResponseWriter, r *http.Request) (status int, err error) {
	siteContent, err := p.Site.Content()
	if err != nil {
		return p.Site.Error(err)
	}

	siteTitle, err := p.Site.Title()
	if err != nil {
		return p.Site.Error(err)
	}

	siteSubtitle, err := p.Site.PluginField(Subtitle)
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

	footer, err := p.Site.PluginField(Footer)
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

	return p.Render.PluginPage(w, r, assets, "template/content/home_edit", p.FuncMap(r), vars)
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

	err = p.Site.SetContent(r.FormValue("content"))
	if err != nil {
		return p.Site.Error(err)
	}

	err = p.Site.SetScheme(r.FormValue("scheme"))
	if err != nil {
		return p.Site.Error(err)
	}

	err = p.Site.SetURL(r.FormValue("domain"))
	if err != nil {
		return p.Site.Error(err)
	}

	err = p.Site.SetPluginField(Subtitle, r.FormValue("subtitle"))
	if err != nil {
		return p.Site.Error(err)
	}

	err = p.Site.SetPluginField(LoginURL, r.FormValue("loginurl"))
	if err != nil {
		return p.Site.Error(err)
	}

	err = p.Site.SetPluginField(Footer, r.FormValue("footer"))
	if err != nil {
		return p.Site.Error(err)
	}

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
