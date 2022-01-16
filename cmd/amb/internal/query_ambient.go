package internal

import (
	"github.com/c-bata/go-prompt"
)

// pluginNames returns a list of plugin names as suggestions.
func pluginNames() []prompt.Suggest {
	arr := make([]prompt.Suggest, 0)
	arr = append(arr, prompt.Suggest{Text: "all", Description: ""})

	// Get the plugin names.
	pluginNames := make([]string, 0)
	err := rc.Get("/plugins", &pluginNames)
	if err != nil {
		log.Error("amb: could not get plugin names: %v", err.Error())
		return nil
	}

	// Build a list of suggestions.
	for _, pluginName := range pluginNames {
		arr = append(arr, prompt.Suggest{Text: pluginName, Description: ""})
	}

	return arr
}
