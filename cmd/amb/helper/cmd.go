package helper

import (
	"github.com/c-bata/go-prompt"
)

var (
	// Commands available.
	execCreateApp = "createapp"
	execEnable    = "enable"
	execGrants    = "grant"
	execEncrypt   = "encryptstorage"
	execDecrypt   = "decryptstorage"
	execExit      = "exit"

	// Prompts should match 1:1 with the commands above.
	promptSuggestions = []prompt.Suggest{
		{Text: execCreateApp, Description: "Create new Ambient app"},
		{Text: execEnable, Description: "Enable plugin..."},
		{Text: execGrants, Description: "Add grants for plugin..."},
		{Text: execEncrypt, Description: "Encrypt storage"},
		{Text: execDecrypt, Description: "Decrypt storage"},
		{Text: execExit, Description: "Exit the CLI (or press Ctrl+C)"},
	}
)

// pluginSuggestions returns a list of suggestions for plugins.
func pluginSuggestions() []prompt.Suggest {
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
	//arr = append(arr, ConvertToSuggestions(pluginNames)...)

	return arr
}
