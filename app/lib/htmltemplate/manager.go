package htmltemplate

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/josephspurrier/ambient/html"
)

type HeaderInjector interface {
	InjectPlugins(t *template.Template) (*template.Template, error)
}

// PartialTemplate -
func (tm *Engine) PartialTemplate(r *http.Request, mainTemplate string, partialTemplate string) (*template.Template, error) {
	// Functions available in the templates.
	fm := html.FuncMap(r, tm.storage, tm.sess)

	baseTemplate := fmt.Sprintf("%v.tmpl", mainTemplate)
	headerTemplate := "partial/head.tmpl"
	contentTemplate := fmt.Sprintf("partial/%v.tmpl", partialTemplate)

	// Parse the main template with the functions.
	t, err := template.New(path.Base(baseTemplate)).Funcs(fm).ParseFS(html.Templates, baseTemplate,
		headerTemplate, contentTemplate)
	if err != nil {
		return nil, err
	}

	t, err = tm.hi.InjectPlugins(t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// PostTemplate -
func (tm *Engine) PostTemplate(r *http.Request, mainTemplate string) (*template.Template, error) {
	// Functions available in the templates.
	fm := html.FuncMap(r, tm.storage, tm.sess)

	baseTemplate := fmt.Sprintf("%v.tmpl", mainTemplate)
	headerTemplate := "partial/head.tmpl"

	// Parse the main template with the functions.
	t, err := template.New(path.Base(baseTemplate)).Funcs(fm).ParseFS(html.Templates, baseTemplate, headerTemplate)
	if err != nil {
		return nil, err
	}

	t, err = tm.hi.InjectPlugins(t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// PluginTemplate -
func (tm *Engine) PluginTemplate2(r *http.Request, assets embed.FS, mainTemplate string, partialTemplate string) (*template.Template, error) {
	// Functions available in the templates.
	fm := html.FuncMap(r, tm.storage, tm.sess)

	baseTemplate := fmt.Sprintf("%v.tmpl", mainTemplate)
	headerTemplate := "partial/head.tmpl"

	// Parse the main template with the functions.
	t, err := template.New(path.Base(baseTemplate)).Funcs(fm).ParseFS(html.Templates, baseTemplate,
		headerTemplate)
	if err != nil {
		return nil, err
	}

	// Parse the plugin template.
	t, err = t.ParseFS(assets, partialTemplate)
	if err != nil {
		return nil, err
	}

	t, err = tm.hi.InjectPlugins(t)
	if err != nil {
		return nil, err
	}

	return t, nil
}
