// Package htmltemplate provides HTML generation using templates.
package htmltemplate

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/josephspurrier/ambient/app/lib/templatebuffer"
)

// TemplateManager represents a function map for templates.
type TemplateManager interface {
	FuncMap(r *http.Request) template.FuncMap
	Templates() *embed.FS
}

// AssetInjector represents code that can inject files into a template.
type AssetInjector interface {
	Inject(t *template.Template, r *http.Request, pluginNames []string, layoutType string) (*template.Template, error)
}

// Engine represents a HTML template engine.
type Engine struct {
	allowUnsafeHTML bool
	templateManager TemplateManager
	assetInjector   AssetInjector
	pluginNames     []string
}

// NewTemplateEngine returns a HTML template engine.
func NewTemplateEngine(templateManager TemplateManager, assetInjector AssetInjector, pluginNames []string) *Engine {
	allowUnsafeHTML, err := strconv.ParseBool(os.Getenv("AMB_ALLOW_HTML"))
	if err != nil {
		log.Printf("environment variable not able to parse as bool: %v", "AMB_ALLOW_HTML")
		return nil
	}

	return &Engine{
		allowUnsafeHTML: allowUnsafeHTML,
		templateManager: templateManager,
		assetInjector:   assetInjector,
		pluginNames:     pluginNames,
	}
}

// Post renders using the post layout.
func (te *Engine) Post(w http.ResponseWriter, r *http.Request, postContent string, vars map[string]interface{}) (status int, err error) {
	return te.post(w, r, "layout/post", "post", postContent, vars)
}

// Page renders using the page layout.
func (te *Engine) Page(w http.ResponseWriter, r *http.Request, partialTemplate string, vars map[string]interface{}) (status int, err error) {
	return te.partial(w, r, "layout/bloglist", "page", partialTemplate, http.StatusOK, vars)
}

// Bloglist renders using the bloglist layout.
func (te *Engine) Bloglist(w http.ResponseWriter, r *http.Request, partialTemplate string, vars map[string]interface{}) (status int, err error) {
	return te.partial(w, r, "layout/bloglist", "bloglist", partialTemplate, http.StatusOK, vars)
}

// Dashboard renders using the dashboard layout.
func (te *Engine) Dashboard(w http.ResponseWriter, r *http.Request, partialTemplate string, vars map[string]interface{}) (status int, err error) {
	return te.partial(w, r, "layout/dashboard", "dashboard", partialTemplate, http.StatusOK, vars)
}

// Error renders HTML to a response writer and returns a 404 status code
// and an error if one occurs.
func (te *Engine) Error(w http.ResponseWriter, r *http.Request, partialTemplate string, vars map[string]interface{}) (status int, err error) {
	return te.partial(w, r, "layout/page", "page", partialTemplate, http.StatusNotFound, vars)
}

// partialTemplate converts content from markdown to HTML and then outputs to
// a response writer. Returns an HTTP status code and an error if one occurs.
func (te *Engine) partial(w http.ResponseWriter, r *http.Request, mainTemplate string, layoutType string, partialTemplate string, statusCode int, vars map[string]interface{}) (status int, err error) {
	// Parse the main template with the functions.
	t, err := te.generateTemplate(r, mainTemplate, layoutType)
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
	err = templatebuffer.ParseExistingTemplate(w, r, t, statusCode, vars)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return
}

// Post converts a site post from markdown to HTML and then outputs to response
// writer. Returns an HTTP status code and an error if one occurs.
func (te *Engine) post(w http.ResponseWriter, r *http.Request, mainTemplate string, layoutType string, postContent string, vars map[string]interface{}) (status int, err error) {
	// Parse the main template with the functions.
	t, err := te.generateTemplate(r, mainTemplate, layoutType)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Parse the content.
	t, err = te.sanitizedContent(t, postContent)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Execute the template and write out if no error.
	err = templatebuffer.ParseExistingTemplate(w, r, t, http.StatusOK, vars)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return
}

func (te *Engine) generateTemplate(r *http.Request, mainTemplate string, layoutType string) (*template.Template, error) {
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
	t, err = te.assetInjector.Inject(t, r, te.pluginNames, layoutType)
	if err != nil {
		return nil, err
	}

	return t, nil
}
