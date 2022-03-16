package grpcp_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/ambientkit/ambient/internal/testutil"
	"github.com/ambientkit/ambient/pkg/grpcp"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/assert"
)

func grpcSetup(t *testing.T) (grpcp.PluginCore, *plugin.Client, http.Handler) {
	// Set the test relative to the project directory since the plugin path
	// is relative to that.
	path, _ := os.Getwd()
	basePath := strings.TrimSuffix(path, "/pkg/grpcp")
	if err := os.Chdir(basePath); err != nil {
		assert.FailNow(t, err.Error())
	}

	// Set up the application.
	core, pluginClient, mux, err := testutil.Setup()
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	return core, pluginClient, mux
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
	_, pluginClient, mux := grpcSetup(t)
	if pluginClient != nil {
		defer pluginClient.Kill()
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
	assert.Equal(t, "Granted: true", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/neighborPluginGrantedBad", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Granted: false", string(body))
}
