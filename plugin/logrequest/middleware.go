package logrequest

import (
	"net/http"
	"time"
)

// LogRequest will log the HTTP requests.
func (p *Plugin) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p.Log.Info("%v %v %v %v", time.Now().Format("2006-01-02 03:04:05 PM"), r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
