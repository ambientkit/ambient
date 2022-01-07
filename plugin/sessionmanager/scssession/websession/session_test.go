package websession_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/josephspurrier/ambient/lib/aesdata"
	"github.com/josephspurrier/ambient/plugin/sessionmanager/scssession/websession"
	"github.com/josephspurrier/ambient/plugin/storage/localstorage/store"
	"github.com/stretchr/testify/assert"
)

func TestNewSession(t *testing.T) {
	// Set up the session storage provider.
	f := "data.bin"
	err := ioutil.WriteFile(f, []byte(""), 0644)
	assert.NoError(t, err)
	ss := store.NewLocalStorage(f)
	secretkey := "82a18fbbfed2694bb15d512a70c53b1a088e669966918d3d474564b2ac44349b"
	en := aesdata.NewEncryptedStorage(secretkey)
	store, err := websession.NewJSONSession(ss, en)
	assert.NoError(t, err)

	// Initialize a new session manager and configure the session lifetime.
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = false
	sessionManager.Store = store
	sess := websession.New("session", sessionManager)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Test user
		u := "foo"
		sess.Login(r, u)
		user, err := sess.AuthenticatedUser(r)
		assert.True(t, err == nil)
		assert.Equal(t, u, user)

		// Test Logout
		sess.Logout(r)
		_, err = sess.AuthenticatedUser(r)
		assert.False(t, err == nil)

		// Test persistence
		assert.Equal(t, sessionManager.Cookie.Persist, false)
		sess.Persist(r, true)
		assert.Equal(t, sessionManager.Cookie.Persist, true)

		// Test CSRF
		assert.False(t, sess.CSRF(r))
		token := sess.SetCSRF(r)
		r.Form = url.Values{}
		r.Form.Set("token", token)
		assert.True(t, sess.CSRF(r))
	})

	mw := sessionManager.LoadAndSave(mux)
	mw.ServeHTTP(w, r)

	os.Remove(f)
}
