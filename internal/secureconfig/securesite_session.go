package secureconfig

import (
	"net/http"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/internal/config"
)

// AuthenticatedUser returns if the current user is authenticated.
func (ss *SecureSite) AuthenticatedUser(r *http.Request) (string, error) {
	if !ss.Authorized(ambient.GrantUserAuthenticatedRead) {
		return "", config.ErrAccessDenied
	}

	return ss.sess.AuthenticatedUser(r)
}

// UserLogin sets the current user as authenticated.
func (ss *SecureSite) UserLogin(r *http.Request, username string) error {
	if !ss.Authorized(ambient.GrantUserAuthenticatedWrite) {
		return config.ErrAccessDenied
	}

	ss.sess.Login(r, username)

	return nil
}

// UserPersist sets the user session to retain after browser close.
func (ss *SecureSite) UserPersist(r *http.Request, persist bool) error {
	if !ss.Authorized(ambient.GrantUserPersistWrite) {
		return config.ErrAccessDenied
	}

	ss.sess.Persist(r, persist)

	return nil
}

// UserLogout logs out the current user.
func (ss *SecureSite) UserLogout(r *http.Request) error {
	if !ss.Authorized(ambient.GrantUserAuthenticatedWrite) {
		return config.ErrAccessDenied
	}

	ss.sess.Logout(r)

	return nil
}

// LogoutAllUsers logs out all users.
func (ss *SecureSite) LogoutAllUsers(r *http.Request) error {
	if !ss.Authorized(ambient.GrantAllUserAuthenticatedWrite) {
		return config.ErrAccessDenied
	}

	ss.sess.LogoutAll(r)

	return nil
}

// SetCSRF sets the session with a token and returns the token for use in a form
// or header.
func (ss *SecureSite) SetCSRF(r *http.Request) string {
	return ss.sess.SetCSRF(r)
}

// CSRF returns true if the CSRF token is valid.
func (ss *SecureSite) CSRF(r *http.Request, token string) bool {
	return ss.sess.CSRF(r, token)
}

// SessionValue returns session value by name.
func (ss *SecureSite) SessionValue(r *http.Request, name string) string {
	return ss.sess.SessionValue(r, name)
}

// SetSessionValue sets a value on the current session.
func (ss *SecureSite) SetSessionValue(r *http.Request, name string, value string) error {
	return ss.sess.SetSessionValue(r, name, value)
}

// DeleteSessionValue deletes a session value on the current session.
func (ss *SecureSite) DeleteSessionValue(r *http.Request, name string) {
	ss.sess.DeleteSessionValue(r, name)
}
