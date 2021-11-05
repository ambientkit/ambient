package etagcache

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
)

// Handler returns middleware.
func (p *Plugin) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the MaxAge setting.
		maxAge, err := p.Site.PluginSettingString(MaxAge)
		if err != nil || len(maxAge) == 0 {
			// Bypass.
			next.ServeHTTP(w, r)
			return
		}

		// Run the request by pass in a custom response writer.
		recorder := NewCustomResponseWriter()
		next.ServeHTTP(recorder, r)

		// Set the etag for cache control.
		newEtag := fmt.Sprintf(`"%x"`, md5.Sum(recorder.body))
		w.Header().Set("Etag", newEtag)
		w.Header().Set("Cache-Control", "max-age="+maxAge)

		if newEtag == r.Header.Get("Etag") {
			// If the etag is the same, send a not modified status code.
			w.WriteHeader(http.StatusNotModified)
			return
		} else if match := r.Header.Get("If-None-Match"); match != "" {
			// Else if the header is set and one of the values are set,
			// then send a not modified status code.
			if strings.Contains(match, newEtag) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		// Write to the response writer.
		recorder.WriteTo(w)
	})
}
