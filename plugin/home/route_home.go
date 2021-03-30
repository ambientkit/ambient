package home

import (
	"net/http"
)

func (p *Plugin) index(w http.ResponseWriter, r *http.Request) (status int, err error) {
	content, err := p.Site.Content()
	if err != nil {
		return p.Site.Error(err)
	}

	if content == "" {
		content = "*No content yet.*"
	}

	vars := make(map[string]interface{})
	return p.Render.PluginPageContent(w, r, content, vars)
}
