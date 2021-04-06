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
	grants     map[Grant]bool

	log     IAppLogger
	storage *Storage
	sess    ISession
	mux     IAppRouter
}

// NewSecureSite -
func NewSecureSite(pluginName string, grants map[Grant]bool, log IAppLogger, storage *Storage, session ISession, mux IAppRouter) *SecureSite {
	return &SecureSite{
		pluginName: pluginName,
		grants:     grants,

		log:     log,
		storage: storage,
		sess:    session,
		mux:     mux,
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
	if !ss.Authorized(GrantSiteLoadTrigger) {
		return ErrAccessDenied
	}

	return ss.storage.Load()
}

// Authorized determines if the current context has access.
func (ss *SecureSite) Authorized(grant Grant) bool {
	// If has star, then allow all access.
	if allowed, ok := ss.grants[GrantAll]; ok && allowed {
		return true
	}

	// If the grant was found, then allow access.
	if allowed, ok := ss.grants[grant]; ok && allowed {
		return true
	}

	ss.log.Info("securesite: denied plugin (%v) access to the data item, requires grant: %v\n", ss.pluginName, grant)

	return false
}
