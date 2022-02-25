package ambient

import (
	"embed"
	"html/template"
	"net/http"
	"os"
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
func globalFuncMapCallable(fm template.FuncMap) func(r *http.Request) template.FuncMap {
	return func(r *http.Request) template.FuncMap {
		return globalFuncMap(fm)
	}
}

// Page renders a page.
func (rr *TemplateRenderer) Page(w http.ResponseWriter, r *http.Request,
	assets embed.FS, templateName string, fm func(r *http.Request) template.FuncMap,
	vars map[string]interface{}) (status int, err error) {
	return rr.render.Page(w, r, assets, templateName, globalFuncMapCallable(fm(r)), vars)
}

// PageContent renders page content.
func (rr *TemplateRenderer) PageContent(w http.ResponseWriter, r *http.Request,
	content string, fm func(r *http.Request) template.FuncMap,
	vars map[string]interface{}) (status int, err error) {
	return rr.render.PageContent(w, r, content, globalFuncMapCallable(fm(r)), vars)
}

// Post renders a post.
func (rr *TemplateRenderer) Post(w http.ResponseWriter, r *http.Request,
	assets embed.FS, templateName string, fm func(r *http.Request) template.FuncMap,
	vars map[string]interface{}) (status int, err error) {
	return rr.render.Post(w, r, assets, templateName, globalFuncMapCallable(fm(r)), vars)
}

// PostContent renders post content.
func (rr *TemplateRenderer) PostContent(w http.ResponseWriter, r *http.Request,
	content string, fm func(r *http.Request) template.FuncMap,
	vars map[string]interface{}) (status int, err error) {
	return rr.render.PostContent(w, r, content, globalFuncMapCallable(fm(r)), vars)
}

// Error renders an error.
func (rr *TemplateRenderer) Error(w http.ResponseWriter, r *http.Request,
	content string, statusCode int, fm func(r *http.Request) template.FuncMap,
	vars map[string]interface{}) (status int, err error) {
	return rr.render.Error(w, r, content, statusCode, globalFuncMapCallable(fm(r)), vars)
}
