package htmltemplate

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/josephspurrier/ambient/html"
)

// AssetInjector represents code that can inject files into a template.
type AssetInjector interface {
	InjectPlugins(t *template.Template, r *http.Request, pluginNames []string, pageURL string) (*template.Template, error)
}

func (te *Engine) partialTemplate(r *http.Request, mainTemplate string, partialTemplate string) (*template.Template, error) {
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

	t, err = te.assetInjector.InjectPlugins(t, r, te.pluginNames, r.URL.Path)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (te *Engine) postTemplate(r *http.Request, mainTemplate string) (*template.Template, error) {
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

	t, err = te.assetInjector.InjectPlugins(t, r, te.pluginNames, r.URL.Path)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (te *Engine) pluginTemplate(r *http.Request, assets embed.FS, mainTemplate string, partialTemplate string) (*template.Template, error) {
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

	t, err = te.assetInjector.InjectPlugins(t, r, te.pluginNames, r.URL.Path)
	if err != nil {
		return nil, err
	}

	return t, nil
}
