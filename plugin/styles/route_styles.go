package styles

import (
	"fmt"
	"net/http"
)

// index returns CSS file.
func (p *Plugin) index(w http.ResponseWriter, r *http.Request) (status int, err error) {
	// Get the styles.
	s, err := p.Site.PluginSetting(Styles)
	if err != nil {
		return p.Site.Error(err)
	}

	w.Header().Set("Content-Type", "text/css")

	fmt.Fprint(w, s)
	return
}
