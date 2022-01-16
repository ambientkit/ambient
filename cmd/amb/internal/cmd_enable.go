package internal

import (
	"fmt"

	"github.com/c-bata/go-prompt"
)

// CmdEnable represents a command object.
type CmdEnable struct {
	CmdBase
}

// Command returns the initial command.
func (c *CmdEnable) Command() string {
	return "enable"
}

// Suggestion returns the suggestion for the initial command.
func (c *CmdEnable) Suggestion() prompt.Suggest {
	return prompt.Suggest{Text: c.Command(), Description: "Enable plugin..."}
}

// Executer executes the command.
func (c *CmdEnable) Executer(args []string) {
	if len(args) < 2 {
		log.Info("amb: command not recognized")
		return
	}

	if args[1] == "all" {
		// Enable all plugins.
		log.Info("amb: enabling all trusted plugins")

		err := rc.Post("/plugins/enable", nil, nil)
		if err != nil {
			log.Error("amb: could not enable all plugins: %v", err.Error())
		}
	} else {
		// Enable one plugin.
		pluginName := args[1]
		log.Info("amb: enabling plugin: %v", pluginName)

		err := rc.Post(fmt.Sprintf("/plugins/%v/enable", pluginName), nil, nil)
		if err != nil {
			log.Error("amb: could not enable plugin, %v: %v", pluginName, err.Error())
		}
	}
}

// Completer returns a list of suggestions based on the user input.
func (c *CmdEnable) Completer(d prompt.Document, args []string) []prompt.Suggest {
	// Show a secondary list of plugin suggestions.
	if len(args) == 2 {
		return prompt.FilterHasPrefix(pluginNames(), args[1], true)
	}

	// Else return nothing.
	return prompt.FilterHasPrefix([]prompt.Suggest{}, d.TextBeforeCursor(), true)
}
