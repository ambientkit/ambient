package pluginsafe

import (
	"embed"
	"html/template"
	"net/http"
	"os"

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

// globalFuncMap adds the URL prefix to the FuncMap.
func globalFuncMap(fm template.FuncMap) template.FuncMap {
	if fm == nil {
		fm = template.FuncMap{}
	}

	fm["URLPrefix"] = func() string {
		return os.Getenv("AMB_URL_PREFIX")
	}

	return fm
}

// globalFuncMapCallable returns a callable function.
func globalFuncMapCallable(r *http.Request, fm func(r *http.Request) template.FuncMap) func(r *http.Request) template.FuncMap {
	var f = template.FuncMap{}
	if fm != nil {
		f = fm(r)
	}
	return func(r *http.Request) template.FuncMap {
		return globalFuncMap(f)
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
	assets embed.FS, templateName string, fm func(r *http.Request) template.FuncMap,
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
