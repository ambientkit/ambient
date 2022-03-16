package grpcp_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/ambientapp"
	"github.com/ambientkit/ambient/pkg/grpcp/testdata/plugin/neighbor"
	"github.com/ambientkit/plugin/logger/zaplogger"
	"github.com/ambientkit/plugin/router/awayrouter"
	"github.com/ambientkit/plugin/sessionmanager/scssession"
	"github.com/ambientkit/plugin/storage/memorystorage"
	"github.com/ambientkit/plugin/templateengine/htmlengine"
	"github.com/stretchr/testify/assert"
)

func grpcSetup(t *testing.T) *ambientapp.App {
	// Set the test relative to the project directory since the plugin path
	// is relative to that.
	path, _ := os.Getwd()
	basePath := strings.TrimSuffix(path, "/pkg/grpcp")
	if err := os.Chdir(basePath); err != nil {
		assert.FailNow(t, err.Error())
	}

	// Set up the application.
	app, err := Setup()
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	return app
}

func doRequest(t *testing.T, mux http.Handler, r *http.Request) (*http.Response, string) {
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	return resp, string(body)
}

func TestMain(t *testing.T) {
	// Setup gRPC server.
	app := grpcSetup(t)
	// Stop plugins when done.
	defer app.StopGRPCClients()

	ps := app.PluginSystem()
	assert.NoError(t, ps.SetEnabled("hello", true))
	assert.NoError(t, ps.SetGrant("hello", ambient.GrantRouterRouteWrite))
	assert.NoError(t, ps.SetGrant("hello", ambient.GrantUserAuthenticatedRead))
	assert.NoError(t, ps.SetGrant("hello", ambient.GrantUserAuthenticatedWrite))
	assert.NoError(t, ps.SetGrant("hello", ambient.GrantPluginNeighborGrantRead))
	assert.NoError(t, ps.SetGrant("hello", ambient.GrantPluginNeighborGrantWrite))

	mux, err := app.Handler()
	if err != nil {
		t.Fatal(err.Error())
	}

	resp, body := doRequest(t, mux, httptest.NewRequest("GET", "/", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "hello world", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/another", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "hello world - another", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/name/foo", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "hello: foo", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/name/bar", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "hello: bar", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/nameold/foo", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "hello: foo", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/error", nil))
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	assert.Equal(t, "Forbidden\n", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/created", nil))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "created: ", string(body))

	r := httptest.NewRequest("GET", "/headers", nil)
	r.Header.Set("foo", "123")
	r.Header.Set("bar", "who")
	resp, body = doRequest(t, mux, r)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "headers: 2", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/form", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "\n<!DOCTYPE html>\n<html lang=\"en\">\n<head></head>\n<body>\n\t<form method=\"post\">\n\t<label for=\"fname\">First name:</label>\n\t<input type=\"text\" id=\"fname\" name=\"fname\" value=\"a\"><br><br>\n\t<label for=\"lname\">Last name:</label>\n\t<input type=\"text\" id=\"lname\" name=\"lname\" value=\"b\"><br><br>\n\t<input type=\"submit\" value=\"Submit\">\n\t</form>\n</body>\n</html>\n", string(body))

	form := url.Values{}
	form.Add("a", "foo")
	form.Add("b", "bar")
	r = httptest.NewRequest("POST", "/form", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, body = doRequest(t, mux, r)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, `body: "a=foo&b=bar"`, string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/loggedin", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "login: () (user not found)", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/login", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "login: (<nil>) (username) (<nil>)", string(body))

	// Test with authenticated cookie.
	r = httptest.NewRequest("GET", "/loggedin", nil)
	for _, v := range resp.Cookies() {
		r.AddCookie(v)
	}
	resp, body = doRequest(t, mux, r)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "login: (username) (<nil>)", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/errors", nil))
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	assert.Equal(t, "request does not exist for the grant\n", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/neighborPluginGrantList", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Grants: 18", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/neighborPluginGrantListBad", nil))
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "item was not found\n", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/neighborPluginGrants", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Grants: 18", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/neighborPluginGranted", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Granted: false", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/neighborPluginGrantedBad", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Granted: false", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/setNeighborPluginGrantFalse", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Granted: false", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/setNeighborPluginGrantTrue", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Granted: true", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/neighborPluginGranted", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Granted: true", string(body))
}

// Setup sets up a test gRPC server.
func Setup() (*ambientapp.App, error) {
	h := func(log ambient.Logger, renderer ambient.Renderer, w http.ResponseWriter, r *http.Request, err error) {
		if err != nil {
			switch e := err.(type) {
			case ambient.Error:
				errText := e.Error()
				if len(errText) == 0 {
					errText = http.StatusText(e.Status())
				}
				http.Error(w, errText, e.Status())
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}
	}

	sessPlugin := scssession.New("5ba3ad678ee1fd9c4fddcef0d45454904422479ed762b3b0ddc990e743cb65e0")
	plugins := &ambient.PluginLoader{
		// Core plugins are implicitly trusted.
		Router:         awayrouter.New(h),
		TemplateEngine: htmlengine.New(),
		SessionManager: sessPlugin,
		// Trusted plugins are those that are typically needed to boot so they
		// will be enabled and given full access.
		TrustedPlugins: map[string]bool{},
		Plugins: []ambient.Plugin{
			neighbor.New(),
		},
		GRPCPlugins: []ambient.GRPCPlugin{
			{Name: "hello", Path: "./pkg/grpcp/testdata/plugin/hello/cmd/plugin/hello"},
		},
		Middleware: []ambient.MiddlewarePlugin{
			// Middleware - executes bottom to top.
			sessPlugin,
		},
	}
	app, _, err := ambientapp.NewApp("myapp", "1.0",
		zaplogger.New(),
		ambient.StoragePluginGroup{
			Storage: memorystorage.New(),
		},
		plugins)
	return app, err
}