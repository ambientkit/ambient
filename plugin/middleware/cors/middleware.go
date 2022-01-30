package cors

import (
	"net/http"
	"strings"
)

// CORS will allow any source to interact with the API.
func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Apply CORS to /api/ only.
		if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS, POST, PUT")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Accept, Accept-Encoding, X-Requested-With, Content-Type")
			if r.Method == "OPTIONS" {
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}
