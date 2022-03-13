package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/internal/config"
	"github.com/ambientkit/ambient/internal/injector"
	"github.com/ambientkit/ambient/internal/pluginsafe"
	"github.com/ambientkit/ambient/internal/secureconfig"
	"github.com/ambientkit/ambient/pkg/grpcp"
	"github.com/ambientkit/plugin/logger/zaplogger"
	"github.com/ambientkit/plugin/router/awayrouter"
	"github.com/ambientkit/plugin/sessionmanager/scssession"
	"github.com/ambientkit/plugin/storage/memorystorage"
	"github.com/ambientkit/plugin/templateengine/htmlengine"
	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
)

func main() {
	z := zaplogger.New()
	logger, err := z.Logger("grpcplugin", "1.0.0", nil)
	if err != nil {
		log.Fatalln(err.Error())
	}

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

	r := awayrouter.New(h)
	router, err := r.Router(logger, nil)
	if err != nil {
		logger.Fatal(err.Error())
	}

	ms := memorystorage.New()
	ds, ss, err := ms.Storage(logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	sessPlugin := scssession.New("5ba3ad678ee1fd9c4fddcef0d45454904422479ed762b3b0ddc990e743cb65e0")
	sess, err := sessPlugin.SessionManager(logger, ss)
	if err != nil {
		logger.Fatal(err.Error())
	}

	tePlugin := htmlengine.New()

	pl := &ambient.PluginLoader{
		Router:         r,
		TemplateEngine: tePlugin,
		SessionManager: sessPlugin,
		Plugins:        []ambient.Plugin{},
		Middleware: []ambient.MiddlewarePlugin{
			sessPlugin,
		},
	}

	storage, err := config.NewStorage(logger, ds, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Initialize the plugin system.
	pluginsystem, err := config.NewPluginSystem(logger, storage, pl)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Set up the template injector.
	pi := injector.NewPlugininjector(logger, pluginsystem, sess, false, true)

	te, err := tePlugin.TemplateEngine(logger, pi)
	if err != nil {
		log.Fatal(err.Error())
	}

	recorder := pluginsafe.NewRouteRecorder(logger, pluginsystem, router)

	// Create secure site for the core app and use "ambient" so it gets
	// full permissions.
	securesite := secureconfig.NewSecureSite("ambient", logger, pluginsystem, sess, router, te, recorder)

	mw := sessPlugin.Middleware()[0]

	err = connectPlugin(logger, router, securesite, mw, "hello", "./cmd/plugin/hello/cmd/plugin/hello")
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func connectPlugin(logger ambient.AppLogger, router ambient.AppRouter, site grpcp.SecureSite, mw func(next http.Handler) http.Handler,
	pluginName string, pluginPath string) error {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: grpcp.Handshake,
		Plugins: map[string]plugin.Plugin{
			pluginName: &grpcp.GenericPlugin{},
		},
		Cmd: exec.Command(pluginPath),
		Logger: hclog.New(&hclog.LoggerOptions{
			Level:      hclog.Debug,
			Output:     os.Stderr,
			JSONFormat: true,
		}),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC,
		},
	})
	defer client.Kill()

	// Connect via RPC.
	rpcClient, err := client.Client()
	if err != nil {
		return fmt.Errorf("server: could not get gRPC client: %v", err.Error())
	}

	// Request the plugin.
	raw, err := rpcClient.Dispense(pluginName)
	if err != nil {
		return fmt.Errorf("server: could not get connect to plugin (%v): %v", pluginName, err.Error())
	}

	p := raw.(grpcp.PluginCore)
	// if !ok {
	// 	fmt.Println("The plugin is not the right format.")
	// 	return
	// }

	toolkit := &grpcp.Toolkit{
		Log:  logger,
		Mux:  router,
		Site: site,
	}

	name, err := p.PluginName()
	if err != nil {
		return fmt.Errorf("server: could not get plugin name: %v", err.Error())
	}
	logger.Info("Plugin name: %v", name)

	version, err := p.PluginVersion()
	if err != nil {
		return fmt.Errorf("server: could not get plugin version: %v", err.Error())
	}
	logger.Info("Plugin version: %v", version)

	err = p.Enable(toolkit)
	if err != nil {
		return fmt.Errorf("server: could not enable: %v", err.Error())
	}

	err = p.Routes()
	if err != nil {
		return fmt.Errorf("server: could not get routes: %v", err.Error())
	}

	return http.ListenAndServe(":8080", mw(router))

	// for {
	// 	<-time.After(5 * time.Second)
	// 	err = p.Routes()
	// 	if err != nil {
	// 		return fmt.Errorf("server: could not get routes: %v", err.Error())
	// 	}

	// 	err = p.Disable()
	// 	if err != nil {
	// 		return fmt.Errorf("server: could not disable plugin: %v", err.Error())
	// 	}
	// }
}
