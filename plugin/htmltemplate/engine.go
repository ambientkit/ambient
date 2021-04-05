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

func (te *Engine) generateTemplate(r *http.Request, mainTemplate string, layoutType string) (*template.Template, error) {
	// Functions available in the templates.
	fm := te.templateManager.FuncMap(r)

	// Generate list of templates.
	baseTemplate := fmt.Sprintf("%v.tmpl", mainTemplate)
	templates := []string{
		baseTemplate,
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
