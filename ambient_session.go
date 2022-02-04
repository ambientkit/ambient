package ambient

import (
	"net/http"
)

// AppSession represents a user session.
type AppSession interface {
	AuthenticatedUser(r *http.Request) (string, error)
	Login(r *http.Request, username string)
	Logout(r *http.Request)
	LogoutAll(r *http.Request) error
	Persist(r *http.Request, persist bool)
	SetCSRF(r *http.Request) string
	CSRF(r *http.Request, token string) bool
	SessionValue(r *http.Request, name string) string
	SetSessionValue(r *http.Request, name string, value string) error
	DeleteSessionValue(r *http.Request, name string)
}
