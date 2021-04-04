package htmltemplate

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/josephspurrier/ambient/app/lib/templatebuffer"
)

// PluginDashboard renders using the dashboard layout.
func (te *Engine) PluginDashboard(w http.ResponseWriter, r *http.Request, assets embed.FS, partialTemplate string, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
	return te.pluginPartial(w, r, "layout/dashboard", "dashboard", assets, partialTemplate, fm, vars)
}

// PluginPage renders using the page layout.
func (te *Engine) PluginPage(w http.ResponseWriter, r *http.Request, assets embed.FS, partialTemplate string, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
	return te.pluginPartial(w, r, "layout/page", "page", assets, partialTemplate, fm, vars)
}

// PluginPageContent renders using the page content.
func (te *Engine) PluginPageContent(w http.ResponseWriter, r *http.Request, content string, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
	return te.pluginContent(w, r, "layout/page", "page", content, fm, vars)
}

func (te *Engine) pluginPartial(w http.ResponseWriter, r *http.Request, mainTemplate string, layoutType string, assets embed.FS, partialTemplate string, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
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
	err = templatebuffer.ParseExistingTemplate(w, r, t, http.StatusOK, vars)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return
}

// pluginContent converts a site post from markdown to HTML and then outputs to response
// writer. Returns an HTTP status code and an error if one occurs.
func (te *Engine) pluginContent(w http.ResponseWriter, r *http.Request, mainTemplate string, layoutType string, postContent string, fm template.FuncMap, vars map[string]interface{}) (status int, err error) {
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
	err = templatebuffer.ParseExistingTemplate(w, r, t, http.StatusOK, vars)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return
}
