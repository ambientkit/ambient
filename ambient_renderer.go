package ambient

import (
	"embed"
	"html/template"
	"net/http"
)

// Renderer represents a template renderer.
type Renderer interface {
	Page(w http.ResponseWriter, r *http.Request, assets FileSystemReader, templateName string,
		fm func(r *http.Request) template.FuncMap, vars map[string]interface{}) (err error)
	PageContent(w http.ResponseWriter, r *http.Request, content string,
		fm func(r *http.Request) template.FuncMap, vars map[string]interface{}) (err error)
	Post(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string,
		fm func(r *http.Request) template.FuncMap, vars map[string]interface{}) (err error)
	PostContent(w http.ResponseWriter, r *http.Request, content string,
		fm func(r *http.Request) template.FuncMap, vars map[string]interface{}) (err error)
	Error(w http.ResponseWriter, r *http.Request, content string, statusCode int,
		fm func(r *http.Request) template.FuncMap, vars map[string]interface{}) (err error)
}
