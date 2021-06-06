package ambient

import (
	"net/http"
)

// AppSession represents a user session.
type AppSession interface {
	UserAuthenticated(r *http.Request) (bool, error)
	Login(r *http.Request, username string)
	Logout(r *http.Request)
	Persist(r *http.Request, persist bool)
	SetCSRF(r *http.Request) string
	CSRF(r *http.Request) bool
}
