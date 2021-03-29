package robots

import (
	"fmt"
	"net/http"
)

// Robots returns a page for web crawlers.
func (p *Plugin) index(w http.ResponseWriter, r *http.Request) (status int, err error) {
	w.Header().Set("Content-Type", "text/plain")
	text :=
		`User-agent: *
Allow: /`
	fmt.Fprint(w, text)
	return
}
