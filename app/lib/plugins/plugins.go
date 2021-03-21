package plugins

import (
	"fmt"
	"log"

	"github.com/josephspurrier/ambient/app/lib/ambsystem"
	"github.com/josephspurrier/ambient/app/lib/datastorage"
	"github.com/josephspurrier/ambient/app/lib/router"
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
			ps[name] = ambsystem.PluginSettings{
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
func Pages(storage *datastorage.Storage, mux *router.Mux, plugins *ambsystem.PluginSystem) error {
	// Set up the plugin routes.
	shouldSave := false
	ps := storage.Site.Plugins
	for name, plugin := range ps {
		if !plugin.Enabled {
			continue
		}

		// Determine if the plugin that is in stored is found in the system.
		v, found := plugins.Plugins[name]

		// If the found setting is different, then update it for saving.
		if found != plugin.Found {
			shouldSave = true
			plugin.Found = found
			ps[name] = plugin
		}

		// If the plugin is not found, then skip over trying to read from it.
		if !found {
			continue
		}

		// Load the pages.
		err := v.SetPages(mux)
		if err != nil {
			log.Printf("problem loading pages from plugin %v: %v", name, err.Error())
		}
	}

	if shouldSave {
		// Save the plugin state if something changed.
		storage.Site.Plugins = ps
		err := storage.Save()
		if err != nil {
			return err
		}
	}

	return nil
}
