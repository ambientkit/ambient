package ambient

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

	log          AppLogger
	pluginsystem *PluginSystem
	sess         AppSession
	mux          AppRouter
	render       Renderer
	recorder     *RouteRecorder
}

// NewSecureSite returns a new secure site.
func NewSecureSite(pluginName string, log AppLogger, ps *PluginSystem, session AppSession, mux AppRouter, render Renderer, recorder *RouteRecorder) *SecureSite {
	return &SecureSite{
		pluginName: pluginName,

		log:          log,
		sess:         session,
		mux:          mux,
		pluginsystem: ps,
		render:       render,
		recorder:     recorder,
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

	return ss.pluginsystem.Load()
}

// Authorized determines if the current context has access.
func (ss *SecureSite) Authorized(grant Grant) bool {
	return ss.pluginsystem.Authorized(ss.pluginName, grant)
}
