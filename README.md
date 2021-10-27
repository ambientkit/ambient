# Ambient 🏖️

[![GitHub Actions status](https://github.com/josephspurrier/ambient/actions/workflows/unit-tests.yml/badge.svg)](https://github.com/josephspurrier/ambient/actions)

## Overview

### What is it?

Ambient is framework in Go for building web apps using plugins. You can use the plugins already included to stand up a blog just like the [Bear Blog](https://bearblog.dev/) or create your own plugins to build your own web app.

### Why was this created?

Each time I write a new web app, I reuse much of the same foundational code. I wrote Ambient to help me standardize existing code, enable/disable packages on demand, modify plugin behaviors using settings, and build new functionality in a reusable way.

### Who is this for?

Ambient will probably appeal to individual developers or small development teams who need to build one or many web apps using the same backend framework. Large teams will probably want a framework more established.

### How does it work?

- Ambient is a web server that accepts an app name, app version, logger, storage system, and a collection of plugins (which must include a session manager, router, template engine).
- Plugins have to satisfy interfaces in order to work with Ambient.
- Plugins must request permissions and the admin must grant each permission.
- Plugins can modify almost any part of a web application:
  - logger
  - session manager
  - router
  - middleware
  - template engine
  - pages or API endpoints
  - content for HTML head, content, navigation, footer, etc.
- Plugin manager allows you to:
  - Enable/disable a plugin
  - Grant permissions to a plugin
  - Modify the settings for a plugin

## Quickstart

You can follow the [tutorial](cmd/myapp/README.md) to quickly get the application up and running locally.

## Screenshots

Below are screens of the sample application that you'll see if you follow the tutorial.

The terminal shows the [logger plugin](plugin/logruslogger/logruslogger.go) that outputs based on log level.

![Terminal](doc/screenshot/terminal.png)

The home screen is from the [simplelogin plugin](plugin/simplelogin/simplelogin.go) and demonstrates the styling from the [bearcss plugin](plugin/bearcss/bearcss.go). Routing is handled through the [awayrouter plugin](plugin/awayrouter).

![Home](doc/screenshot/home.png)

The login page takes a username and password (handled by the [simplelogin plugin](plugin/simplelogin/simplelogin.go)). The password hash is read from the environment variable: `AMB_PASSWORD_HASH`. The [scssession plugin](plugin/scssession/scssession.go) handles the session creation and stores to the local filesystem, but supports any storage system via a plugin that satisfies the [`SessionStorer`](ambient_sessionstorer.go) interface.

![Login](doc/screenshot/login.png)

The [pluginmanager plugin](plugin/pluginmanager/pluginmanager.go) provides an easy way to modify plugins.

![Plugin Manager](doc/screenshot/pluginmanager.png)

The settings page (part of the [pluginmanager plugin](plugin/pluginmanager/pluginmanager.go)) allows you to customize the value that gets displayed in the meta tag that is set by the [author plugin](plugin/author/author.go).

![Settings](doc/screenshot/settings.png)

The grants page (part of the [pluginmanager plugin](plugin/pluginmanager/pluginmanager.go)) allows you to allow or deny modifications to the application by the [author plugin](plugin/author/author.go).

![Grants](doc/screenshot/grants.png)

Once enabled, the [author plugin](plugin/author/author.go) modifies the HTML header (through the [htmlengine plugin](plugin/htmlengine/htmlengine.go)) to add in a meta tag with the value from the settings page.

![HTML](doc/screenshot/htmlauthor.png)

The backend storage is provided by the [gcpbucketstorage plugin](plugin/gcpbucketstorage/gcpbucketstorage.go) and is stored in a JSON file on the local filesystem, but supports any storage system via a plugin that satisfies the [`DataStorer`](ambient_datastorer.go) interface.

![Storage](doc/screenshot/storage.png)