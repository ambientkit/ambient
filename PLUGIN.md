# Plugin Development Guide <!-- omit in toc -->

This guide will walk you through creating a plugin for Ambient.

- [Minimum Viable Plugin (MVP)](#minimum-viable-plugin-mvp)
- [Types of Plugins](#types-of-plugins)
	- [Logger](#logger)
	- [Storage System](#storage-system)
	- [Session Manager](#session-manager)
	- [Template engine](#template-engine)
	- [Router](#router)
	- [Middleware](#middleware)
	- [Generic Plugin](#generic-plugin)
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
- generic plugin

The main difference between the plugins is what functions are called in them

### Logger

A [logger](plugin/logruslogger/logruslogger.go) is what is used to output messages at different levels: fatal, error, warn, info, and debug. It's helpful when you can provide more information during troubleshooting by changing the log level because you can get to the bottom of issues quicker. The logger is used by the Ambient internal system and is made available 

The logger plugin must include the MVP code as well as the `Logger` function.

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

A [storage system](plugin/gcpbucketstorage/gcpbucketstorage.go) is what is used to store the web app settings (title, content, scheme, URL, etc.) as well as plugin status (enabled/disabled), settings, and permissions granted.



### Session Manager

### Template engine

### Router

### Middleware

### Generic Plugin

## Good Practices

- Use the Ambient logger with all it's different levels: fatal, error, warn, info, and debug. You shouldn't use `log` or `fmt` package to output any messages because they are not standardized.
- If you run background jobs in your plugin, make sure you implement the `Disable()` function to stop the background job.
- When creating a funcmap, you must prefix each one with your plugin name so there are no collisions in the templates. An error message will be throw if any of the funcmaps are not named properly.
- You must return every permission your plugin needs to use in the `GrantRequests()` function. Otherwise, the plugin will not work properly when enabled.

## Misc

- How to use the logger
- How to use the router