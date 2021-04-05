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

// InjectHead -
func (te *Engine) InjectHead(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error) {
	return te.inject(t, "PluginHeadContent", content, fm, data)
}

// InjectHeader -
func (te *Engine) InjectHeader(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error) {
	return te.inject(t, "PluginHeaderContent", content, fm, data)
}

// InjectMain -
func (te *Engine) InjectMain(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error) {
	return te.inject(t, "PluginMainContent", content, fm, data)
}

// InjectFooter -
func (te *Engine) InjectFooter(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error) {
	return te.inject(t, "PluginFooterContent", content, fm, data)
}

// InjectBody -
func (te *Engine) InjectBody(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error) {
	return te.inject(t, "PluginBodyContent", content, fm, data)
}
