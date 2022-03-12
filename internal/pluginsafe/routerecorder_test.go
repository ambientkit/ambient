package pluginsafe_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/internal/mock"
	"github.com/ambientkit/ambient/internal/pluginsafe"
	"github.com/ambientkit/ambient/pkg/ambientapp"
	"github.com/ambientkit/away/router"
	"github.com/stretchr/testify/assert"
)

func TestRouteRecorder(t *testing.T) {
	mp1 := mock.NewPlugin("mp1", "1.0.0")
	mp1.MockGrants = []ambient.GrantRequest{
		{Grant: ambient.GrantRouterRouteWrite, Description: "Access to create default route."},
	}

	mp2 := mock.NewPlugin("mp2", "1.0.0")
	mp2.MockGrants = []ambient.GrantRequest{
		{Grant: ambient.GrantRouterRouteWrite, Description: "Access to create default route."},
	}

	// Set up the lighweight app.
	app, logger, err := ambientapp.NewApp("myapp", "1.0",
		mock.NewLoggerPlugin(nil),
		ambient.StoragePluginGroup{
			Storage: mock.NewStoragePlugin(),
		},
		&ambient.PluginLoader{
			Router:         nil,
			TemplateEngine: nil,
			SessionManager: nil,
			TrustedPlugins: map[string]bool{
				"mp1": true,
				"mp2": true,
			},
			Plugins: []ambient.Plugin{
				mp1,
				mp2,
			},
			Middleware: []ambient.MiddlewarePlugin{},
		})
	assert.NoError(t, err)

	ps := app.PluginSystem()

	mux := router.New()
	rr := pluginsafe.NewRouteRecorder(logger, ps, mux)

	pr1 := rr.WithPlugin("mp1")
	called1 := false
	pr1.Get("/", func(http.ResponseWriter, *http.Request) (err error) {
		called1 = true
		return
	})

	pr2 := rr.WithPlugin("mp2")
	called2 := false
	pr2.Get("/", func(http.ResponseWriter, *http.Request) (err error) {
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

	err = ps.SetEnabled("mp1", false)
	assert.NoError(t, err)

	called1 = false
	called2 = false

	r = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	resp = w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.False(t, called1)
	assert.True(t, called2)
}
