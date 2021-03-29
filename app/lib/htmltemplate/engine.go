// Package htmltemplate provides HTML generation using templates.
package htmltemplate

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/josephspurrier/ambient/app/lib/templatebuffer"
)

// TemplateManager represents a function map for templates.
type TemplateManager interface {
	FuncMap(r *http.Request) template.FuncMap
	Templates() embed.FS
}

// AssetInjector represents code that can inject files into a template.
type AssetInjector interface {
	Inject(t *template.Template, r *http.Request, pluginNames []string, pageURL string) (*template.Template, error)
}

// Engine represents a HTML template engine.
type Engine struct {
	allowUnsafeHTML bool
	templateManager TemplateManager
	assetInjector   AssetInjector
	pluginNames     []string
}

// New returns a HTML template engine.
func New(allowUnsafeHTML bool, templateManager TemplateManager, assetInjector AssetInjector, pluginNames []string) *Engine {
	return &Engine{
		allowUnsafeHTML: allowUnsafeHTML,
		templateManager: templateManager,
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

// partialTemplate converts content from markdown to HTML and then outputs to
// a response writer. Returns an HTTP status code and an error if one occurs.
func (te *Engine) partial(w http.ResponseWriter, r *http.Request, mainTemplate string,
	partialTemplate string, statusCode int, vars map[string]interface{}) (status int, err error) {
	// Set the status to passed in value.
	status = statusCode

	// Parse the main template with the functions.
	t, err := te.generateTemplate(r, mainTemplate)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Parse the partial template.
	contentTemplate := fmt.Sprintf("content/%v.tmpl", partialTemplate)
	t, err = t.ParseFS(te.templateManager.Templates(), contentTemplate)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Execute the template and write out if no error.
	err = templatebuffer.ParseExistingTemplate(w, t, status, vars)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return
}

// Post converts a site post from markdown to HTML and then outputs to response
// writer. Returns an HTTP status code and an error if one occurs.
func (te *Engine) Post(w http.ResponseWriter, r *http.Request, mainTemplate string,
	postContent string, vars map[string]interface{}) (status int, err error) {
	// Set the status to OK starting out.
	status = http.StatusOK

	// Parse the main template with the functions.
	t, err := te.generateTemplate(r, mainTemplate)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Parse the content.
	t, err = te.sanitizedContent(t, postContent)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Execute the template and write out if no error.
	err = templatebuffer.ParseExistingTemplate(w, t, status, vars)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return
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

	// Execute the template and write out if no error.
	err = templatebuffer.ParseExistingTemplate(w, t, status, vars)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return
}

func (te *Engine) generateTemplate(r *http.Request, mainTemplate string) (*template.Template, error) {
	// Functions available in the templates.
	fm := te.templateManager.FuncMap(r)

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
	t, err := template.New(path.Base(baseTemplate)).Funcs(fm).ParseFS(te.templateManager.Templates(), templates...)
	if err != nil {
		return nil, err
	}

	// Inject the plugins.
	t, err = te.assetInjector.Inject(t, r, te.pluginNames, r.URL.Path)
	if err != nil {
		return nil, err
	}

	return t, nil
}
