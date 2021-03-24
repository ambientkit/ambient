package htmltemplate

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/josephspurrier/ambient/html"
)

// HeaderInjector represents code that can inject files into a template.
type HeaderInjector interface {
	InjectPlugins(t *template.Template, r *http.Request) (*template.Template, error)
}

// PartialTemplate -
func (te *Engine) PartialTemplate(r *http.Request, mainTemplate string, partialTemplate string) (*template.Template, error) {
	// Functions available in the templates.
	fm := html.FuncMap(r, te.storage, te.sess)

	baseTemplate := fmt.Sprintf("%v.tmpl", mainTemplate)
	headTemplate := "partial/head.tmpl"
	headerTemplate := "partial/header.tmpl"
	navTemplate := "partial/nav.tmpl"
	footerTemplate := "partial/footer.tmpl"
	contentTemplate := fmt.Sprintf("content/%v.tmpl", partialTemplate)

	// Parse the main template with the functions.
	t, err := template.New(path.Base(baseTemplate)).Funcs(fm).ParseFS(html.Templates, baseTemplate,
		headTemplate, headerTemplate, navTemplate, footerTemplate, contentTemplate)
	if err != nil {
		return nil, err
	}

	t, err = te.hi.InjectPlugins(t, r)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// PostTemplate -
func (te *Engine) PostTemplate(r *http.Request, mainTemplate string) (*template.Template, error) {
	// Functions available in the templates.
	fm := html.FuncMap(r, te.storage, te.sess)

	baseTemplate := fmt.Sprintf("%v.tmpl", mainTemplate)
	headTemplate := "partial/head.tmpl"
	headerTemplate := "partial/header.tmpl"
	navTemplate := "partial/nav.tmpl"
	footerTemplate := "partial/footer.tmpl"

	// Parse the main template with the functions.
	t, err := template.New(path.Base(baseTemplate)).Funcs(fm).ParseFS(html.Templates, baseTemplate, headTemplate, headerTemplate, navTemplate, footerTemplate)
	if err != nil {
		return nil, err
	}

	t, err = te.hi.InjectPlugins(t, r)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// PluginTemplate2 -
func (te *Engine) PluginTemplate2(r *http.Request, assets embed.FS, mainTemplate string, partialTemplate string) (*template.Template, error) {
	// Functions available in the templates.
	fm := html.FuncMap(r, te.storage, te.sess)

	baseTemplate := fmt.Sprintf("%v.tmpl", mainTemplate)
	headTemplate := "partial/head.tmpl"
	headerTemplate := "partial/header.tmpl"
	navTemplate := "partial/nav.tmpl"
	footerTemplate := "partial/footer.tmpl"

	// Parse the main template with the functions.
	t, err := template.New(path.Base(baseTemplate)).Funcs(fm).ParseFS(html.Templates, baseTemplate, headTemplate, headerTemplate, navTemplate, footerTemplate)
	if err != nil {
		return nil, err
	}

	// Parse the plugin template.
	t, err = t.ParseFS(assets, partialTemplate)
	if err != nil {
		return nil, err
	}

	t, err = te.hi.InjectPlugins(t, r)
	if err != nil {
		return nil, err
	}

	return t, nil
}
