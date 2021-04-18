package ambient

import (
	"embed"
	"html/template"
	"net/http"
)

// TemplateRenderer represents a plugin template enginer.
type TemplateRenderer struct {
	render Renderer
}

// NewRenderer returns a new template engine for plugins.
func NewRenderer(render Renderer) *TemplateRenderer {
	return &TemplateRenderer{
		render: render,
	}
}

// Page renders a page.
func (rr *TemplateRenderer) Page(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
	return rr.render.Page(w, r, assets, templateName, fm, vars)
}

// PageContent renders page content.
func (rr *TemplateRenderer) PageContent(w http.ResponseWriter, r *http.Request, content string, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
	return rr.render.PageContent(w, r, content, fm, vars)
}

// Post renders a post.
func (rr *TemplateRenderer) Post(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
	return rr.render.Post(w, r, assets, templateName, fm, vars)
}

// PostContent renders post content.
func (rr *TemplateRenderer) PostContent(w http.ResponseWriter, r *http.Request, content string, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
	return rr.render.PostContent(w, r, content, fm, vars)
}

// Error renders an error.
func (rr *TemplateRenderer) Error(w http.ResponseWriter, r *http.Request, content string, statusCode int, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
	return rr.render.Error(w, r, content, statusCode, fm, vars)
}
