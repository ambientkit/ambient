package uptimerobotok

import (
	"net/http"
)

// HeadReply will return a 200 for the uptimerobot.
func HeadReply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && r.Method == http.MethodHead {
			return
		}

		next.ServeHTTP(w, r)
	})
}
