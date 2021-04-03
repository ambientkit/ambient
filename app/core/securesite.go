package core

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/josephspurrier/ambient/app/lib/datastorage"
	"github.com/josephspurrier/ambient/app/lib/logger"
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
	sess       ISession
	mux        *router.Mux
	grants     map[string]bool
	log        *logger.Logger
}

// NewSecureSite -
func NewSecureSite(pluginName string, log *logger.Logger, storage *datastorage.Storage, session ISession, mux *router.Mux, grants map[string]bool) *SecureSite {
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

func escapeValue(s string) string {
	last := s
	before := s
	after := ""
	for before != after {
		before = last
		after = last
		after = strings.ReplaceAll(after, "{{", "")
		after = strings.ReplaceAll(after, "}}", "")
		last = after
	}
	return after
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

	ss.storage.Site.Title = escapeValue(title)

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

	fields.Fields[name] = escapeValue(value)
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

	fields.Fields[name] = escapeValue(value)
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

// Updated returns the home last updated timestamp.
func (ss *SecureSite) Updated() (time.Time, error) {
	grant := "site.updated:read"

	if !ss.Authorized(grant) {
		return time.Now(), ErrAccessDenied
	}

	return ss.storage.Site.Updated, nil
}

// PostsAndPages returns the list of posts and pages.
func (ss *SecureSite) PostsAndPages(onlyPublished bool) (model.PostWithIDList, error) {
	grant := "site.postsandpages:read"

	if !ss.Authorized(grant) {
		return nil, ErrAccessDenied
	}

	return ss.storage.Site.PostsAndPages(onlyPublished), nil
}

// Tags returns the list of tags.
func (ss *SecureSite) Tags(onlyPublished bool) (model.TagList, error) {
	grant := "site.tags:read"

	if !ss.Authorized(grant) {
		return nil, ErrAccessDenied
	}

	return ss.storage.Site.Tags(onlyPublished), nil
}

// Description returns the site description.
// func (ss *SecureSite) Description() (string, error) {
// 	grant := "site.description:read"

// 	if !ss.Authorized(grant) {
// 		return "", ErrAccessDenied
// 	}

// 	return ss.storage.Site.Description, nil
// }

// Content returns the site home page content.
func (ss *SecureSite) Content() (string, error) {
	grant := "site.content:read"

	if !ss.Authorized(grant) {
		return "", ErrAccessDenied
	}

	return ss.storage.Site.Content, nil
}

// UserAuthenticated returns if the current user is authenticated.
func (ss *SecureSite) UserAuthenticated(r *http.Request) (bool, error) {
	grant := "user.authenticated:read"

	if !ss.Authorized(grant) {
		return false, ErrAccessDenied
	}

	return ss.sess.UserAuthenticated(r)
}
