package testutil

import (
	"net/http"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/ambientapp"
	"github.com/ambientkit/ambient/pkg/grpcp/testdata/plugin/neighbor"
	"github.com/ambientkit/plugin/logger/zaplogger"
	"github.com/ambientkit/plugin/router/awayrouter"
	"github.com/ambientkit/plugin/sessionmanager/scssession"
	"github.com/ambientkit/plugin/storage/memorystorage"
	"github.com/ambientkit/plugin/templateengine/htmlengine"
)

// Setup sets up a test gRPC server.
func Setup() (*ambientapp.App, error) {
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

	sessPlugin := scssession.New("5ba3ad678ee1fd9c4fddcef0d45454904422479ed762b3b0ddc990e743cb65e0")
	plugins := &ambient.PluginLoader{
		// Core plugins are implicitly trusted.
		Router:         awayrouter.New(h),
		TemplateEngine: htmlengine.New(),
		SessionManager: sessPlugin,
		// Trusted plugins are those that are typically needed to boot so they
		// will be enabled and given full access.
		TrustedPlugins: map[string]bool{},
		Plugins: []ambient.Plugin{
			//hello.New(),
			neighbor.New(),
		},
		GRPCPlugins: []ambient.GRPCPlugin{
			{PluginName: "hello", PluginPath: "./pkg/grpcp/testdata/plugin/hello/cmd/plugin/hello"},
		},
		Middleware: []ambient.MiddlewarePlugin{
			// Middleware - executes bottom to top.
			sessPlugin,
		},
	}
	app, _, err := ambientapp.NewApp("myapp", "1.0",
		zaplogger.New(),
		ambient.StoragePluginGroup{
			Storage: memorystorage.New(),
		},
		plugins)
	return app, err
}
