// Package websession provides session management.
package websession

import (
	"crypto/rand"
	"errors"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

// Sessionstorer reads and writes data to an object.
type Sessionstorer interface {
	Save([]byte) error
	Load() ([]byte, error)
}

// Session stores session level information
type Session struct {
	Name    string
	manager *scs.SessionManager
}

// New returns a new session cookie store.
func New(name string, manager *scs.SessionManager) *Session {
	return &Session{
		manager: manager,
		Name:    name,
	}
}

// Persist -
func (s *Session) Persist(r *http.Request, persist bool) {
	s.manager.Cookie.Persist = persist
}

// Logout -
func (s *Session) Logout(r *http.Request) {
	s.manager.Destroy(r.Context())
}

// AuthenticatedUser returns the user ID if authenticated or an error.
func (s *Session) AuthenticatedUser(r *http.Request) (string, error) {
	u := s.manager.GetString(r.Context(), "user")

	if len(u) == 0 {
		return "", errors.New("user not found")
	}

	return u, nil
}

// Login -
func (s *Session) Login(r *http.Request, value string) {
	s.manager.Put(r.Context(), "user", value)
}

// String -
func (s *Session) String(r *http.Request, name string) string {
	return s.manager.GetString(r.Context(), name)
}

// SetString -
func (s *Session) SetString(r *http.Request, name string, value string) {
	s.manager.Put(r.Context(), name, value)
}

// SetCSRF -
func (s *Session) SetCSRF(r *http.Request) string {
	token := generate(32)
	path := "csrf_" + r.URL.Path
	s.SetString(r, path, token)
	return token
}

// CSRF -
func (s *Session) CSRF(r *http.Request, token string) bool {
	path := "csrf_" + r.URL.Path
	v := s.String(r, path)

	if len(v) > 0 {
		s.manager.Remove(r.Context(), path)
		if v == token && len(token) > 0 {
			return true
		}
	}

	return false
}

// Generate a token.
// Source: https://devpy.wordpress.com/2013/10/24/create-random-string-in-golang/
func generate(length int) string {
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, length)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}

	return string(bytes)
}
