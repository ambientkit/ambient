# Plugin Development Guide <!-- omit in toc -->

This guide will walk you through creating a plugin for Ambient.

- [Minimum Viable Plugin (MVP)](#minimum-viable-plugin-mvp)
- [Types of Plugins](#types-of-plugins)
	- [Logger](#logger)
	- [Storage System](#storage-system)
	- [Session Manager](#session-manager)
	- [Template Engine](#template-engine)
	- [Router](#router)
	- [Middleware](#middleware)
	- [Other Plugin](#other-plugin)
- [Good Practices](#good-practices)
- [Misc](#misc)

## Minimum Viable Plugin (MVP)

To create the smallest package that can be used as a plugin, you can paste this code into a new file like this: `plugin/mvp/mvp.go`:

```go
// Package mvp provides a template for building a plugin for Ambient apps.
package mvp

import (
	"github.com/josephspurrier/ambient"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	// PluginBase is a struct that provides empty functions to satisfy the
	// `Plugin` interface in ambient.go. This allows plugins to work with newer
	// versions of the interface as long as functions defined below still match
	// the interface.
	*ambient.PluginBase
	// Toolkit is an object that should be assigned when the plugin is enabled
	// so the plugin can interact with the Logger, Router, Renderer, and
	// SecureSite. If the plugin tries to interact with the toolkit before the
	// plugin is enabled, it will be nil.
	*ambient.Toolkit
}

// New returns a new mvp plugin. This should be modified to include any values
// or objects that need to be passed in before it's enabled. Any example would
// be a password hash or a flag for debug mode.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
	}
}

// PluginName returns the plugin name. This name should match the package name.
// PluginName should be globally unique. Only lowercase letters, numbers,
// and underscores are permitted. Must start with with a letter.
func (p *Plugin) PluginName() string {
	return "mvp"
}

// PluginVersion returns the plugin version. This version must follow
// https://semver.org/.
func (p *Plugin) PluginVersion() string {
	return "1.0.0"
}

// Enable accepts the toolkit and stores it for use when enabled.
func (p *Plugin) Enable(toolkit *ambient.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}
```

To add the plugin to your Ambient app, you should add `mvp.New(),` to the `plugin.go` file, under the `Plugins` section array:

```go
// Plugins defines the plugins.
func Plugins() *ambient.PluginLoader {
	// Get the environment variables.
	secretKey := os.Getenv("AMB_SESSION_KEY")
	if len(secretKey) == 0 {
		log.Fatalf("app: environment variable missing: %v\n", "AMB_SESSION_KEY")
	}

	passwordHash := os.Getenv("AMB_PASSWORD_HASH")
	if len(passwordHash) == 0 {
		log.Fatalf("app: environment variable is missing: %v\n", "AMB_PASSWORD_HASH")
	}

	return &ambient.PluginLoader{
		Router:         awayrouter.New(nil),
		TemplateEngine: htmlengine.New(),
		// Trusted plugins are required to boot the application so they will be
		// given full access.
		TrustedPlugins: map[string]bool{
			"scssession":    true, // Session manager.
			"pluginmanager": true, // Page to manage plugins.
			"simplelogin":   true, // Simple login page.
			"bearcss":       true, // Bear Blog styling.
		},
		Plugins: []ambient.Plugin{
			pluginmanager.New(),           // Page to manage plugins.
			simplelogin.New(passwordHash), // Simple login page.
			bearcss.New(),                 // Bear Blog styling.
			mvp.New(),                     // Your new plugin.
		},
		Middleware: []ambient.MiddlewarePlugin{
			// Middleware - executes bottom to top.
			scssession.New(secretKey), // Session manager.
			logrequest.New(),          // Log every request as INFO.
		},
	}
}
```

When you start the application, the plugin will not be enabled. You must login and then navigate to: http://localhost:8080/dashboard/plugins. Put a checkmark next to the plugin and then click the **Save** button at the bottom of the page. Your plugin is now enabled! It doesn't do anything so don't get too excited, but you've got the scaffolding of a plugin so on to the next step where you choose the type of plugin to create.

## Types of Plugins

Ambient supports the following types of plugins:

- logger
- storage system
- session manager
- template engine
- router
- middleware
- other plugin

The main difference between the plugins is what functions are called in them. All the functions listed in the "other plugin" type can be used TBD

### Logger

A [logger](plugin/logruslogger/logruslogger.go) outputs messages at different levels: fatal, error, warn, info, and debug. It's helpful when you can provide more information during troubleshooting by changing the log level because you can get to the bottom of issues quicker. The logger is used by the Ambient internal system and is made available to plugins as well.

The logger plugin must include the MVP code as well as the `Logger()` function.

```go
// Logger returns a logger.
func (p *Plugin) Logger(appName string, appVersion string) (ambient.AppLogger, error) {
	// Create the logger.
	p.log = NewLogger(appName, appVersion)
	return p.log, nil
}
```

The function should return an object that satisfies the [`AppLogger`](ambient_logger.go) interface.

```go
package ambient

// AppLogger represents the log service for the application.
type AppLogger interface {
	Logger

	// Fatal is reserved for the application level only.
	Fatal(format string, v ...interface{})
	SetLogLevel(level LogLevel)
}

// Logger represents the log service for the plugins.
type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
}

// LogLevel is a log level.
type LogLevel int

const (
	// LogLevelDebug is for debugging output. It's very verbose.
	LogLevelDebug LogLevel = iota
	// LogLevelInfo is for informational messages. It shows messages on services
	// starting, stopping, and users logging in.
	LogLevelInfo
	// LogLevelWarn is for behavior that may need to be fixed. It shows
	// permission warnings for plugins.
	LogLevelWarn
	// LogLevelError is for messages when something is wrong with the
	// application and it needs to be corrected.
	LogLevelError
	// LogLevelFatal is for messages when the application cannot continue and
	// will halt.
	LogLevelFatal
)
```

### Storage System

A [storage system](plugin/gcpbucketstorage/gcpbucketstorage.go) stores the web app settings (title, content, scheme, URL, etc.) as well as plugin status (enabled/disabled), settings, and permissions granted.

The storage system plugin must include the MVP code as well as the `Storage()` function.

```go
// Storage returns data and session storage.
func (p *Plugin) Storage(logger ambient.Logger) (ambient.DataStorer, ambient.SessionStorer, error) {
	var ds ambient.DataStorer
	var ss ambient.SessionStorer

	if envdetect.RunningLocalDev() {
		// Use local filesytem when developing.
		ds = store.NewLocalStorage(p.sitePath)
		ss = store.NewLocalStorage(p.sessionPath)
	} else {
		bucket, err := p.Site.PluginSettingString(Bucket)
		if err != nil {
			return nil, nil, err
		}

		// Use Google when running in GCP.
		ds = store.NewGCPStorage(bucket, p.sitePath)
		ss = store.NewGCPStorage(bucket, p.sessionPath)
	}

	return ds, ss, nil
}
```

The function should return objects that satisfy the [`DataStorer`](ambient_datastorer.go) interface and the [`SessionStorer`](ambient_sessionstorer.go) interface.

```go
// DataStorer reads and writes data to an object.
type DataStorer interface {
	Save([]byte) error
	Load() ([]byte, error)
}

// SessionStorer reads and writes data to an object.
type SessionStorer interface {
	Save([]byte) error
	Load() ([]byte, error)
}

```

### Session Manager

A [session manager](plugin/scssession/scssession.go) authenticates and verify users.

The session manager plugin must include the MVP code as well as the `SessionManager()` function. It should also include a `Middleware()` function to verify users when they try to access authenticated routes. You can get more information in the the [Middleware](#middleware) section below.

```go
// SessionManager returns the session manager.
func (p *Plugin) SessionManager(logger ambient.Logger, ss ambient.SessionStorer) (ambient.AppSession, error) {
	// Set up the session storage provider.
	en := websession.NewEncryptedStorage(p.sessionKey)
	store, err := websession.NewJSONSession(ss, en)
	if err != nil {
		return nil, err
	}

	sessionName := "session"

	p.sessionManager = scs.New()
	p.sessionManager.Lifetime = 24 * time.Hour
	p.sessionManager.Cookie.Persist = false
	p.sessionManager.Store = store
	p.sess = websession.New(sessionName, p.sessionManager)

	return p.sess, nil
}

// Middleware returns router middleware.
func (p *Plugin) Middleware() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{
		p.sessionManager.LoadAndSave,
	}
}
```

The `SessionManager()` function should return an object that satisfies the [`AppSession`](ambient_session.go) interface. The `Middleware()` function should return an object that satisfies the `http.Handler` interface.

```go
// AppSession represents a user session.
type AppSession interface {
	AuthenticatedUser(r *http.Request) (string, error)
	Login(r *http.Request, username string)
	Logout(r *http.Request)
	Persist(r *http.Request, persist bool)
	SetCSRF(r *http.Request) string
	CSRF(r *http.Request) bool
}

// Handler from the http standard library package.
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

### Template Engine

A [template engine](plugin/htmlengine/htmlengine.go) renders content to the `ResponseWriter`.

The template engine plugin must include the MVP code as well as the `TemplateEngine()` function.

```go
// TemplateEngine returns a template engine.
func (p *Plugin) TemplateEngine(logger ambient.Logger, injector ambient.AssetInjector) (ambient.Renderer, error) {
	tmpl := NewTemplateEngine(logger, injector)
	return tmpl, nil
}
```

The function should return an object that satisfies the [`Renderer`](ambient_renderer.go) interface.

```go
// Renderer represents a template renderer.
type Renderer interface {
	Page(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string,
		fm template.FuncMap, vars map[string]interface{}) (status int, err error)
	PageContent(w http.ResponseWriter, r *http.Request, content string,
		fm template.FuncMap, vars map[string]interface{}) (status int, err error)
	Post(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string,
		fm template.FuncMap, vars map[string]interface{}) (status int, err error)
	PostContent(w http.ResponseWriter, r *http.Request, content string,
		fm template.FuncMap, vars map[string]interface{}) (status int, err error)
	Error(w http.ResponseWriter, r *http.Request, content string, statusCode int,
		fm template.FuncMap, vars map[string]interface{}) (status int, err error)
}
```

### Router

A [router](plugin/awayrouter/awayrouter.go) handles web requests based on the HTTP method and route.

A router plugin must include the MVP code as well as the `Router()` function.

```go
// Router returns a router.
func (p *Plugin) Router(logger ambient.Logger, te ambient.Renderer) (ambient.AppRouter, error) {
	// Set up the default router.
	mux := router.New()

	// Set the NotFound and custom ServeHTTP handlers.
	p.setupRouter(logger, mux, te)

	return mux, nil
}
```

The function should return an object that satisfies the [`AppRouter`](ambient_router.go) interface.

```go
// AppRouter represents a router.
type AppRouter interface {
	Router

	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Clear(method string, path string)
	SetNotFound(notFound http.Handler)
	SetServeHTTP(h func(w http.ResponseWriter, r *http.Request, status int, err error))
}

// Router represents a router.
type Router interface {
	Get(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Post(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Patch(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Put(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Delete(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Head(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Options(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Error(status int, w http.ResponseWriter, r *http.Request)
	Param(r *http.Request, name string) string
	Wrap(handler http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) (status int, err error)
}
```

### Middleware

[Middleware](plugin/scssession/scssession.go) runs code before or after a route is served. It's useful for tasks like loging a request, checking if the user is authenticated, and compressing the response.

A middleware plugin must include the MVP code as well as the `Middleware()` function.

```go
// Middleware returns router middleware.
func (p *Plugin) Middleware() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{
		p.sessionManager.LoadAndSave,
	}
}
```

The `Middleware()` function should return an object that satisfies the `http.Handler` interface.

```go
// Handler from the http standard library package.
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

### Other Plugin

An [other plugin](plugin/author/author.go) is a plugin that doesn't fall under any of the other categories.

Any plugin that modifies the application must use the `GrantRequests()` function and return the `[]GrantRequest` object. This does apply to any of the plugin types above.

To add configurable settings, the plugin must use the `Setting()` function and return the `[]Setting` object. To modify the HTML of responses, the plugin can use the `Assets()` function and return a `[]Asset` object and an `embed.FS` object.

```go
// GrantRequests returns a list of grants requested by the plugin.
func (p *Plugin) GrantRequests() []ambient.GrantRequest {
	return []ambient.GrantRequest{
		{Grant: ambient.GrantPluginSettingRead, Description: "Access to the author name."},
		{Grant: ambient.GrantSiteAssetWrite, Description: "Access to write a meta tag to the header."},
	}
}

// Settings returns a list of user settable fields.
func (p *Plugin) Settings() []ambient.Setting {
	return []ambient.Setting{
		{
			Name: Author,
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]ambient.Asset, *embed.FS) {
	name, err := p.Site.PluginSettingString(Author)
	if err != nil || len(name) == 0 {
		// Otherwise don't set the assets.
		return nil, nil
	}

	return []ambient.Asset{
		{
			Filetype:   ambient.AssetGeneric,
			Location:   ambient.LocationHead,
			TagName:    "meta",
			ClosingTag: false,
			Attributes: []ambient.Attribute{
				{
					Name:  "name",
					Value: "author",
				},
				{
					Name:  "content",
					Value: name,
				},
			},
		},
	}, nil
}
```

## Good Practices

- Use the Ambient logger with all it's different levels: fatal, error, warn, info, and debug. You shouldn't use `log` or `fmt` package to output any messages because they are not standardized.
- If you run background jobs in your plugin, make sure you implement the `Disable()` function to stop the background job.
- When creating a funcmap, you must prefix each one with your plugin name so there are no collisions in the templates. An error message will be throw if any of the funcmaps are not named properly.
- You must return every permission your plugin needs to use in the `GrantRequests()` function. Otherwise, the plugin will not work properly when enabled.

## Misc

- How to use the logger
- How to use the router