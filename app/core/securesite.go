package core

import (
	"errors"
	"net/http"
	"strings"
	"time"
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

// Title returns the title.
func (ss *SecureSite) Title() (string, error) {
	grant := "site.title:read"

	if !ss.Authorized(grant) {
		return "", ErrAccessDenied
	}

	return ss.storage.Site.Title, nil
}

// SetScheme sets the site scheme.
func (ss *SecureSite) SetScheme(scheme string) error {
	grant := "site.scheme:write"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	ss.storage.Site.Scheme = scheme

	return ss.storage.Save()
}

// Scheme returns the site scheme.
func (ss *SecureSite) Scheme() (string, error) {
	grant := "site.scheme:read"

	if !ss.Authorized(grant) {
		return "", ErrAccessDenied
	}

	return ss.storage.Site.Scheme, nil
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

// SetURL sets the site URL.
func (ss *SecureSite) SetURL(URL string) error {
	grant := "site.url:write"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	ss.storage.Site.URL = URL

	return ss.storage.Save()
}

// URL returns the URL without the scheme at the beginning.
func (ss *SecureSite) URL() (string, error) {
	grant := "site.url:read"

	if !ss.Authorized(grant) {
		return "", ErrAccessDenied
	}

	return ss.storage.Site.URL, nil
}

// FullURL returns the URL with the scheme at the beginning.
func (ss *SecureSite) FullURL() (string, error) {
	grant1 := "site.url:read"
	grant2 := "site.scheme:read"

	if !ss.Authorized(grant1) || !ss.Authorized(grant2) {
		return "", ErrAccessDenied
	}

	return ss.storage.Site.SiteURL(), nil
}

// Plugins returns the plugin list.
func (ss *SecureSite) Plugins() (map[string]PluginSettings, error) {
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
		fields = PluginFields{
			Fields: make(map[string]string),
		}
	}

	fields.Fields[name] = escapeValue(value)
	ss.storage.Site.PluginFields[ss.pluginName] = fields

	return ss.storage.Save()
}

func (ss *SecureSite) pluginField(pluginName string, fieldName string) (string, error) {
	// See if the value is set.
	fields, ok := ss.storage.Site.PluginFields[pluginName]
	if ok {
		value, ok := fields.Fields[fieldName]
		if ok {
			if len(value) > 0 {
				return value, nil
			}
		}
	}

	// See if there is a default value.
	plugin, ok := ss.storage.Site.PluginSettings[pluginName]
	if ok {
		field, ok := plugin.Fields[fieldName]
		if ok {
			if len(field.Default) > 0 {
				return field.Default, nil
			}
		}
	}

	return "", nil
}

// PluginFieldChecked gets a checked variable for the plugin.
func (ss *SecureSite) PluginFieldChecked(name string) (bool, error) {
	grant := "plugin:getfield"

	if !ss.Authorized(grant) {
		return false, ErrAccessDenied
	}

	value, err := ss.pluginField(ss.pluginName, name)

	return value == "true", err
}

// PluginField gets a variable for the plugin.
func (ss *SecureSite) PluginField(fieldName string) (string, error) {
	grant := "plugin:getfield"

	if !ss.Authorized(grant) {
		return "", ErrAccessDenied
	}

	return ss.pluginField(ss.pluginName, fieldName)
}

// SetNeighborPluginField sets a variable for a neighbor plugin.
func (ss *SecureSite) SetNeighborPluginField(pluginName string, fieldName string, value string) error {
	grant := "plugin:setneighborfield"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	fields, ok := ss.storage.Site.PluginFields[pluginName]
	if !ok {
		fields = PluginFields{
			Fields: make(map[string]string),
		}
	}

	fields.Fields[fieldName] = escapeValue(value)
	ss.storage.Site.PluginFields[pluginName] = fields

	return ss.storage.Save()
}

// NeighborPluginField gets a variable for a neighbor plugin.
func (ss *SecureSite) NeighborPluginField(pluginName string, fieldName string) (string, error) {
	grant := "plugin:getneighborfield"

	if !ss.Authorized(grant) {
		return "", ErrAccessDenied
	}

	return ss.pluginField(pluginName, fieldName)
}

// Updated returns the home last updated timestamp.
func (ss *SecureSite) Updated() (time.Time, error) {
	grant := "site.updated:read"

	if !ss.Authorized(grant) {
		return time.Now(), ErrAccessDenied
	}

	return ss.storage.Site.Updated, nil
}

// SavePost saves a post.
func (ss *SecureSite) SavePost(ID string, post Post) error {
	grant := "site.post:write"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	ss.storage.Site.Posts[ID] = post

	return ss.storage.Save()
}

// PostsAndPages returns the list of posts and pages.
func (ss *SecureSite) PostsAndPages(onlyPublished bool) (PostWithIDList, error) {
	grant := "site.postsandpages:read"

	if !ss.Authorized(grant) {
		return nil, ErrAccessDenied
	}

	return ss.storage.Site.PostsAndPages(onlyPublished), nil
}

// PublishedPosts returns the list of published posts.
func (ss *SecureSite) PublishedPosts() ([]Post, error) {
	grant := "site.posts:read" // TODO: Differentiate between posts and published posts?

	if !ss.Authorized(grant) {
		return nil, ErrAccessDenied
	}

	return ss.storage.Site.PublishedPosts(), nil
}

// PublishedPages returns the list of published pages.
func (ss *SecureSite) PublishedPages() ([]Post, error) {
	grant := "site.pages:read" // TODO: Differentiate between posts and published posts?

	if !ss.Authorized(grant) {
		return nil, ErrAccessDenied
	}

	return ss.storage.Site.PublishedPages(), nil
}

// PostBySlug returns the post by slug.
func (ss *SecureSite) PostBySlug(slug string) (PostWithID, error) {
	grant := "site.postbyslug:read"

	if !ss.Authorized(grant) {
		return PostWithID{}, ErrAccessDenied
	}

	return ss.storage.Site.PostBySlug(slug), nil
}

// PostByID returns the post by ID.
func (ss *SecureSite) PostByID(ID string) (Post, error) {
	grant := "site.postbyid:read"

	if !ss.Authorized(grant) {
		return Post{}, ErrAccessDenied
	}

	post, ok := ss.storage.Site.Posts[ID]
	if !ok {
		return Post{}, ErrNotFound
	}

	return post, nil
}

// DeletePostByID deletes a post.
func (ss *SecureSite) DeletePostByID(ID string) error {
	grant := "site.deletepostbyid:write"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	delete(ss.storage.Site.Posts, ID)

	return ss.storage.Save()
}

// Tags returns the list of tags.
func (ss *SecureSite) Tags(onlyPublished bool) (TagList, error) {
	grant := "site.tags:read"

	if !ss.Authorized(grant) {
		return nil, ErrAccessDenied
	}

	return ss.storage.Site.Tags(onlyPublished), nil
}

// SetContent sets the home page content.
func (ss *SecureSite) SetContent(content string) error {
	grant := "site.content:write"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	ss.storage.Site.Content = content

	return ss.storage.Save()
}

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
