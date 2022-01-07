# Plugin Development Guide <!-- omit in toc -->

This guide will walk you through creating a plugin for Ambient.

- [Minimum Viable Plugin (MVP)](#minimum-viable-plugin-mvp)
- [Plugin Boot Process](#plugin-boot-process)
- [Plugin Functions](#plugin-functions)
	- [Logger](#logger)
	- [Storage System](#storage-system)
	- [Session Manager](#session-manager)
	- [Template Engine](#template-engine)
	- [Router](#router)
	- [Middleware](#middleware)
	- [Routes](#routes)
	- [Grant Requests](#grant-requests)
	- [Settings](#settings)
	- [Assets](#assets)
	- [Funcmaps](#funcmaps)
- [Good Plugin Practices](#good-plugin-practices)
- [FAQs](#faqs)

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

	// Define the session manager so it can be used as a core plugin and
	// middleware.
	sessionManager := scssession.New(secretKey)

	return &ambient.PluginLoader{
		// Core plugins are implicitly trusted.
		Router:         awayrouter.New(nil),
		TemplateEngine: htmlengine.New(),
		SessionManager: sessionManager,
		// Trusted plugins are required to boot the app so they will be
		// given full access.
		TrustedPlugins: map[string]bool{
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
			sessionManager,   // Session manager middleware.
			logrequest.New(), // Log every request as INFO.
		},
	}
}
```

When you start the app, the plugin will not be enabled. You must login and then navigate to: http://localhost:8080/dashboard/plugins. Put a checkmark next to the plugin and then click the **Save** button at the bottom of the page. Your plugin is now enabled! It doesn't do anything so don't get too excited, but you've got the scaffolding of a plugin so on to the next step where you choose the type of plugin function to use.

## Plugin Boot Process

An Ambient app follows this process when it boots:

- Load logger plugin by calling `Logger()` func.
- Load storage plugin by calling `Storage()` func.
- Enable and grant permissions to trusted plugins.
- Load session manager plugin by calling `SessionManager()` func.
- Load template engine plugin by calling `TemplateEngine()` func.
- Load router plugin by calling `Router()` func.
- Load each plugin (except those above) only if enabled in `site.bin` file:
  - Enable plugin by calling `Enable()` and passing in `ambient.Toolkit`.
  - Add routes from plugin to router by calling `Routes()`.
  - Load assets from plugin by calling `Assets()`.
- Load each middleware plugin (first handler is the handler from router):
  - Wrap the middleware around the previous handler by calling `Middleware()` func and adding a conditional so it's only run if enabled in `site.bin` file.
- Pass the handler to `ListenAndServe()` func.

An Ambient app can have a plugin enabled or disabled while it's running (through the pluginmanager). It will then load or unload the plugin making changes to routes, assets, and middleware.

When a change to the app is made or data is read or modified in `site.bin` file, the permissions of the plugin are checked first to ensure the user granted the plugin permissions to perform their action. The permissions are stored in the `site.bin` file.

A few things to note:
- Logger plugin and storage plugin are automatically trusted because they are loaded before the plugin system boots.
- Router plugin and template engine plugin are automatically trusted because they are explicitly passed to the plugin system.
- Logger, storage, template engine, and router won't have the `Enable()` func called so it will only be able to use parts of the toolkit that are passed in when their respective functions are called. You can also remove the `*ambient.PluginBase` and `*ambient.Toolkit` from the main struct since they won't be used. You can see [zaplogger](plugin/logger/zaplogger/zaplogger.go) as an example.
- Session manager should always have a middleware component to it so shouldn't be listed in the Plugins section, but it should be listed in the Middleware section. Be sure to define it only once and then use it as both a parameter for `ambient.PluginLoader.SessionManager` and `ambient.PluginLoader.Middleware`. You define it in middleware so you can control when it gets called relative to other middleware.
- Plugin manager should be in the trusted plugins list since it's required to enable other plugins.

## Plugin Functions

Ambient supports the following types of plugin functions via the [Plugin interface](ambient.go):

- logger
- storage system
- session manager
- template engine
- router
- middleware
- routes
- grant requests
- settings
- assets
- funcmaps

There are also plugins that fall outside this list (most of them). They use the remainder of the functions to modify or interact with the app. Since a single interface is used for all plugins, a single plugin could essentially serve all the purposes, but then it wouldn't reall

The main difference between the plugins is what functions are called in them.

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

The function should return an object that satisfies the [`AppLogger`](ambient_logger.go) interface. You should probably also add in an option to output in either human readable format (tabs) or JSON to make it easy to work with in development or in production.

```go
package ambient

// AppLogger represents the log service for the app.
type AppLogger interface {
	Logger

	// Fatal is reserved for the app level only.
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
	// app and it needs to be corrected.
	LogLevelError
	// LogLevelFatal is for messages when the app cannot continue and
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

The function should return objects that satisfy the [`DataStorer`](ambient_datastorer.go) interface and the [`SessionStorer`](ambient_sessionstorer.go) interface. Notice you don't have to worry about the type of data. This makes it easy to read or write to any medium.

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

The `SessionManager()` function should return an object that satisfies the [`AppSession`](ambient_session.go) interface. The `Middleware()` function should return an object that satisfies the `[]func(next http.Handler) http.Handler` definition.

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

The function should return an object that satisfies the [`Renderer`](ambient_renderer.go) interface. The page and post allow you to define two different formats when rendering content so you can have the assets affect each differently.

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

The function should return an object that satisfies the [`AppRouter`](ambient_router.go) interface. Note the router needs to support clearing routes which may require extending popular router packages.

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

The `Middleware()` function should return an object that satisfies the `[]func(next http.Handler) http.Handler` definition. The middleware load from bottom to top so be sure to organize them accordingly. They will also be ordered based on the position in the `plugin.go` file.

### Routes

The `Routes()` function registers HTTP handlers with the router.

A [plugin with routes](plugin/simplelogin/simplelogin.go) defined must include the MVP code as well as the `Routes()` function.

```go
// Routes gets routes for the plugin.
func (p *Plugin) Routes() {
	p.Mux.Get("/", p.home)
	p.Mux.Get("/login", p.login)
	p.Mux.Post("/login", p.loginPost)
	p.Mux.Get("/dashboard", p.dashboard)
	p.Mux.Get("/dashboard/logout", p.logout)
}
```

The function doesn't return any objects and shouldn't fail either. It also takes a special kind of HTTP handler - one that returns an HTTP status code and an error. A standard HTTP handler doesn't have any returns, but that makes it more difficult to standardize how to write out status codes and output errors so this new function definition improves on it.

```go
// Home renders the home template.
func (p *Plugin) Home(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "Home"
	return p.Render.Page(w, r, assets, "template/content/home", p.funcMap(r), vars)
}
```

### Grant Requests

The `GrantRequests()` function returns a list of permissions required by the plugin. The admin of the app must enable each of the permissions.

A [plugin](plugin/prism/prism.go) that needs to make changes to the app or interact with its data must include the MVP code as well as the `GrantRequests()` function.

```go
// GrantRequests returns a list of grants requested by the plugin.
func (p *Plugin) GrantRequests() []ambient.GrantRequest {
	return []ambient.GrantRequest{
		{Grant: ambient.GrantSiteAssetWrite, Description: "Access to add stylesheets and javascript to each page."},
		{Grant: ambient.GrantRouterRouteWrite, Description: "Access to create routes for accessing stylesheets."},
		{Grant: ambient.GrantPluginSettingRead, Description: "Read own plugin settings."},
	}
}
```

The function returns a `[]GrantRequest` object. You can see the full list of permissions in [model_grant.go](model_grant.go).

### Settings

The `Settings()` function returns a list of settings that can be edited from the pluginmanager UI.

A [plugin](plugin/simplelogin/simplelogin.go) that has configurable settings should use MVP code as well as the `Settings()` function.

```go
// Settings returns a list of user settable fields.
func (p *Plugin) Settings() []ambient.Setting {
	return []ambient.Setting{
		{
			Name:    Username,
			Default: "admin",
		},
		{
			Name:    Password,
			Default: p.passwordHash,
			Type:    ambient.InputPassword,
			Hide:    true,
		},
		{
			Name: MFAKey,
			Type: ambient.InputPassword,
		},
		{
			Name:    LoginURL,
			Default: "admin",
			Hide:    true,
		},
		{
			Name: Author,
		},
		{
			Name: Subtitle,
			Hide: true,
		},
		{
			Name: Description,
			Type: ambient.Textarea,
		},
		{
			Name: Footer,
			Type: ambient.Textarea,
			Hide: true,
		},
		{
			Name: AllowHTMLinMarkdown,
			Type: ambient.Checkbox,
		},
	}
}
```

You can see all of the available setting types in the [model_setting.go](model_setting.go) file.

### Assets

The `Assets()` function returns a list of assets that can modify the template output. They can be used to link to local resources like stylesheets and javascript files or they can link to external resources. You can also define the header and footer for template output. Assets support templating from a string or from a template. You can even add HTML tags to the template output. They are pretty powerful and support a bunch of use cases.

A [plugin](plugin/simplelogin/simplelogin.go) that has assets should use MVP code as well as the `Assets()` function.

```go
// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]ambient.Asset, *embed.FS) {
	arr := make([]ambient.Asset, 0)

	siteTitle, err := p.Site.Title()
	if err == nil && len(siteTitle) > 0 {
		arr = append(arr, ambient.Asset{
			Filetype: ambient.AssetGeneric,
			Location: ambient.LocationHead,
			TagName:  "title",
			Inline:   true,
			Content:  fmt.Sprintf(`{{if .pagetitle}}{{.pagetitle}} | %v{{else}}%v{{end}}`, siteTitle, siteTitle),
		})
	}

	siteDescription, err := p.Site.PluginSettingString(Description)
	if err == nil && len(siteDescription) > 0 {
		arr = append(arr, ambient.Asset{
			Filetype:   ambient.AssetGeneric,
			Location:   ambient.LocationHead,
			TagName:    "meta",
			ClosingTag: false,
			Attributes: []ambient.Attribute{
				{
					Name:  "name",
					Value: "description",
				},
				{
					Name:  "content",
					Value: fmt.Sprintf("{{if .pagedescription}}{{.pagedescription}}{{else}}%v{{end}}", siteDescription),
				},
			},
		})
	}

	siteAuthor, err := p.Site.PluginSettingString(Author)
	if err == nil && len(siteAuthor) > 0 {
		arr = append(arr, ambient.Asset{
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
					Value: siteAuthor,
				},
			},
		})
	}

	arr = append(arr, ambient.Asset{
		Path:     "template/partial/nav.tmpl",
		Filetype: ambient.AssetGeneric,
		Location: ambient.LocationHeader,
		Inline:   true,
	})

	arr = append(arr, ambient.Asset{
		Path:     "template/partial/footer.tmpl",
		Filetype: ambient.AssetGeneric,
		Location: ambient.LocationFooter,
		Inline:   true,
	})

	return arr, &assets
}
```

You can see all of the available setting types in the [asset.go](asset.go) file.

### Funcmaps

The `FuncMap()` function returns a `template.FuncMap` that can be used in the templates. They can also be used in assets which is pretty cool.

A [plugin](plugin/disqus/diqus.go) that needs a FuncMap for templates should use MVP code as well as the `Assets()` function.

```go
// FuncMap returns a callable function when passed in a request.
func (p *Plugin) FuncMap() func(r *http.Request) template.FuncMap {
	return func(r *http.Request) template.FuncMap {
		fm := make(template.FuncMap)
		fm["disqus_PageURL"] = func() string {
			return r.URL.Path
		}

		return fm
	}
}
```

## Good Plugin Practices

- Use the Ambient logger with all it's different levels: fatal, error, warn, info, and debug. You shouldn't use `log` or `fmt` package to output any messages because they are not standardized.
- If you run background jobs in your plugin, make sure you implement the `Disable()` function to stop the background job.
- When creating a funcmap, you must prefix each one with your plugin name so there are no collisions in the templates. An error message will be throw if any of the funcmaps are not named properly.
- You must return every permission your plugin needs to use in the `GrantRequests()` function. Otherwise, the plugin will not work properly when enabled.
- Routes can be overwritten by other plugins so it's best to namespace with the plugin name since plugin names are unique.

## FAQs

**What should determine the plugin boundary?**

You should group required items together in a plugin. For instance, a session manager should implement both the `SessionManager()` and `Middleware()` functions - it wouldn't make sense to separate into different plugins since they require each other.

You could group by feature. A login plugin will have multiple routes (`GET /login`, `POST /login`, `GET /logout`) so those should all be in the same plugin.

You could group by type as well. You may want to keep most of your middleware or code that modifies your header in the same plugin. For middleware, you may want to group these into a single plugin so you can use a single settings page to enable or disable each one: request logging, trail slash removal, and response gzip compression.

**How should I approach wrapping a package as a plugin?**

Regardless if you're using Ambient or not, Go packages should be designed so they can be reused by other people and projects without tightly coupling dependencies. If you are creating a new plugin, build the package so it follows [Go best practices](https://talks.golang.org/2013/bestpractices.slide) and then import it to the plugin package. That way you can use it with other apps as well.

**Should I be worred about creating too many plugins?**

If you feel like you have so many plugins that it's hard to find what you're looking for, it's probably too many plugins.

**Does the app run slower with more plugins?**

It will require a big more startup time with many plugins, but it should be neligible unless you have hundreds of plugins.