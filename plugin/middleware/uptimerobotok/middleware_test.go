package uptimerobotok_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ambientkit/ambient/plugin/middleware/uptimerobotok"
	"github.com/stretchr/testify/assert"
)

func TestNewSession(t *testing.T) {
	r := httptest.NewRequest("HEAD", "/", nil)
	w := httptest.NewRecorder()
	mux := http.NewServeMux()
	mw := uptimerobotok.HeadReply(mux)
	mw.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}
