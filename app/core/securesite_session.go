package core

import "net/http"

// UserAuthenticated returns if the current user is authenticated.
func (ss *SecureSite) UserAuthenticated(r *http.Request) (bool, error) {
	if !ss.Authorized(GrantUserAuthenticatedRead) {
		return false, ErrAccessDenied
	}

	return ss.sess.UserAuthenticated(r)
}
