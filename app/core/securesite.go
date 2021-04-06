package core

import (
	"errors"
	"net/http"
)

var (
	// ErrAccessDenied is when access is not allowed to the data item.
	ErrAccessDenied = errors.New("access denied to the data item")
	// ErrNotFound is when an item is not found.
	ErrNotFound = errors.New("item was not found")
)

// SecureSite is a secure data access for the site.
type SecureSite struct {
	pluginName string
	storage    *Storage
	sess       ISession
	mux        IAppRouter
	grants     map[string]bool
	log        IAppLogger
}

// NewSecureSite -
func NewSecureSite(pluginName string, log IAppLogger, storage *Storage, session ISession, mux IAppRouter, grants map[string]bool) *SecureSite {
	return &SecureSite{
		pluginName: pluginName,
		storage:    storage,
		sess:       session,
		mux:        mux,
		grants:     grants,
		log:        log,
	}
}

// Error handles returning the proper error.
func (ss *SecureSite) Error(siteError error) (status int, err error) {
	switch siteError {
	case ErrAccessDenied:
		return http.StatusForbidden, siteError
	case ErrNotFound:
		return http.StatusNotFound, siteError
	default:
		return http.StatusInternalServerError, siteError
	}
}

// ErrorAccessDenied return true if the error is AccessDenied.
func (ss *SecureSite) ErrorAccessDenied(err error) bool {
	return err == ErrAccessDenied
}

// ErrorNotFound return true if the error is NotFound.
func (ss *SecureSite) ErrorNotFound(err error) bool {
	return err == ErrNotFound
}

// Load forces a reload of the data.
func (ss *SecureSite) Load() error {
	grant := "site.load:write"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	return ss.storage.Load()
}

// Authorized determines if the current context has access.
func (ss *SecureSite) Authorized(grant string) bool {
	//return true
	if allowed, ok := ss.grants[grant]; ok && allowed {
		return true
	}

	ss.log.Info("securesite: denied plugin (%v) access to the data action: %v\n", ss.pluginName, grant)

	return false
}

// func escapeValue(s string) string {
// 	last := s
// 	before := s
// 	after := ""
// 	for before != after {
// 		before = last
// 		after = last
// 		after = strings.ReplaceAll(after, "{{", "")
// 		after = strings.ReplaceAll(after, "}}", "")
// 		last = after
// 	}
// 	return after
// }
