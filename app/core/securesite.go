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
	// ErrGrantNotRequested is when a grant is attempted to be enabled on a
	// plugin, but the plugin didn't explicitly request the grant so don't allow
	// it.
	ErrGrantNotRequested = errors.New("request does not exist for the grant")
	// ErrSettingNotSpecified is when a setting is attempted to be set on a
	// plugin, but the plugin didn't explicity specify it as a setting.
	ErrSettingNotSpecified = errors.New("setting does not exist for the plugin")
)

// SecureSite is a secure data access for the site.
type SecureSite struct {
	pluginName string

	log          IAppLogger
	storage      *Storage
	sess         ISession
	mux          IAppRouter
	pluginsystem *PluginSystem
	render       IRender
}

// NewSecureSite -
func NewSecureSite(pluginName string, log IAppLogger, storage *Storage, session ISession, mux IAppRouter, render IRender, ps *PluginSystem) *SecureSite {
	return &SecureSite{
		pluginName: pluginName,

		log:          log,
		storage:      storage,
		sess:         session,
		mux:          mux,
		pluginsystem: ps,
		render:       render,
	}
}

// Error handles returning the proper error.
func (ss *SecureSite) Error(siteError error) (status int, err error) {
	switch siteError {
	case ErrAccessDenied, ErrGrantNotRequested, ErrSettingNotSpecified:
		return http.StatusForbidden, siteError
	case ErrNotFound:
		return http.StatusNotFound, siteError
	default:
		return http.StatusInternalServerError, siteError
	}
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
	// Always allow ambient application to get full access.
	if ss.pluginName == "ambient" {
		ss.log.Debug("securesite: granted plugin (%v) GrantAll access to the data item for grant: %v\n", "ambient", grant)
		return true
	}

	// If has star, then allow all access.
	if granted := ss.pluginsystem.Granted(ss.pluginName, GrantAll); granted {
		ss.log.Debug("securesite: granted plugin (%v) GrantAll access to the data item for grant: %v\n", ss.pluginName, grant)
		return true
	}

	// If the grant was found, then allow access.
	if granted := ss.pluginsystem.Granted(ss.pluginName, grant); granted {
		ss.log.Debug("securesite: granted plugin (%v) access to the data item for grant: %v\n", ss.pluginName, grant)
		return true
	}

	ss.log.Warn("securesite: denied plugin (%v) access to the data item, requires grant: %v\n", ss.pluginName, grant)

	return false
}
