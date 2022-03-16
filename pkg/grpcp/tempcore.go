package grpcp

import (
	"net/http"

	"github.com/ambientkit/ambient"
)

// PluginCore represents the core of a plugin.
type PluginCore interface {
	// PluginName should be globally unique. It must start with a lowercase
	// letter and then contain only lowercase letters and numbers.
	PluginName() (string, error)
	// PluginVersion must follow https://semver.org/.
	PluginVersion() (string, error)
	//	GrantRequests() []GrantRequest
	Enable(*Toolkit) error
	Disable() error
	Routes() error
}

// Toolkit provides utilities to plugins.
type Toolkit struct {
	Log  ambient.Logger
	Mux  ambient.Router
	Site SecureSite
}

// SecureSite provides plugin functions.
type SecureSite interface {
	// Error handles returning the proper error.
	Error(siteError error) (err error)
	// Load forces a reload of the data.
	Load() error
	// Authorized determines if the current context has access.
	Authorized(grant ambient.Grant) bool
	// // NeighborPluginGrantList gets the grants requests for a neighbor plugin.
	// NeighborPluginGrantList(pluginName string) ([]GrantRequest, error)
	// // NeighborPluginGrants gets the map of granted permissions.
	// NeighborPluginGrants(pluginName string) (map[Grant]bool, error)
	// // NeighborPluginGranted returns true if the plugin has the grant.
	// NeighborPluginGranted(pluginName string, grantName Grant) (bool, error)
	// // SetNeighborPluginGrant sets a grant for a neighbor plugin.
	// SetNeighborPluginGrant(pluginName string, grantName Grant, granted bool) error
	// // Plugins returns the plugin list.
	// Plugins() (map[string]PluginData, error)
	// // PluginNames returns the list of plugin name.
	// PluginNames() ([]string, error)
	// // DeletePlugin deletes a plugin.
	// DeletePlugin(name string) error
	// // EnablePlugin enables a plugin.
	// EnablePlugin(pluginName string, loadPlugin bool) error
	// // LoadAllPluginPages loads all of the pages from the plugins.
	// LoadAllPluginPages() error
	// // DisablePlugin disables a plugin.
	// DisablePlugin(pluginName string, unloadPlugin bool) error
	// // LoadAllPluginMiddleware returns a handler that is wrapped in conditional
	// // middlware from the plugins. This only needs to be run once at start up
	// // and should never be called again.
	// LoadAllPluginMiddleware() http.Handler
	// // SavePost saves a post.
	// SavePost(ID string, post Post) error
	// // PostsAndPages returns the list of posts and pages.
	// PostsAndPages(onlyPublished bool) (PostWithIDList, error)
	// // PublishedPosts returns the list of published posts.
	// PublishedPosts() ([]Post, error)
	// // PublishedPages returns the list of published pages.
	// PublishedPages() ([]Post, error)
	// // PostBySlug returns the post by slug.
	// PostBySlug(slug string) (PostWithID, error)
	// // PostByID returns the post by ID.
	// PostByID(ID string) (Post, error)
	// // DeletePostByID deletes a post.
	// DeletePostByID(ID string) error
	// // PluginNeighborRoutesList gets the routes for a neighbor plugin.
	// PluginNeighborRoutesList(pluginName string) ([]Route, error)
	// AuthenticatedUser returns if the current user is authenticated.
	AuthenticatedUser(r *http.Request) (string, error)
	// UserLogin sets the current user as authenticated.
	UserLogin(r *http.Request, username string) error
	// // UserPersist sets the user session to retain after browser close.
	// UserPersist(r *http.Request, persist bool) error
	// // UserLogout logs out the current user.
	// UserLogout(r *http.Request) error
	// // LogoutAllUsers logs out all users.
	// LogoutAllUsers(r *http.Request) error
	// // SetCSRF sets the session with a token and returns the token for use in a form
	// // or header.
	// SetCSRF(r *http.Request) string
	// // CSRF returns true if the CSRF token is valid.
	// CSRF(r *http.Request, token string) bool
	// // SessionValue returns session value by name.
	// SessionValue(r *http.Request, name string) string
	// // SetSessionValue sets a value on the current session.
	// SetSessionValue(r *http.Request, name string, value string) error
	// // DeleteSessionValue deletes a session value on the current session.
	// DeleteSessionValue(r *http.Request, name string)
	// // PluginNeighborSettingsList gets the grants requests for a neighbor plugin.
	// PluginNeighborSettingsList(pluginName string) ([]Setting, error)
	// // SetPluginSetting sets a variable for the plugin.
	// SetPluginSetting(settingName string, value string) error
	// // PluginSettingBool returns a plugin setting as a bool.
	// PluginSettingBool(name string) (bool, error)
	// // PluginSettingString returns a setting for the plugin as a string.
	// PluginSettingString(fieldName string) (string, error)
	// // PluginSetting returns a setting for the plugin as an interface{}.
	// PluginSetting(fieldName string) (interface{}, error)
	// // SetNeighborPluginSetting sets a setting for a neighbor plugin.
	// SetNeighborPluginSetting(pluginName string, settingName string, value string) error
	// // NeighborPluginSettingString returns a setting for a neighbor plugin as a string.
	// NeighborPluginSettingString(pluginName string, fieldName string) (string, error)
	// // NeighborPluginSetting returns a setting for a neighbor plugin as an interface{}.
	// NeighborPluginSetting(pluginName string, fieldName string) (interface{}, error)
	// // PluginTrusted returns whether a plugin is trusted or not.
	// PluginTrusted(pluginName string) (bool, error)
	// // SetTitle sets the title.
	// SetTitle(title string) error
	// // Title returns the title.
	// Title() (string, error)
	// // SetScheme sets the site scheme.
	// SetScheme(scheme string) error
	// // Scheme returns the site scheme.
	// Scheme() (string, error)
	// // SetURL sets the site URL.
	// SetURL(URL string) error
	// // URL returns the URL without the scheme at the beginning.
	// URL() (string, error)
	// // FullURL returns the URL with the scheme at the beginning.
	// FullURL() (string, error)
	// // Updated returns the home last updated timestamp.
	// Updated() (time.Time, error)
	// // Tags returns the list of tags.
	// Tags(onlyPublished bool) (TagList, error)
	// // SetContent sets the home page content.
	// SetContent(content string) error
	// // Content returns the site home page content.
	// Content() (string, error)
}