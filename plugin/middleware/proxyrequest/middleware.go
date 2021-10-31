package proxyrequest

import (
	"net/http"
	"strings"
)

// ProxyRequest will send all requests prefixed with the specified API path to
// the Ambient app while all other requests to the proxy URL.
func (p *Plugin) ProxyRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the path starts with the specific string, serve the API.
		if strings.HasPrefix(r.URL.Path, p.prefixForAPI) {
			next.ServeHTTP(w, r)
			return
		}

		// Else serve the proxy for the UI.
		p.handlerUI.ServeHTTP(w, r)
	})
}
