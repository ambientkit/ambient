package core

import (
	"embed"
	"html/template"
	"net/http"
)

// TemplateManager represents a function map for templates.
type TemplateManager interface {
	FuncMap(r *http.Request) template.FuncMap
	Templates() *embed.FS
}
