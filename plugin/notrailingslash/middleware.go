package notrailingslash

import (
	"net/http"
	"strings"
)

// StripSlash will strip trailing slashes from requests.
func (p *Plugin) stripSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Don't allow access to files with a slash at the end.
		if strings.Contains(r.URL.Path, ".") && strings.HasSuffix(r.URL.Path, "/") {
			p.Mux.Error(http.StatusNotFound, w, r)
			return
		}

		// Strip trailing slash.
		if r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, strings.TrimRight(r.URL.Path, "/"), http.StatusPermanentRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}
