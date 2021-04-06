package htmltemplate

import (
	"html/template"

	"github.com/josephspurrier/ambient/plugin/htmltemplate/lib/templatebuffer"
)

func (te *Engine) inject(t *template.Template, field string, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error) {
	body, err := templatebuffer.ParseTemplate(content, fm, data)
	if err != nil {
		return nil, err
	}

	// Escape the content.
	t, err = te.escapeContent(t, field, body)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Head -
func (te *Engine) Head(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error) {
	return te.inject(t, "PluginHeadContent", content, fm, data)
}

// Header -
func (te *Engine) Header(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error) {
	return te.inject(t, "PluginHeaderContent", content, fm, data)
}

// Main -
func (te *Engine) Main(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error) {
	return te.inject(t, "PluginMainContent", content, fm, data)
}

// Footer -
func (te *Engine) Footer(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error) {
	return te.inject(t, "PluginFooterContent", content, fm, data)
}

// Body -
func (te *Engine) Body(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error) {
	return te.inject(t, "PluginBodyContent", content, fm, data)
}
