package pluginsafe

import (
	"html/template"
	"net/http"

	"github.com/ambientkit/ambient"
)

// TemplateRenderer represents a plugin template enginer.
type TemplateRenderer struct {
	render ambient.Renderer
}

// NewRenderer returns a new template engine for plugins.
func NewRenderer(render ambient.Renderer) *TemplateRenderer {
	return &TemplateRenderer{
		render: render,
	}
}

// globalFuncMapCallable returns a callable function.
func globalFuncMapCallable(r *http.Request, fm func(r *http.Request) template.FuncMap) func(r *http.Request) template.FuncMap {
	var f = template.FuncMap{}
	if fm != nil {
		f = fm(r)
	}
	return func(r *http.Request) template.FuncMap {
		return ambient.GlobalFuncMap(f)
	}
}

// Page renders a page.
func (rr *TemplateRenderer) Page(w http.ResponseWriter, r *http.Request,
	assets ambient.FileSystemReader, templateName string, fm func(r *http.Request) template.FuncMap,
	vars map[string]interface{}) (err error) {
	return rr.render.Page(w, r, assets, templateName, globalFuncMapCallable(r, fm), vars)
}

// PageContent renders page content.
func (rr *TemplateRenderer) PageContent(w http.ResponseWriter, r *http.Request,
	content string, fm func(r *http.Request) template.FuncMap,
	vars map[string]interface{}) (err error) {
	return rr.render.PageContent(w, r, content, globalFuncMapCallable(r, fm), vars)
}

// Post renders a post.
func (rr *TemplateRenderer) Post(w http.ResponseWriter, r *http.Request,
	assets ambient.FileSystemReader, templateName string, fm func(r *http.Request) template.FuncMap,
	vars map[string]interface{}) (err error) {
	return rr.render.Post(w, r, assets, templateName, globalFuncMapCallable(r, fm), vars)
}

// PostContent renders post content.
func (rr *TemplateRenderer) PostContent(w http.ResponseWriter, r *http.Request,
	content string, fm func(r *http.Request) template.FuncMap,
	vars map[string]interface{}) (err error) {
	return rr.render.PostContent(w, r, content, globalFuncMapCallable(r, fm), vars)
}

// Error renders an error.
func (rr *TemplateRenderer) Error(w http.ResponseWriter, r *http.Request,
	content string, statusCode int, fm func(r *http.Request) template.FuncMap,
	vars map[string]interface{}) (err error) {
	return rr.render.Error(w, r, content, statusCode, globalFuncMapCallable(r, fm), vars)
}
