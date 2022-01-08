package ambient

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/josephspurrier/ambient/plugin/router/awayrouter/router"
)

// DevConsole represents a web interface to receive commands from the amb tool.
type DevConsole struct {
	port         string
	log          AppLogger
	storage      *Storage
	pluginsystem *PluginSystem
}

// NewDevConsole returns the dev console object to receive commands from the amb
// tool.
func NewDevConsole(port string, log AppLogger, storage *Storage, pluginsystem *PluginSystem) *DevConsole {
	return &DevConsole{
		port:         port,
		log:          log,
		storage:      storage,
		pluginsystem: pluginsystem,
	}
}

// JSON sends data to the writer.
func JSON(w http.ResponseWriter, data interface{}) (int, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(b))

	return http.StatusOK, nil
}

// EnableDevConsole turns on the dev console web listener.
func (dc *DevConsole) EnableDevConsole() {
	dc.log.Info("ambient: dev console started on: %v", dc.port)

	// Create secure site for the core app and use "ambient" so it gets
	// full permissions.
	securestorage := NewSecureSite("ambient", dc.log, dc.pluginsystem, nil, nil, nil)

	go func() {
		mux := router.New()

		// Encrypt the site JSON file on disk.
		mux.Post("/storage/encrypt", func(w http.ResponseWriter, r *http.Request) (int, error) {
			dc.log.Debug("ambient: dev console - site.bin encrypted")
			err := dc.storage.LoadDecrypted()
			if err != nil {
				return http.StatusInternalServerError, err
			}
			err = dc.storage.Save()
			if err != nil {
				return http.StatusInternalServerError, err
			}

			return http.StatusOK, nil
		})

		// Decrypt the site JSON file on disk.
		mux.Post("/storage/decrypt", func(w http.ResponseWriter, r *http.Request) (int, error) {
			dc.log.Debug("ambient: dev console - site.bin decrypted")
			err := dc.storage.SaveDecrypted()
			if err != nil {
				return http.StatusInternalServerError, err
			}

			return http.StatusOK, nil
		})

		// Return a list of plugin names.
		mux.Get("/plugins", func(w http.ResponseWriter, r *http.Request) (int, error) {
			dc.log.Debug("ambient: dev console - get plugin names")
			return JSON(w, dc.pluginsystem.TrustedPluginNames())
		})

		// Enable one plugin.
		mux.Post("/plugins/:pluginName/enable", func(w http.ResponseWriter, r *http.Request) (int, error) {
			pluginName := mux.Param(r, "pluginName")
			dc.log.Info("ambient: dev console - enable plugin: %v", pluginName)

			err := securestorage.EnablePlugin(pluginName, false)
			if err != nil {
				return http.StatusBadRequest, err
			}

			return http.StatusOK, nil
		})

		// Enable all plugins.
		mux.Post("/plugins/enable", func(w http.ResponseWriter, r *http.Request) (int, error) {
			dc.log.Info("ambient: dev console - enable all plugins")

			// Loop through all the trusted plugins.
			for _, pluginName := range dc.pluginsystem.TrustedPluginNames() {
				err := securestorage.EnablePlugin(pluginName, false)
				if err != nil {
					// TODO: Should return an error at the end if at least one fails.
					dc.log.Error("ambient: dev console - failed to enable plugin (%v): %v", pluginName, err.Error())
					// Continue on
				}
			}

			return http.StatusOK, nil
		})

		// Enable all grants for one plugin.
		mux.Post("/plugins/:pluginName/grant", func(w http.ResponseWriter, r *http.Request) (int, error) {
			pluginName := mux.Param(r, "pluginName")
			dc.log.Info("ambient: dev console - enable plugin grants: %v", pluginName)

			p, err := dc.pluginsystem.Plugin(pluginName)
			if err != nil {
				return http.StatusBadRequest, fmt.Errorf("failed to get plugin (%v) for grants: %v", pluginName, err.Error())
			}

			for _, request := range p.GrantRequests() {
				dc.log.Info("ambient: dev console - plugin, %v - add grant: %v", pluginName, request.Grant)
				err := securestorage.SetNeighborPluginGrant(pluginName, request.Grant, true)
				if err != nil {
					return http.StatusBadRequest, fmt.Errorf("failed to enable plugin (%v) for grant, %v: %v", pluginName, request.Grant, err.Error())
				}
			}

			return http.StatusOK, nil
		})

		// Enable all grants for all plugins.
		mux.Post("/plugins/grant", func(w http.ResponseWriter, r *http.Request) (int, error) {
			pluginName := mux.Param(r, "pluginName")
			dc.log.Info("ambient: dev console - enable plugin grant: %v", pluginName)

			// Loop through all the trusted plugins.
			for _, pluginName := range dc.pluginsystem.TrustedPluginNames() {
				p, err := dc.pluginsystem.Plugin(pluginName)
				if err != nil {
					return http.StatusBadRequest, fmt.Errorf("failed to get plugin (%v) for grants: %v", pluginName, err.Error())
				}

				for _, request := range p.GrantRequests() {
					dc.log.Info("ambient: dev console - plugin, %v - add grant: %v", pluginName, request.Grant)
					err := securestorage.SetNeighborPluginGrant(pluginName, request.Grant, true)
					if err != nil {
						return http.StatusBadRequest, fmt.Errorf("failed to enable plugin (%v) for grant, %v: %v", pluginName, request.Grant, err.Error())
					}
				}
			}

			return http.StatusOK, nil
		})

		err := http.ListenAndServe(":"+dc.port, mux)
		if err != nil {
			dc.log.Error("ambient: dev config server cannot start: %v", err.Error())
		}
	}()
}
