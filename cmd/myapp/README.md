# Sample Application Tutorial

This folder contains a sample application to demonstrate how to use Ambient.

## Quickstart on Local

To test out the sample web app from the `cmd/myapp` folder

- Clone the repository: `git clone git@github.com:josephspurrier/ambient.git`
- Create a new file called `.env` in the root of the repository with this content:

```bash
# Local Development
## Set this to any value to allow you to do testing locally without GCP access.
## See 'Local Development Flags' section below for more information.
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

## Optional: set the web server port.
# PORT=8080

## Optional: set the URL prefix if behind a proxy.
# AMB_URL_PREFIX=/api
```

- To create the session and site files in the storage folder, run: `make local-init`
- To start the webserver on port 8080, run: `make local-run`

The login page is located at: http://localhost:8080/login.

To login, you'll need:

- the default username is: `admin`
- the password from the .env file for which the `AMB_PASSWORD_HASH` was derived

Once you are logged in, you should see a new menu option call `Plugins`. From this screen, you'll be able to use the Plugin Manager to make changes to the plugin state, permissions, and settings.

### Local Development Flags

You can set the web server `PORT` to values other than `8080`.

When `AMB_LOCAL` is set to `true`:

- data storage will be the local filesystem instead of in Google Cloud Storage
- if you try to access the application, it will listen on all IPs/addresses, instead of redirecting like it does in production

You can use `envdetect.RunningLocalDev()` to detect if the flag is set to true or not.

When `AMB_TIMEZONE` is set to a timezone like `America/New_York`, the application will use that timezone. This is required if using time-based packages like MFA.

When `AMB_URL_PREFIX` is set to a path like `/api`, the application will server requests from `/api/...`. This is helpful if you are running behind a proxy or are hosting multiple websites from a single URL.

### Application Settings

In the main.go file, you can modify your log level:

```go
ambientApp, err := ambient.NewApp(...)
ambientApp.SetLogLevel(ambient.LogLevelDebug)
ambientApp.SetLogLevel(ambient.LogLevelInfo)
ambientApp.SetLogLevel(ambient.LogLevelError)
ambientApp.SetLogLevel(ambient.LogLevelFatal)
```

In the main.go file, you can enable `span` tags around HTML elements to determine which content is loaded from which plugins:

```go
ambientApp, err := ambient.NewApp(...)
ambientApp.SetDebugTemplates(true)
```

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