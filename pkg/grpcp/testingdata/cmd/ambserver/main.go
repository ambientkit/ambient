package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/ambientapp"
	"github.com/ambientkit/plugin/logger/zaplogger"
	"github.com/ambientkit/plugin/router/awayrouter"
	"github.com/ambientkit/plugin/sessionmanager/scssession"
	"github.com/ambientkit/plugin/storage/memorystorage"
	"github.com/ambientkit/plugin/templateengine/htmlengine"
)

func main() {
	// Create the ambient app.
	plugins := Plugins()
	ambientApp, logger, err := ambientapp.NewApp("myapp", "1.0",
		zaplogger.New(),
		ambient.StoragePluginGroup{
			Storage: memorystorage.New(),
		},
		plugins)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Load the plugins and return the handler.
	mux, err := ambientApp.Handler()
	if err != nil {
		logger.Fatal(err.Error())
	}

	// Start the web listener.
	os.Setenv("PORT", "8081")
	ambientApp.ListenAndServe(mux)
}

// Plugins defines the plugins.
func Plugins() *ambient.PluginLoader {
	// Define the session manager so it can be used as a core plugin and
	// middleware.
	sessionManager := scssession.New("5ba3ad678ee1fd9c4fddcef0d45454904422479ed762b3b0ddc990e743cb65e0")

	h := func(log ambient.Logger, renderer ambient.Renderer, w http.ResponseWriter, r *http.Request, err error) {
		if err != nil {
			switch e := err.(type) {
			case ambient.Error:
				errText := e.Error()
				if len(errText) == 0 {
					errText = http.StatusText(e.Status())
				}
				http.Error(w, errText, e.Status())
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}
	}

	return &ambient.PluginLoader{
		// Core plugins are implicitly trusted.
		Router:         awayrouter.New(h),
		TemplateEngine: htmlengine.New(),
		SessionManager: sessionManager,
		// Trusted plugins are those that are typically needed to boot so they
		// will be enabled and given full access.
		TrustedPlugins: map[string]bool{
			"hello": true,
		},
		Plugins: []ambient.Plugin{
			New(),
		},
		Middleware: []ambient.MiddlewarePlugin{
			// Middleware - executes bottom to top.
			sessionManager, // Session manager middleware.
		},
	}
}
