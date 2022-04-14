package ambient

import (
	"context"
	"html/template"
	"net/http"
)

// AssetInjector represents code that can inject files into a template.
type AssetInjector interface {
	Inject(ctx context.Context, injector LayoutInjector, t *template.Template, r *http.Request, layoutType LayoutType, vars map[string]interface{}) (*template.Template, error)
	DebugTemplates() bool
	EscapeTemplates() bool
}

// LayoutInjector represents an injector that the AssetInjector will call to inject assets in the correct place.
type LayoutInjector interface {
	Head(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error)
	Header(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error)
	Main(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error)
	Footer(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error)
	Body(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error)
}
