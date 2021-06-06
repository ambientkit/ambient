package ambient

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
)

// Toolkit provides utilities to plugins.
type Toolkit struct {
	Log    Logger
	Mux    Router
	Render Renderer
	Site   *SecureSite
}

// Redirect to a page with the proper URL prefix.
func (t *Toolkit) Redirect(w http.ResponseWriter, r *http.Request, url string, code int) {
	http.Redirect(w, r, t.Path(url), code)
}

// Path to a page with the proper URL prefix.
func (t *Toolkit) Path(url string) string {
	return path.Join(os.Getenv("AMB_URL_PREFIX"), url)
}

// JSON sends a JSON response.
func (t *Toolkit) JSON(w http.ResponseWriter, status int, response interface{}) (int, error) {
	// Convert to JSON bytes.
	b, err := json.Marshal(response)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Set the header.
	w.Header().Set("Content-Type", "application/json")

	// Write out the response.
	_, err = fmt.Fprint(w, string(b))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Return the status.
	return status, nil
}

// Renderer represents a template renderer.
type Renderer interface {
	Page(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
	PageContent(w http.ResponseWriter, r *http.Request, content string, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
	Post(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
	PostContent(w http.ResponseWriter, r *http.Request, content string, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
	Error(w http.ResponseWriter, r *http.Request, content string, statusCode int, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
}

// AppRouter represents a router.
type AppRouter interface {
	Router

	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Clear(method string, path string)
	SetNotFound(notFound http.Handler)
	SetServeHTTP(h func(w http.ResponseWriter, r *http.Request, status int, err error))
}

// Router represents a router.
type Router interface {
	Get(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Post(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Patch(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Put(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Delete(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Head(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Options(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Error(status int, w http.ResponseWriter, r *http.Request)
	Param(r *http.Request, name string) string
	Wrap(handler http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) (status int, err error)
}

// AppSession represents a user session.
type AppSession interface {
	UserAuthenticated(r *http.Request) (bool, error)
	Login(r *http.Request, username string)
	Logout(r *http.Request)
	Persist(r *http.Request, persist bool)
	SetCSRF(r *http.Request) string
	CSRF(r *http.Request) bool
}
