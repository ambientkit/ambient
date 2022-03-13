package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	// Set the test relative to the project directory since the plugin path
	// is relative to that.
	path, _ := os.Getwd()
	basePath := strings.TrimSuffix(path, "/cmd/server")
	if err := os.Chdir(basePath); err != nil {
		assert.FailNow(t, err.Error())
	}

	// Set up the application.
	_, pluginClient, mux, err := setup()
	if pluginClient != nil {
		defer pluginClient.Kill()
	}
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	// Test reuqest.
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "hello world", string(body))
}
