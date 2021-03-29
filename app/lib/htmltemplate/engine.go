// Package htmltemplate provides HTML generation using templates.
package htmltemplate

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/josephspurrier/ambient/app/lib/datastorage"
	"github.com/josephspurrier/ambient/app/lib/websession"
	"github.com/josephspurrier/ambient/app/model"
	"github.com/josephspurrier/ambient/html"
	"github.com/oxtoacart/bpool"
)

// AssetInjector represents code that can inject files into a template.
type AssetInjector interface {
	InjectPlugins(t *template.Template, r *http.Request, pluginNames []string, pageURL string) (*template.Template, error)
}

// Engine represents a HTML template engine.
type Engine struct {
	allowUnsafeHTML bool
	storage         *datastorage.Storage
	sess            *websession.Session
	assetInjector   AssetInjector
	pluginNames     []string
}

// New returns a HTML template engine.
func New(allowUnsafeHTML bool, storage *datastorage.Storage, sess *websession.Session,
	pluginNames []string, assetInjector AssetInjector) *Engine {
	return &Engine{
		allowUnsafeHTML: allowUnsafeHTML,
		storage:         storage,
		sess:            sess,
		assetInjector:   assetInjector,
		pluginNames:     pluginNames,
	}
}

// Template renders HTML to a response writer and returns a 200 status code and
// an error if one occurs.
func (te *Engine) Template(w http.ResponseWriter, r *http.Request, mainTemplate string,
	partialTemplate string, vars map[string]interface{}) (status int, err error) {
	return te.partial(w, r, mainTemplate, partialTemplate, http.StatusOK, vars)
}

// ErrorTemplate renders HTML to a response writer and returns a 404 status code
// and an error if one occurs.
func (te *Engine) ErrorTemplate(w http.ResponseWriter, r *http.Request, mainTemplate string,
	partialTemplate string, vars map[string]interface{}) (status int, err error) {
	return te.partial(w, r, mainTemplate, partialTemplate, http.StatusNotFound, vars)
}

// bufpool is used to write out HTML after it's been executed and before it's
// written to the ResponseWriter to catch any partially written templates.
var bufpool *bpool.BufferPool = bpool.NewBufferPool(64)

// partialTemplate converts content from markdown to HTML and then outputs to
// a response writer. Returns an HTTP status code and an error if one occurs.
func (te *Engine) partial(w http.ResponseWriter, r *http.Request, mainTemplate string,
	partialTemplate string, statusCode int, vars map[string]interface{}) (status int, err error) {
	// Parse the main template with the functions.
	t, err := te.generateTemplate(r, mainTemplate)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Parse the partial template.
	contentTemplate := fmt.Sprintf("content/%v.tmpl", partialTemplate)
	t, err = t.ParseFS(html.Templates, contentTemplate)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Output the status code.
	w.WriteHeader(statusCode)

	// Write temporarily to a buffer pool
	buf := bufpool.Get()
	defer bufpool.Put(buf)

	// Execute the template.
	if err := t.Execute(buf, vars); err != nil {
		return http.StatusInternalServerError, err
	}

	// Write out the template.
	buf.WriteTo(w)

	return statusCode, nil
}

// Post converts a site post from markdown to HTML and then outputs to response
// writer. Returns an HTTP status code and an error if one occurs.
func (te *Engine) Post(w http.ResponseWriter, r *http.Request, mainTemplate string,
	post model.Post, vars map[string]interface{}) (status int, err error) {
	// Display 404 if not found.
	if post.URL == "" {
		return http.StatusNotFound, nil
	}

	// Parse the main template with the functions.
	t, err := te.generateTemplate(r, mainTemplate)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Parse the content.
	t, err = te.sanitizedContent(t, post.Content)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Write temporarily to a buffer pool
	buf := bufpool.Get()
	defer bufpool.Put(buf)

	// Execute the template.
	if err := t.Execute(buf, vars); err != nil {
		return http.StatusInternalServerError, err
	}

	// Write out the template.
	buf.WriteTo(w)

	return http.StatusOK, nil
}

// PluginTemplate -
func (te *Engine) PluginTemplate(w http.ResponseWriter, r *http.Request, assets embed.FS,
	partialTemplate string, vars map[string]interface{}) (status int, err error) {
	// Set the status to OK starting out.
	status = http.StatusOK

	// Parse the main template with the functions.
	// FIXME: Shouldn't just use the dashboard layout.
	t, err := te.generateTemplate(r, "layout/dashboard")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Parse the plugin template.
	t, err = t.ParseFS(assets, partialTemplate)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Write temporarily to a buffer pool
	buf := bufpool.Get()
	defer bufpool.Put(buf)

	// Execute the template.
	if err := t.Execute(buf, vars); err != nil {
		return http.StatusInternalServerError, err
	}

	// Output the status code.
	w.WriteHeader(status)

	// Write out the template.
	buf.WriteTo(w)

	return
}

func (te *Engine) generateTemplate(r *http.Request, mainTemplate string) (*template.Template, error) {
	// Functions available in the templates.
	fm := html.FuncMap(r, te.storage, te.sess)

	// Generate list of templates.
	baseTemplate := fmt.Sprintf("%v.tmpl", mainTemplate)
	templates := []string{
		baseTemplate,
		"partial/head.tmpl",
		"partial/header.tmpl",
		"partial/nav.tmpl",
		"partial/footer.tmpl",
	}

	// Parse the main template with the functions.
	t, err := template.New(path.Base(baseTemplate)).Funcs(fm).ParseFS(html.Templates, templates...)
	if err != nil {
		return nil, err
	}

	// Inject the plugins.
	t, err = te.assetInjector.InjectPlugins(t, r, te.pluginNames, r.URL.Path)
	if err != nil {
		return nil, err
	}

	return t, nil
}
