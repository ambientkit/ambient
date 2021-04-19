// Package extend is an extension of the plugin system.
package extend

import (
	"fmt"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/app/extend/plugin/hello"
)

// PluginLoaderExtended -
type PluginLoaderExtended struct {
	ambient.PluginLoader

	CustomPlugins []PluginExtended
}

// PluginExtended represents an extended plugin.
type PluginExtended interface {
	ambient.PluginCore

	CustomFunction() string
}

func newApp() {
	pl := PluginLoaderExtended{
		//PluginLoader: ,
		CustomPlugins: []PluginExtended{
			hello.New(),
		},
	}

	fmt.Println(pl)
}
