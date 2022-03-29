package devconsole

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/internal/secureconfig"
	"github.com/ambientkit/ambient/pkg/envdetect"
	"github.com/ambientkit/away/router"
)

// DevConsole represents a web interface to receive commands from the amb tool.
type DevConsole struct {
	log           ambient.AppLogger
	storage       ambient.Storage
	pluginsystem  ambient.PluginSystem
	securestorage *secureconfig.SecureSite
}

// NewDevConsole returns the dev console object to receive commands from the amb
// tool.
func NewDevConsole(logger ambient.AppLogger, ps ambient.PluginSystem, storage ambient.Storage, site *secureconfig.SecureSite) *DevConsole {
	return &DevConsole{
		log:           logger,
		storage:       storage,
		pluginsystem:  ps,
		securestorage: site,
	}
}

// JSON sends data to the writer.
func JSON(w http.ResponseWriter, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return ambient.StatusError{Code: http.StatusInternalServerError, Err: err}
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(b))

	return nil
}

// EnableDevConsole turns on the dev console web listener.
func (dc *DevConsole) EnableDevConsole() {
	dc.log.Info("started and available at: %v/%v", envdetect.DevConsoleURL(), envdetect.DevConsolePort())

	go func() {
		mux := router.New()

		// Encrypt the site JSON file on disk.
		mux.Post("/storage/encrypt", func(w http.ResponseWriter, r *http.Request) error {
			dc.log.Debug("site.bin encrypted")
			err := dc.storage.LoadDecrypted()
			if err != nil {
				return ambient.StatusError{Code: http.StatusInternalServerError, Err: err}
			}
			err = dc.storage.Save()
			if err != nil {
				return ambient.StatusError{Code: http.StatusInternalServerError, Err: err}
			}

			return nil
		})

		// Decrypt the site JSON file on disk.
		mux.Post("/storage/decrypt", func(w http.ResponseWriter, r *http.Request) error {
			dc.log.Debug("site.bin decrypted")
			err := dc.storage.SaveDecrypted()
			if err != nil {
				return ambient.StatusError{Code: http.StatusInternalServerError, Err: err}
			}

			return nil
		})

		// Return a list of plugin names.
		mux.Get("/plugins", func(w http.ResponseWriter, r *http.Request) error {
			dc.log.Debug("get plugin names")
			return JSON(w, dc.pluginsystem.TrustedPluginNames())
		})

		// Enable one plugin.
		mux.Post("/plugins/{pluginName}/enable", func(w http.ResponseWriter, r *http.Request) error {
			pluginName := mux.Param(r, "pluginName")
			dc.log.Debug("enable plugin: %v", pluginName)

			err := dc.securestorage.EnablePlugin(pluginName, true)
			if err != nil {
				return ambient.StatusError{Code: http.StatusBadRequest, Err: err}
			}

			return nil
		})

		// Enable all plugins.
		mux.Post("/plugins/enable", func(w http.ResponseWriter, r *http.Request) error {
			dc.log.Debug("enable all plugins")

			// Loop through all the trusted plugins.
			for _, pluginName := range dc.pluginsystem.TrustedPluginNames() {
				err := dc.securestorage.EnablePlugin(pluginName, true)
				if err != nil {
					// TODO: Should return an error at the end if at least one fails.
					dc.log.Error("failed to enable plugin (%v): %v", pluginName, err.Error())
					// Continue on
				}
			}

			return nil
		})

		// Enable all grants for one plugin.
		mux.Post("/plugins/{pluginName}/grant", func(w http.ResponseWriter, r *http.Request) error {
			pluginName := mux.Param(r, "pluginName")
			dc.log.Debug("enable plugin grants: %v", pluginName)

			p, err := dc.pluginsystem.Plugin(pluginName)
			if err != nil {
				return ambient.StatusError{Code: http.StatusBadRequest,
					Err: fmt.Errorf("failed to get plugin (%v) for grants: %v", pluginName, err.Error())}
			}

			for _, request := range p.GrantRequests() {
				dc.log.Debug("plugin (%v), add grant: %v", pluginName, request.Grant)
				err := dc.securestorage.SetNeighborPluginGrant(pluginName, request.Grant, true)
				if err != nil {
					return ambient.StatusError{Code: http.StatusBadRequest,
						Err: fmt.Errorf("failed to enable plugin (%v) for grant, %v: %v", pluginName, request.Grant, err.Error())}
				}
			}

			return nil
		})

		// Enable all grants for all plugins.
		mux.Post("/plugins/grant", func(w http.ResponseWriter, r *http.Request) error {
			pluginName := mux.Param(r, "pluginName")
			dc.log.Debug("enable plugin grant: %v", pluginName)

			// Loop through all the trusted plugins.
			for _, pluginName := range dc.pluginsystem.TrustedPluginNames() {
				p, err := dc.pluginsystem.Plugin(pluginName)
				if err != nil {
					return ambient.StatusError{Code: http.StatusBadRequest,
						Err: fmt.Errorf("failed to get plugin (%v) for grants: %v", pluginName, err.Error())}
				}

				for _, request := range p.GrantRequests() {
					dc.log.Debug("plugin (%v), add grant: %v", pluginName, request.Grant)
					err := dc.securestorage.SetNeighborPluginGrant(pluginName, request.Grant, true)
					if err != nil {
						return ambient.StatusError{Code: http.StatusBadRequest,
							Err: fmt.Errorf("failed to enable plugin (%v) for grant, %v: %v", pluginName, request.Grant, err.Error())}
					}
				}
			}

			return nil
		})

		err := http.ListenAndServe(":"+envdetect.DevConsolePort(), mux)
		if err != nil {
			dc.log.Error("listener cannot start: %v", err.Error())
		}
	}()
}
