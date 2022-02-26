package ambient

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ambientkit/away/router"
	"github.com/stretchr/testify/assert"
)

func TestRouteRecorder(t *testing.T) {
	mp1 := NewMockPlugin("mp1", "1.0.0")
	mp1.MockGrants = []GrantRequest{
		{Grant: GrantRouterRouteWrite, Description: "Access to create default route."},
	}

	// Set up the lighweight app.
	app, _, err := NewApp("myapp", "1.0",
		NewMockLoggerPlugin(),
		StoragePluginGroup{
			Storage: NewMockStoragePlugin(),
		},
		&PluginLoader{
			Router:         nil,
			TemplateEngine: nil,
			SessionManager: nil,
			TrustedPlugins: map[string]bool{
				"mp1": true,
			},
			Plugins: []Plugin{
				mp1,
			},
			Middleware: []MiddlewarePlugin{},
		})
	assert.NoError(t, err)

	mux := router.New()
	rr := NewRouteRecorder(app.log, app.pluginsystem, mux)

	pr1 := rr.withPlugin("mp1")
	called1 := false
	pr1.Get("/", func(http.ResponseWriter, *http.Request) (status int, err error) {
		called1 = true
		return
	})

	pr2 := rr.withPlugin("mp2")
	called2 := false
	pr2.Get("/", func(http.ResponseWriter, *http.Request) (status int, err error) {
		called2 = true
		return
	})

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, called1)
	assert.False(t, called2)
}
