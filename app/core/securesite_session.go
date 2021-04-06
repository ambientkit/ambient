package core

import "net/http"

// UserAuthenticated returns if the current user is authenticated.
func (ss *SecureSite) UserAuthenticated(r *http.Request) (bool, error) {
	grant := "user.authenticated:read"

	if !ss.Authorized(grant) {
		return false, ErrAccessDenied
	}

	return ss.sess.UserAuthenticated(r)
}
