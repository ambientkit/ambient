package testutil

import (
	"log"
	"net/http"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/ambientapp"
	"github.com/ambientkit/ambient/pkg/grpcp"
	"github.com/ambientkit/ambient/pkg/grpcp/testdata/plugin/hello"
	"github.com/ambientkit/ambient/pkg/grpcp/testdata/plugin/neighbor"
	"github.com/ambientkit/plugin/logger/zaplogger"
	"github.com/ambientkit/plugin/router/awayrouter"
	"github.com/ambientkit/plugin/sessionmanager/scssession"
	"github.com/ambientkit/plugin/storage/memorystorage"
	"github.com/ambientkit/plugin/templateengine/htmlengine"
	"github.com/hashicorp/go-plugin"
)

// Setup sets up a test gRPC server.
func Setup() (ambient.Plugin, *plugin.Client, http.Handler, error) {
	// z := zaplogger.New()
	// logger, err := z.Logger("grpcplugin", "1.0.0", nil)
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }
	// logger.SetLogLevel(ambient.LogLevelDebug)

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
			hello.New(),
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
	app, logger, err := ambientapp.NewApp("myapp", "1.0",
		zaplogger.New(),
		ambient.StoragePluginGroup{
			Storage: memorystorage.New(),
		},
		plugins)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// r := awayrouter.New(h)
	// router, err := r.Router(logger, nil)
	// if err != nil {
	// 	logger.Fatal(err.Error())
	// }

	// ms := memorystorage.New()
	// ds, ss, err := ms.Storage(logger)
	// if err != nil {
	// 	logger.Fatal(err.Error())
	// }

	// sessPlugin := scssession.New("5ba3ad678ee1fd9c4fddcef0d45454904422479ed762b3b0ddc990e743cb65e0")
	// sess, err := sessPlugin.SessionManager(logger, ss)
	// if err != nil {
	// 	logger.Fatal(err.Error())
	// }

	// tePlugin := htmlengine.New()

	// pl := &ambient.PluginLoader{
	// 	Router:         r,
	// 	TemplateEngine: tePlugin,
	// 	SessionManager: sessPlugin,
	// 	Plugins: []ambient.Plugin{
	// 		hello.New(),
	// 		neighbor.New(),
	// 	},
	// 	Middleware: []ambient.MiddlewarePlugin{
	// 		sessPlugin,
	// 	},
	// }

	// storage, err := config.NewStorage(logger, ds, nil)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// // Initialize the plugin system.
	// pluginsystem, err := config.NewPluginSystem(logger, storage, pl)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// // Set up the template injector.
	// pi := injector.NewPlugininjector(logger, pluginsystem, sess, false, true)

	// te, err := tePlugin.TemplateEngine(logger, pi)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// recorder := pluginsafe.NewRouteRecorder(logger, pluginsystem, router)

	// // Create secure site for the core app.
	// securesite := secureconfig.NewSecureSite("hello", logger, pluginsystem, sess, router, te, recorder)
	// rr := recorder.WithPlugin("hello")

	// toolkit := &ambient.Toolkit{
	// 	Log:  logger,
	// 	Mux:  rr,
	// 	Site: securesite,
	// }
	// secureconfig.SaveRoutesForPlugin("hello", rr, pluginsystem)

	// mw := sessPlugin.Middleware()[0]

	p, pluginClient, err := grpcp.ConnectPlugin("hello", "./pkg/grpcp/testdata/plugin/hello/cmd/plugin/hello")
	if err != nil {
		logger.Fatal(err.Error())
	}

	// name := p.PluginName()
	// // if err != nil {
	// // 	return nil, pluginClient, nil, fmt.Errorf("server: could not get plugin name: %v", err.Error())
	// // }
	// logger.Info("Plugin name: %v", name)

	// version := p.PluginVersion()
	// // if err != nil {
	// // 	return nil, pluginClient, nil, fmt.Errorf("server: could not get plugin version: %v", err.Error())
	// // }
	// logger.Info("Plugin version: %v", version)

	// err = p.Enable(toolkit)
	// if err != nil {
	// 	return nil, pluginClient, nil, fmt.Errorf("server: could not enable: %v", err.Error())
	// }

	// p.Routes()
	// if err != nil {
	// 	return nil, pluginClient, nil, fmt.Errorf("server: could not get routes: %v", err.Error())
	// }

	// 	err = p.Disable()
	// 	if err != nil {
	// 		return fmt.Errorf("server: could not disable plugin: %v", err.Error())
	// 	}
	// }

	handler, err := app.Handler()
	if err != nil {
		logger.Fatal(err.Error())
	}

	return p, pluginClient, handler, err
}
