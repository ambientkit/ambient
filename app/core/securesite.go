package core

import (
	"errors"
	"log"
	"net/http"

	"github.com/josephspurrier/ambient/app/lib/datastorage"
	"github.com/josephspurrier/ambient/app/lib/router"
	"github.com/josephspurrier/ambient/app/model"
)

var (
	// ErrAccessDenied is when access is not allowed to the data item.
	ErrAccessDenied = errors.New("access denied to the data item")
	// ErrNotFound is when an item is not found.
	ErrNotFound = errors.New("item was no found")
)

// SecureSite is a secure data access for the site.
type SecureSite struct {
	pluginName string
	storage    *datastorage.Storage
	mux        *router.Mux
	grants     map[string]bool
}

// NewSecureSite -
func NewSecureSite(pluginName string, storage *datastorage.Storage, mux *router.Mux, grants map[string]bool) *SecureSite {
	return &SecureSite{
		pluginName: pluginName,
		storage:    storage,
		mux:        mux,
		grants:     grants,
	}
}

// Error handles returning the proper error.
func (ss *SecureSite) Error(err error) (status int, errr error) {
	switch err {
	case ErrAccessDenied:
		return http.StatusForbidden, err
	case ErrNotFound:
		return http.StatusNotFound, err
	default:
		return http.StatusInternalServerError, err
	}
}

// Authorized determines if the current context has access.
func (ss *SecureSite) Authorized(grant string) bool {
	//return true
	if allowed, ok := ss.grants[grant]; ok && allowed {
		return true
	}

	log.Printf("denied plugin (%v) access to the data action: %v\n", ss.pluginName, grant)

	return false
}

// Title returns the title.
func (ss *SecureSite) Title() (string, error) {
	grant := "site.title:read"

	if !ss.Authorized(grant) {
		return "", ErrAccessDenied
	}

	return ss.storage.Site.Title, nil
}

// SetTitle sets the title.
func (ss *SecureSite) SetTitle(title string) error {
	grant := "site.title:write"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	ss.storage.Site.Title = title

	return ss.storage.Save()
}

// URL returns the site URL.
func (ss *SecureSite) URL() (string, error) {
	grant := "site.url:read"

	if !ss.Authorized(grant) {
		return "", ErrAccessDenied
	}

	return ss.storage.Site.SiteURL(), nil
}

// Plugins returns the plugin list.
func (ss *SecureSite) Plugins() (map[string]model.PluginSettings, error) {
	grant := "site.plugins:read"

	if !ss.Authorized(grant) {
		return nil, ErrAccessDenied
	}

	return ss.storage.Site.PluginSettings, nil
}

// DeletePlugin deletes a plugin.
func (ss *SecureSite) DeletePlugin(name string) error {
	grant := "site.plugins:deleteone"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	delete(ss.storage.Site.PluginSettings, name)

	return ss.storage.Save()
}

// EnablePlugin enables a plugin.
func (ss *SecureSite) EnablePlugin(name string) error {
	grant := "site.plugins:enable"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	plugin, ok := ss.storage.Site.PluginSettings[name]
	if !ok {
		return ErrNotFound
	}

	plugin.Enabled = true
	ss.storage.Site.PluginSettings[name] = plugin

	return ss.storage.Save()
}

// DisablePlugin disables a plugin.
func (ss *SecureSite) DisablePlugin(name string) error {
	grant := "site.plugins:disable"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	plugin, ok := ss.storage.Site.PluginSettings[name]
	if !ok {
		return ErrNotFound
	}

	plugin.Enabled = false
	ss.storage.Site.PluginSettings[name] = plugin

	return ss.storage.Save()
}

// ClearRoute clears out an old route.
func (ss *SecureSite) ClearRoute(method string, path string) error {
	grant := "router:clear"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	ss.mux.Clear(method, path)

	return nil
}

// ClearRoutePlugin clears out an old route.
func (ss *SecureSite) ClearRoutePlugin(pluginName string) error {
	grant := "router:clear"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	routes, ok := ss.storage.PluginRoutes.Routes[pluginName]
	if !ok {
		return ErrNotFound
	}

	for _, v := range routes {
		ss.mux.Clear(v.Method, v.Path)
	}

	return nil
}

// SetPluginField sets a variable for the plugin.
func (ss *SecureSite) SetPluginField(name string, value string) error {
	grant := "plugin:setfield"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	fields, ok := ss.storage.Site.PluginFields[ss.pluginName]
	if !ok {
		fields = model.PluginFields{
			Fields: make(map[string]string),
		}
	}

	fields.Fields[name] = value
	ss.storage.Site.PluginFields[ss.pluginName] = fields

	return ss.storage.Save()
}

// PluginField gets a variable for the plugin.
func (ss *SecureSite) PluginField(name string) (string, error) {
	grant := "plugin:getfield"

	if !ss.Authorized(grant) {
		return "", ErrAccessDenied
	}

	fields, ok := ss.storage.Site.PluginFields[ss.pluginName]
	if !ok {
		return "", ErrNotFound
	}

	value, ok := fields.Fields[name]
	if !ok {
		return "", ErrNotFound
	}

	return value, nil
}

// SetNeighborPluginField sets a variable for a neighbor plugin.
func (ss *SecureSite) SetNeighborPluginField(pluginName string, name string, value string) error {
	grant := "plugin:setneighborfield"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	fields, ok := ss.storage.Site.PluginFields[pluginName]
	if !ok {
		fields = model.PluginFields{
			Fields: make(map[string]string),
		}
	}

	fields.Fields[name] = value
	ss.storage.Site.PluginFields[pluginName] = fields

	return ss.storage.Save()
}

// NeighborPluginField gets a variable for a neighbor plugin.
func (ss *SecureSite) NeighborPluginField(pluginName string, name string) (string, error) {
	grant := "plugin:getneighborfield"

	if !ss.Authorized(grant) {
		return "", ErrAccessDenied
	}

	fields, ok := ss.storage.Site.PluginFields[pluginName]
	if !ok {
		return "", ErrNotFound
	}

	value, ok := fields.Fields[name]
	if !ok {
		return "", ErrNotFound
	}

	return value, nil
}

// ErrorAccessDenied return true if the error is AccessDenied.
func (ss *SecureSite) ErrorAccessDenied(err error) bool {
	return err == ErrAccessDenied
}

// ErrorNotFound return true if the error is NotFound.
func (ss *SecureSite) ErrorNotFound(err error) bool {
	return err == ErrNotFound
}
