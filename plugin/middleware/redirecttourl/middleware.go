package redirecttourl

import (
	"fmt"
	"net/http"
	"strings"
)

// StripSlash will strip trailing slashes from requests.
func (p *Plugin) stripSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the site scheme.
		siteScheme, err := p.Site.PluginSettingString(SiteScheme)
		if err != nil || len(siteScheme) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// Get the site URL.
		siteURL, err := p.Site.PluginSettingString(SiteURL)
		if err != nil || len(siteURL) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// Redirect to the correct website if the values are set.
		if len(siteURL) > 0 && len(siteScheme) > 0 && !strings.Contains(r.Host, siteURL) {
			http.Redirect(w, r, fmt.Sprintf("%v://%v%v", siteScheme, siteURL, r.URL.Path), http.StatusPermanentRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}
