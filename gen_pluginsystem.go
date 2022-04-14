// Code generated by ifacemaker. DO NOT EDIT.

package ambient

import (
	"context"
	"time"
)

// PluginSystem provides config functions.
type PluginSystem interface {
	// LoaderPlugins returns the loader plugins, these include initial gRPC plugins
	// as well.
	LoaderPlugins() []Plugin
	// LoaderMiddleware returns the loader middleware, these include initial gRPC plugins
	// as well.
	LoaderMiddleware() []MiddlewarePlugin
	// Plugins returns the map of plugins. This returns a map that is passed by
	// reference and all the values are pointers so any changes to it will
	// be reflected in the plugin system.
	Plugins() map[string]Plugin
	// SessionManager returns the session manager.
	SessionManager() SessionManagerPlugin
	// TemplateEngine returns the template engine.
	TemplateEngine() TemplateEnginePlugin
	// Router returns the router.
	Router() RouterPlugin
	// StorageManager returns the storage manager.
	StorageManager() Storage
	// LoadPlugin loads a single plugin into the plugin system and saves the config.
	LoadPlugin(ctx context.Context, plugin Plugin, middleware bool, grpcPlugin bool) (err error)
	// Load will load the storage and return an error if one occurs.
	Load() error
	// Save will save the storage and return an error if one occurs.
	Save() error
	// InitializePlugin will initialize the plugin in the storage and will return
	// an error if one occurs.
	InitializePlugin(pluginName string, pluginVersion string) error
	// RemovePlugin will delete the plugin from the storage and will return
	// an error if one occurs.
	RemovePlugin(pluginName string) error
	// Names returns a list of plugin names.
	Names() []string
	// MiddlewareNames returns a list of middleware plugin names.
	MiddlewareNames() []string
	// IsMiddleware returns if the plugin is middleware.
	IsMiddleware(name string) bool
	// TrustedPluginNames returns a list of sorted trusted names.
	TrustedPluginNames() []string
	// Trusted returns if a plugin is trusted.
	Trusted(pluginName string) bool
	// Routes returns a list of plugin routes.
	Routes(pluginName string) []Route
	// PluginsData returns the plugin data map.
	PluginsData() map[string]PluginData
	// Plugin returns a plugin by name.
	Plugin(name string) (Plugin, error)
	// PluginData returns a plugin data by name.
	PluginData(name string) (PluginData, error)
	// Enabled returns if the plugin is enabled or not. If it cannot be found, it
	// will still return false.
	Enabled(name string) bool
	// SetEnabled sets a plugin as enabled or not.
	SetEnabled(pluginName string, enabled bool) error
	// GrantRequests returns a list of grant requests.
	GrantRequests(pluginName string, grant Grant) ([]GrantRequest, error)
	// Authorized returns whether a plugin is inherited granted for a plugin.
	Authorized(pluginName string, grant Grant) bool
	// Granted returns whether a plugin is explicitly granted for a plugin.
	Granted(pluginName string, grant Grant) bool
	// SetGrant sets a plugin grant.
	SetGrant(pluginName string, grant Grant) error
	// RemoveGrant removes a plugin grant.
	RemoveGrant(pluginName string, grant Grant) error
	// SetSetting sets a plugin setting.
	SetSetting(pluginName string, settingName string, value interface{}) error
	// Setting returns a setting value.
	Setting(pluginName string, settingName string) (interface{}, error)
	// SettingDefault returns a setting default for a setting.
	SettingDefault(pluginName string, settingName string) (interface{}, error)
	// SetRoute saves a route.
	SetRoute(pluginName string, route []Route)
	// SetTitle sets the title.
	SetTitle(title string) error
	// Title returns the title.
	Title() string
	// SetScheme sets the site scheme.
	SetScheme(scheme string) error
	// Scheme returns the site scheme.
	Scheme() string
	// SetURL sets the site URL.
	SetURL(URL string) error
	// URL returns the URL without the scheme at the beginning.
	URL() string
	// FullURL returns the URL with the scheme at the beginning.
	FullURL() string
	// Updated returns the home last updated timestamp.
	Updated() time.Time
	// Tags returns the list of tags.
	Tags(onlyPublished bool) TagList
	// SetContent sets the home page content.
	SetContent(content string) error
	// Content returns the site home page content.
	Content() string
	// SavePost saves a post.
	SavePost(ID string, post Post) error
	// PostsAndPages returns the list of posts and pages.
	PostsAndPages(onlyPublished bool) PostWithIDList
	// PublishedPosts returns the list of published posts.
	PublishedPosts() []Post
	// PublishedPages returns the list of published pages.
	PublishedPages() []Post
	// PostBySlug returns the post by slug.
	PostBySlug(slug string) PostWithID
	// PostByID returns the post by ID.
	PostByID(ID string) (Post, error)
	// DeletePostByID deletes a post.
	DeletePostByID(ID string) error
}
