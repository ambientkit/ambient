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
	"github.com/ambientkit/ambient/pkg/grpcp/testutil"
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
	app, err := testutil.Setup(false)
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
	assert.NoError(t, ps.SetGrant("hello", ambient.GrantSitePluginRead))
	assert.NoError(t, ps.SetGrant("hello", ambient.GrantSitePluginDelete))
	assert.NoError(t, ps.SetGrant("hello", ambient.GrantSitePluginEnable))
	assert.NoError(t, ps.SetGrant("hello", ambient.GrantSitePluginDisable))
	assert.NoError(t, ps.SetGrant("hello", ambient.GrantSitePostWrite))
	assert.NoError(t, ps.SetGrant("hello", ambient.GrantSitePostRead))
	assert.NoError(t, ps.SetGrant("hello", ambient.GrantSitePostDelete))
	assert.NoError(t, ps.SetGrant("hello", ambient.GrantPluginNeighborRouteRead))
	assert.NoError(t, ps.SetGrant("hello", ambient.GrantUserPersistWrite))
	assert.NoError(t, ps.SetGrant("hello", ambient.GrantAllUserAuthenticatedWrite))

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

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/neighborPluginRequestedGrant", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Requested: true", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/neighborPluginRequestedGrantBad", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Requested: false", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/plugins", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Plugins: 3", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/pluginNames", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Plugin names: 3", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("DELETE", "/deletePlugin", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Delete plugin: <nil>", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("DELETE", "/deletePluginBad", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Delete plugin: plugin name not found", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/enablePlugin", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Enable plugin: <nil>", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/enablePluginBad", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Enable plugin: item was not found", string(body))

	// resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/loadAllPluginPages", nil))
	// assert.Equal(t, http.StatusOK, resp.StatusCode)
	// assert.Equal(t, "Load pages: <nil>", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/disablePlugin", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Disable plugin: <nil>", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/disablePluginBad", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Disable plugin: item was not found", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("POST", "/savePost", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Posts are the same.", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/publishedPosts", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Posts are the same.", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/publishedPages", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Pages are the same.", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/postBySlug", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Pages are the same.", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/postBySlugBad", nil))
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "item was not found\n", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/postByID", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Pages are the same.", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/postByIDBad", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Pages are the same.", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("DELETE", "/deletePostByID", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Works.", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/pluginNeighborRoutesList", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Routes: 1", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/pluginNeighborRoutesListBad", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Routes: 0", string(body))

	resp, _ = doRequest(t, mux, httptest.NewRequest("GET", "/userPersist", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 86400, resp.Cookies()[0].MaxAge)

	resp, _ = doRequest(t, mux, httptest.NewRequest("GET", "/userPersistFalse", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 0, resp.Cookies()[0].MaxAge)

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/grantRequests", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Grant requests: 18", string(body))

	resp, body = doRequest(t, mux, httptest.NewRequest("GET", "/userLogout", nil))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "User cleared.", string(body))
	assert.Equal(t, "", resp.Cookies()[0].Value)
	assert.Equal(t, -1, resp.Cookies()[0].MaxAge)

	// Login user.
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

	// Destroy users.
	r = httptest.NewRequest("GET", "/logoutAllUsers", nil)
	for _, v := range resp.Cookies() {
		r.AddCookie(v)
	}
	resp, body = doRequest(t, mux, r)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Users cleared.", string(body))

	// Test with authenticated cookie again.
	r = httptest.NewRequest("GET", "/loggedin", nil)
	for _, v := range resp.Cookies() {
		r.AddCookie(v)
	}
	resp, body = doRequest(t, mux, r)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "login: () (user not found)", string(body))
}
