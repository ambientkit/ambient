// Package plugintemplateengine provides a template engine for plugins.
package plugintemplateengine

import (
	"embed"
	"html/template"
	"net/http"
)

// Renderer represents a plugin template enginer.
type Renderer struct {
	render IRenderer
}

// IRenderer represents an application template enginer.
type IRenderer interface {
	Page(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
	PageContent(w http.ResponseWriter, r *http.Request, content string, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
	Post(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
	PostContent(w http.ResponseWriter, r *http.Request, content string, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
	Error(w http.ResponseWriter, r *http.Request, content string, statusCode int, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
}

// NewRenderer returns a new template engine for plugins.
func NewRenderer(render IRenderer) *Renderer {
	return &Renderer{
		render: render,
	}
}

// Page renders a page.
func (rr *Renderer) Page(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
	return rr.render.Page(w, r, assets, templateName, fm, vars)
}

// PageContent renders page content.
func (rr *Renderer) PageContent(w http.ResponseWriter, r *http.Request, content string, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
	return rr.render.PageContent(w, r, content, fm, vars)
}

// Post renders a post.
func (rr *Renderer) Post(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
	return rr.render.Post(w, r, assets, templateName, fm, vars)
}

// PostContent renders post content.
func (rr *Renderer) PostContent(w http.ResponseWriter, r *http.Request, content string, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
	return rr.render.PostContent(w, r, content, fm, vars)
}

// Error renders an error.
func (rr *Renderer) Error(w http.ResponseWriter, r *http.Request, content string, statusCode int, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
	return rr.render.Error(w, r, content, statusCode, fm, vars)
}
