package htmltemplate

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/josephspurrier/ambient/app/lib/templatebuffer"
)

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
	t, err := te.generateTemplate(r, mainTemplate, layoutType)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Parse the plugin template separately for security.
	content, err := templatebuffer.ParseTemplateFS(assets, fmt.Sprintf("%v.tmpl", partialTemplate), fm, vars)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	safeContent := fmt.Sprintf(`{{define "content"}}%s{{end}}`, content)
	t, err = t.Parse(safeContent)
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

// pluginContent converts a site post from markdown to HTML and then outputs to response
// writer. Returns an HTTP status code and an error if one occurs.
func (te *Engine) pluginContent(w http.ResponseWriter, r *http.Request, mainTemplate string, layoutType string, postContent string, statusCode int, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
	// Parse the main template with the functions.
	t, err := te.generateTemplate(r, mainTemplate, layoutType)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Parse the plugin template separately for security.
	content, err := templatebuffer.ParseTemplate(te.sanitized(postContent), fm, vars)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	safeContent := fmt.Sprintf(`{{define "content"}}%s{{end}}`, content)
	t, err = t.Parse(safeContent)
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
