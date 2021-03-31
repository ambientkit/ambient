// Package cachecontrol provides caching.
// Resource: https://www.sanarias.com/blog/115LearningHTTPcachinginGo
package cachecontrol

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
)

// Handle returns true if handled.
func Handle(w http.ResponseWriter, r *http.Request, content []byte) bool {
	// Set the etag for cache control.
	etag := fmt.Sprintf(`"%x"`, md5.Sum(content))
	w.Header().Set("Etag", etag)
	w.Header().Set("Cache-Control", "max-age=2592000") // 30 days
	if match := r.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, etag) {
			w.WriteHeader(http.StatusNotModified)
			return true
		}
	}

	return false
}
