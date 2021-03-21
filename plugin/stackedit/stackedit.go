package stackedit

import (
	"fmt"
	"net/http"

	"github.com/josephspurrier/ambient/app/lib/ambsystem"
)

// StackEdit -
type StackEdit struct {
	ambsystem.PluginMeta
}

// Activate installs and enables the plugin.
func Activate() StackEdit {
	return StackEdit{
		PluginMeta: ambsystem.PluginMeta{
			Name:       "stackedit",
			Version:    "1.0.0",
			AppVersion: "1.0.0",
		},
	}
}

// SetPages -
func (pm StackEdit) SetPages(mux ambsystem.IRouter) error {
	fmt.Println("StackEdit page registered.")

	mux.Get("/se", func(w http.ResponseWriter, r *http.Request) (int, error) {
		fmt.Fprint(w, "Plugin, stackedit, is loaded.")
		return http.StatusOK, nil
	})

	return nil
}

// // SetSettings -
// func (pm StackEdit) SetSettings(s ambsystem.ISettings) error {
// 	s.Add("name string", fieldType string, defaultValue string)

// }

// Deactivate deactivates the plugin, but leaves the state in the system.
func Deactivate() error {
	return nil
}

// Uninstall removes all plugin state from the system.
func Uninstall() error {
	return nil
}
