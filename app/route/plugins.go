package route

import (
	"fmt"

	"github.com/josephspurrier/ambient/app/lib/ambsystem"
	"github.com/josephspurrier/ambient/app/lib/router"
	"github.com/josephspurrier/ambient/plugin/stackedit"
)

// LoadPlugins will load the plugins.
func LoadPlugins(mux *router.Mux) {
	// Define the plugins.
	plugins := []ambsystem.IPlugin{
		stackedit.Activate(),
	}

	// Load the plugins.
	for _, v := range plugins {
		fmt.Printf("Load plugin: %v\n", v.PluginName())
		v.SetPages(mux)
	}
}

// func setLogging() {
// 	//temp := os.Stdout
// 	// Turn off logging so it plugins don't have control over output.
// 	os.Stdout = nil
// 	//os.Stdout = temp   // restore it
// }
