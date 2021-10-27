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

## Quickstart on Local

To test out an example web app:

- Clone the repository: `git clone git@github.com:josephspurrier/ambient.git`
- Create a new file called `.env` in the root of the repository with this content:

```bash
# Local Development
## Set this to any value to allow you to do testing locally without GCP access.
## See 'Local Development Flag' section below for more information.
AMB_LOCAL=true

# App Configuration
## Session key to encrypt the cookie store. Generate with: make privatekey
AMB_SESSION_KEY=
## Password hash that is base64 encoded. Generate with: make passhash passwordhere
AMB_PASSWORD_HASH=

# GCP Deployment
## GCP project ID.
AMB_GCP_PROJECT_ID=my-sample-project-191923
## GCP bucket name (this can be one that doesn't exist yet).
AMB_GCP_BUCKET_NAME=sample-bucket
## Name of the docker image that will be created and stored in GCP Repository.
AMB_GCP_IMAGE_NAME=sample-image
## Name of the Cloud Run service to create.
AMB_GCP_CLOUDRUN_NAME=sample-service
## Region (not zone) where the Cloud Run service will be created:
## https://cloud.google.com/compute/docs/regions-zones#available
AMB_GCP_REGION=us-central1

## Optional: set the time zone from here:
## https://golang.org/src/time/zoneinfo_abbrs_windows.go
# AMB_TIMEZONE=America/New_York
```

- To create the session and site files in the storage folder, run: `make local-init`
- To start the webserver on port 8080, run: `make local-run`

The login page is located at: http://localhost:8080/login.

To login, you'll need:

- the default username is: `admin`
- the password from the .env file for which the `AMB_PASSWORD_HASH` was derived

Once you are logged in, you should see a new menu option call `Plugins`. From this screen, you'll be able to use the Plugin Manager to make changes to the plugin state, permissions, and settings.

### Local Development Flag

When `AMB_LOCAL` is set, the following things will happen:

- data storage will be the local filesystem instead of in Google Cloud Storage
- if you try to access the application, it will listen on all IPs/addresses, instead of redirecting like it does in production

## Development Workflow

If you would like to make changes to the code, I recommend these tools to help streamline your workflow.

```bash
# Install air to allow hot reloading so you can make changes quickly.
curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s

# Install direnv and hook into your shell. This allows you to manage 
# https://direnv.net/docs/installation.html
```

Once you have `direnv` installed, create .envrc file:

```bash
# Load the shared environment variables (shared with Makefile).
# Export the vars in .env into the shell.
export $(egrep -v '^#' .env | xargs)

export PATH=$PATH:$(pwd)/bin
```

You can then use this command to start the web server and monitor for changes:

```bash
# Start hot reload. The web application should be available at: http://localhost:8080
air
```