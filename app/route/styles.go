package route

import (
	"net/http"

	"github.com/josephspurrier/ambient/app/core"
)

// Styles -
type Styles struct {
	*core.App
}

func registerStyles(c *Styles) {
	c.Router.Get("/dashboard/styles", c.edit)
	c.Router.Post("/dashboard/styles", c.update)
}

func (c *Styles) edit(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "Site styles"
	vars["token"] = c.Sess.SetCSRF(r)
	vars["favicon"] = c.Storage.Site.Favicon
	vars["styles"] = c.Storage.Site.Styles

	return c.Render.Template(w, r, "layout/dashboard", "styles_edit", vars)
}

func (c *Styles) update(w http.ResponseWriter, r *http.Request) (status int, err error) {
	r.ParseForm()

	// CSRF protection.
	success := c.Sess.CSRF(r)
	if !success {
		return http.StatusBadRequest, nil
	}

	c.Storage.Site.Favicon = r.FormValue("favicon")
	c.Storage.Site.Styles = r.FormValue("styles")

	err = c.Storage.Save()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/dashboard/styles", http.StatusFound)
	return
}
