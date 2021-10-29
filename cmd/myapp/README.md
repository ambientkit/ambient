# Sample App Tutorial

This folder contains a sample app to demonstrate how to use Ambient.

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

- data storage will use the local filesystem instead of Google Cloud Storage
- if you try to access the app, it will listen on all IPs/addresses, instead of redirecting like it does in production

You can use `envdetect.RunningLocalDev()` to detect if the flag is set to true or not.

When `AMB_TIMEZONE` is set to a timezone like `America/New_York`, the app will use that timezone. This is required if using time-based packages like MFA.

When `AMB_URL_PREFIX` is set to a path like `/api`, the app will serve requests from `/api/...`. This is helpful if you are running behind a proxy or are hosting multiple websites from a single URL.

### App Settings

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

If you would like to make changes to the code, I recommend `air` to help streamline your workflow.

```bash
# Install air to allow hot reloading so you can make changes quickly.
curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
```

You can then use this command to start the web server and monitor for changes:

```bash
# Start hot reload. The web app should be available at: http://localhost:8080
air
```