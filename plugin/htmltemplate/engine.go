package htmltemplate

import (
	"crypto/rand"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/josephspurrier/ambient/app/core"
	"github.com/josephspurrier/ambient/app/lib/templatebuffer"
)

// Engine represents a HTML template engine.
type Engine struct {
	allowUnsafeHTML bool
	assetInjector   core.AssetInjector
	pluginNames     []string
	escape          bool
	log             core.ILogger
}

// NewTemplateEngine returns a HTML template engine.
func NewTemplateEngine(logger core.ILogger, assetInjector core.AssetInjector, pluginNames []string) *Engine {
	allowHTML := os.Getenv("AMB_ALLOW_HTML")
	allowUnsafeHTML, err := strconv.ParseBool(allowHTML)
	if err != nil {
		logger.Error("htmltemplate: environment variable is not able to parse as bool (%v=%v): %v", "AMB_ALLOW_HTML", allowHTML, err.Error())
		return nil
	}

	//TODO: Add a setting to enable or disable escaping.

	return &Engine{
		allowUnsafeHTML: allowUnsafeHTML,
		assetInjector:   assetInjector,
		pluginNames:     pluginNames,
		escape:          true,
		log:             logger,
	}
}

// Page renders using the page layout.
func (te *Engine) Page(w http.ResponseWriter, r *http.Request, assets embed.FS, partialTemplate string, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
	return te.pluginPartial(w, r, "layout/page", "page", assets, partialTemplate, http.StatusOK, fm, vars)
}

// PageContent renders using the page content.
func (te *Engine) PageContent(w http.ResponseWriter, r *http.Request, content string, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
	return te.pluginContent(w, r, "layout/page", "page", content, http.StatusOK, fm, vars)
}

// Post renders using the post layout.
func (te *Engine) Post(w http.ResponseWriter, r *http.Request, assets embed.FS, partialTemplate string, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
	return te.pluginPartial(w, r, "layout/page", "post", assets, partialTemplate, http.StatusOK, fm, vars)
}

// PostContent renders using the post content.
func (te *Engine) PostContent(w http.ResponseWriter, r *http.Request, content string, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
	return te.pluginContent(w, r, "layout/page", "post", content, http.StatusOK, fm, vars)
}

// Error renders HTML to a response writer and returns a 404 status code
// and an error if one occurs.
func (te *Engine) Error(w http.ResponseWriter, r *http.Request, content string, statusCode int, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
	return te.pluginContent(w, r, "layout/page", "page", content, statusCode, fm, vars)
}

func (te *Engine) pluginPartial(w http.ResponseWriter, r *http.Request, mainTemplate string, layoutType string, assets embed.FS, partialTemplate string, statusCode int, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
	// Parse the main template with the functions.
	t, err := te.generateTemplate(r, mainTemplate, layoutType, vars)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Parse the plugin template separately for security.
	content, err := templatebuffer.ParseTemplateFS(assets, fmt.Sprintf("%v.tmpl", partialTemplate), fm, vars)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Escape the content.
	t, err = te.escapeContent(t, "content", content)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Execute the template and write out if no error.
	err = templatebuffer.ParseExistingTemplate(w, r, t, statusCode, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return
}

// pluginContent converts a site post from markdown to HTML and then outputs to response
// writer. Returns an HTTP status code and an error if one occurs.
func (te *Engine) pluginContent(w http.ResponseWriter, r *http.Request, mainTemplate string, layoutType string, postContent string, statusCode int, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
	// Parse the main template with the functions.
	t, err := te.generateTemplate(r, mainTemplate, layoutType, vars)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	//TODO: If we were going to use a filter on content, this is where it would go.

	// Parse the plugin template separately for security.
	content, err := templatebuffer.ParseTemplate(postContent, fm, vars)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Escape the content.
	t, err = te.escapeContent(t, "content", content)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Execute the template and write out if no error.
	err = templatebuffer.ParseExistingTemplate(w, r, t, statusCode, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return
}

func (te *Engine) generateTemplate(r *http.Request, mainTemplate string, layoutType string, vars map[string]interface{}) (*template.Template, error) {
	// Generate list of templates.
	baseTemplate := fmt.Sprintf("%v.tmpl", mainTemplate)
	templates := []string{
		baseTemplate,
	}

	// Parse the main template with the functions.
	t, err := template.New(path.Base(baseTemplate)).ParseFS(assets, templates...)
	if err != nil {
		return nil, err
	}

	// Inject the plugins.
	t, err = te.assetInjector.Inject(t, r, te.pluginNames, layoutType, vars)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// escapeContent returns an escaped content block or an error is one occurs.
func (te *Engine) escapeContent(t *template.Template, name string, content string) (*template.Template, error) {
	if !te.escape {
		safeContent := fmt.Sprintf(`{{define "%s"}}%s{{end}}`, name, content)
		var err error
		t, err = t.Parse(safeContent)
		if err != nil {
			return nil, err
		}

		return t, nil
	}

	// Generate a random UUID.
	uuid, err := generateUUID()
	if err != nil {
		return nil, err
	}

	// Choose random delimiters that systems wouldn't use or guess for security.
	startDelim := uuid + "[{[{"
	endDelim := "}]}]" + uuid

	// Change delimiters temporarily so code samples can use Go blocks.
	safeContent := fmt.Sprintf(`%sdefine "content"%s%s%send%s`, startDelim, endDelim, string(content), startDelim, endDelim)
	t = t.Delims(startDelim, endDelim)
	t, err = t.Parse(safeContent)
	if err != nil {
		return nil, err
	}
	// Reset delimiters
	t = t.Delims("{{", "}}")
	return t, nil
}

// generateUUID for use as an random identifier.
func generateUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
