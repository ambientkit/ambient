package plugins

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/josephspurrier/ambient/app/lib/ambsystem"
	"github.com/josephspurrier/ambient/app/lib/datastorage"
	"github.com/josephspurrier/ambient/app/lib/router"
	"github.com/josephspurrier/ambient/app/model"
	"github.com/josephspurrier/ambient/app/modelsecure"
	"github.com/josephspurrier/ambient/app/route"
)

// Load the plugins into storage.
func Load(arr []ambsystem.IPlugin, storage *datastorage.Storage) (*ambsystem.PluginSystem, error) {
	// Create the plugin system.
	pluginsys := ambsystem.NewPluginSystem()

	// Load the plugins.
	needSave := false
	ps := storage.Site.Plugins
	for _, v := range arr {
		name := v.PluginName()
		_, found := ps[name]
		if !found {
			fmt.Printf("Load new plugin: %v\n", name)
			ps[name] = model.PluginSettings{
				Enabled: false,
			}
			needSave = true
		} else {
			fmt.Printf("Plugin already found: %v\n", name)
		}

		// Add to the system.
		pluginsys.Plugins[name] = v
	}

	if needSave {
		// Save the plugins.
		storage.Site.Plugins = ps
		err := storage.Save()
		if err != nil {
			return nil, err
		}
	}

	return pluginsys, nil
}

// Pages loads the pages from the plugins.
func Pages(c *route.Core) error {
	// Set up the plugin routes.
	shouldSave := false
	ps := c.Storage.Site.Plugins
	for name, plugin := range ps {
		// Determine if the plugin that is in stored is found in the system.
		v, found := c.Plugins.Plugins[name]

		// If the found setting is different, then update it for saving.
		if found != plugin.Found {
			shouldSave = true
			plugin.Found = found
			ps[name] = plugin
		}

		// If the plugin is not found or not enabled, then skip over it.
		if !found || !plugin.Enabled {
			continue
		}

		grants := make(map[string]bool)
		grants["site.title:read"] = true
		grants["site.plugins:read"] = true
		grants["site.plugins:enable"] = true
		grants["site.plugins:disable"] = true
		grants["site.plugins:deleteone"] = true
		grants["router:clear"] = true

		recorder := router.NewRecorder(c.Router)

		toolkit := &ambsystem.Toolkit{
			Router:   recorder,
			Render:   c.Render,
			Security: c.Sess,
			Site:     modelsecure.NewSecureSite(name, c.Storage, c.Router, grants),
		}

		// Load the pages.
		err := v.SetPages(toolkit)
		if err != nil {
			log.Printf("problem loading pages from plugin %v: %v", name, err.Error())
		}

		fmt.Println("Routes:", recorder.Routes())

		arr := make([]model.Route, 0)
		for _, route := range recorder.Routes() {
			arr = append(arr, model.Route{
				Method: route.Method,
				Path:   route.Path,
			})
		}
		c.Storage.PluginRoutes.Routes[name] = arr

		// Load the assets.
		assets, files := v.Assets()
		if files == nil {
			continue
		}

		fmt.Println("loading assets for:", name)

		// Handle embedded assets.
		err = EmbeddedAssets(c.Router, name, assets, files)
		if err != nil {
			log.Println(err.Error())
		}

	}

	if shouldSave {
		// Save the plugin state if something changed.
		c.Storage.Site.Plugins = ps
		err := c.Storage.Save()
		if err != nil {
			return err
		}
	}

	return nil
}

func EmbeddedAssets(mux *router.Mux, pluginName string, files []ambsystem.Asset, assets *embed.FS) error {
	for _, v := range files {
		// Skip files that are not embedded.
		if !v.Embedded {
			continue
		}

		fileurl := path.Join("/plugins", pluginName, v.SanitizedPath())

		// TODO: Need to check for missing locations and types.

		exists := fileExists(assets, v.SanitizedPath())
		if !exists {
			return fmt.Errorf("plugin (%v) has missing file, please check 'SetAssets()': %v", pluginName, v)
		}

		mux.Get(fileurl, func(w http.ResponseWriter, r *http.Request) (statusCode int, err error) {
			// Don't allow directory browsing.
			if strings.HasSuffix(r.URL.Path, "/") {
				return http.StatusNotFound, nil
			}

			// Use the root directory.
			fsys, err := fs.Sub(assets, ".")
			if err != nil {
				return http.StatusInternalServerError, err
			}

			// Get the requested file name.
			fname := strings.TrimPrefix(r.URL.Path, path.Join("/plugins", pluginName)+"/")

			// Open the file.
			f, err := fsys.Open(fname)
			if err != nil {
				return http.StatusNotFound, nil
			}
			defer f.Close()

			// Get the file time.
			st, err := f.Stat()
			if err != nil {
				return http.StatusInternalServerError, err
			}

			http.ServeContent(w, r, fname, st.ModTime(), f.(io.ReadSeeker))
			return
		})
	}

	return nil
}

// fileExists determines if an embedded file exists.
func fileExists(assets *embed.FS, filename string) bool {
	// Use the root directory.
	fsys, err := fs.Sub(assets, ".")
	if err != nil {
		return false
	}

	// Open the file.
	f, err := fsys.Open(filename)
	if err != nil {
		return false
	}
	defer f.Close()

	return true
}
