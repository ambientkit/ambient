package ambient

import "net/http"

// AuthenticatedUser returns if the current user is authenticated.
func (ss *SecureSite) AuthenticatedUser(r *http.Request) (string, error) {
	if !ss.Authorized(GrantUserAuthenticatedRead) {
		return "false", ErrAccessDenied
	}

	return ss.sess.AuthenticatedUser(r)
}

// UserLogin sets the current user as authenticated.
func (ss *SecureSite) UserLogin(r *http.Request, username string) error {
	if !ss.Authorized(GrantUserAuthenticatedWrite) {
		return ErrAccessDenied
	}

	ss.sess.Login(r, username)

	return nil
}

// UserPersist sets the user session to retain after browser close.
func (ss *SecureSite) UserPersist(r *http.Request, persist bool) error {
	if !ss.Authorized(GrantUserPersistWrite) {
		return ErrAccessDenied
	}

	ss.sess.Persist(r, persist)

	return nil
}

// UserLogout logs out the current user.
func (ss *SecureSite) UserLogout(r *http.Request) error {
	if !ss.Authorized(GrantUserAuthenticatedWrite) {
		return ErrAccessDenied
	}

	ss.sess.Logout(r)

	return nil
}

// SetCSRF sets the session with a token and returns the token for use in a form.
func (ss *SecureSite) SetCSRF(r *http.Request) string {
	return ss.sess.SetCSRF(r)
}

// CSRF returns true if the CSRF token is valid.
func (ss *SecureSite) CSRF(r *http.Request) bool {
	return ss.sess.CSRF(r)
}
