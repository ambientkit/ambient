package styles

import (
	"fmt"
	"net/http"
)

// index returns CSS file.
func (p *Plugin) index(w http.ResponseWriter, r *http.Request) (status int, err error) {
	// Get the styles.
	s, err := p.Site.PluginField(Styles)
	if err != nil {
		return p.Site.Error(err)
	}

	w.Header().Set("Content-Type", "text/css")

	fmt.Fprint(w, s)
	return
}

func (p *Plugin) edit(w http.ResponseWriter, r *http.Request) (status int, err error) {
	// Get the favicon.
	favicon, err := p.Site.PluginField(Favicon)
	if err != nil {
		return p.Site.Error(err)
	}

	// Get the styles.
	s, err := p.Site.PluginField(Styles)
	if err != nil {
		return p.Site.Error(err)
	}

	vars := make(map[string]interface{})
	vars["title"] = "Site styles"
	vars["token"] = p.Security.SetCSRF(r)
	vars["favicon"] = favicon
	vars["styles"] = s

	return p.Render.PluginDashboard(w, r, assets, "styles_edit", vars)
}

func (p *Plugin) update(w http.ResponseWriter, r *http.Request) (status int, err error) {
	r.ParseForm()

	// CSRF protection.
	success := p.Security.CSRF(r)
	if !success {
		return http.StatusBadRequest, nil
	}

	err = p.Site.SetPluginField(Favicon, r.FormValue("favicon"))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = p.Site.SetPluginField(Styles, r.FormValue("styles"))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/dashboard/styles", http.StatusFound)
	return
}
