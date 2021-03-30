package htmltemplate

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/josephspurrier/ambient/app/lib/templatebuffer"
)

// PluginDashboard renders using the dashboard layout.
func (te *Engine) PluginDashboard(w http.ResponseWriter, r *http.Request, assets embed.FS, partialTemplate string, vars map[string]interface{}) (status int, err error) {
	return te.pluginPartial(w, r, "layout/dashboard", "dashboard", assets, partialTemplate, vars)
}

// PluginPage renders using the page layout.
func (te *Engine) PluginPage(w http.ResponseWriter, r *http.Request, assets embed.FS, partialTemplate string, vars map[string]interface{}) (status int, err error) {
	return te.pluginPartial(w, r, "layout/page", "page", assets, partialTemplate, vars)
}

// PluginPageContent renders using the page content.
func (te *Engine) PluginPageContent(w http.ResponseWriter, r *http.Request, content string, vars map[string]interface{}) (status int, err error) {
	return te.pluginContent(w, r, "layout/page", "page", content, vars)
}

func (te *Engine) pluginPartial(w http.ResponseWriter, r *http.Request, mainTemplate string, layoutType string, assets embed.FS, partialTemplate string, vars map[string]interface{}) (status int, err error) {
	// Set the status to OK starting out.
	status = http.StatusOK

	// Parse the main template with the functions.
	t, err := te.generateTemplate(r, mainTemplate, layoutType)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Parse the plugin template.
	// FIXME: Should parse separately for added security.
	t, err = t.ParseFS(assets, fmt.Sprintf("%v.tmpl", partialTemplate))
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

// pluginContent converts a site post from markdown to HTML and then outputs to response
// writer. Returns an HTTP status code and an error if one occurs.
func (te *Engine) pluginContent(w http.ResponseWriter, r *http.Request, mainTemplate string, layoutType string, postContent string, vars map[string]interface{}) (status int, err error) {
	// Set the status to OK starting out.
	status = http.StatusOK

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
	err = templatebuffer.ParseExistingTemplate(w, t, status, vars)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return
}
