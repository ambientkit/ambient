package securedashboard

import (
	"net/http"
	"strings"
)

// DisallowAnon does not allow anonymous users to access the page.
func (p *Plugin) DisallowAnon(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Don't allow anon users to access the dashboard.
		if strings.HasPrefix(r.URL.Path, "/dashboard") {
			// If user is not authenticated, don't allow them to access the page.
			loggedIn, err := p.Site.UserAuthenticated(r)
			// If there was an error, then return error.
			if err != nil {
				status, _ := p.Site.Error(err)
				p.Mux.Error(status, w, r)
				return
			}
			if !loggedIn {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}
