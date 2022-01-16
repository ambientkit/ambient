package internal

import (
	"fmt"

	"github.com/c-bata/go-prompt"
)

// CmdGrant represents a command object.
type CmdGrant struct {
	CmdBase
}

// Command returns the initial command.
func (c *CmdGrant) Command() string {
	return "grant"
}

// Suggestion returns the suggestion for the initial command.
func (c *CmdGrant) Suggestion() prompt.Suggest {
	return prompt.Suggest{Text: c.Command(), Description: "Add grants for plugin..."}
}

// Executer executes the command.
func (c *CmdGrant) Executer(args []string) {
	if len(args) < 2 {
		log.Info("amb: command not recognized")
		return
	}

	if args[1] == "all" {
		// Enable grants for all plugins.
		log.Info("amb: adding grants for all trusted plugins")

		err := rc.Post("/plugins/grant", nil, nil)
		if err != nil {
			log.Error("amb: cloud not enable all plugins grants: %v", err.Error())
		}
	} else {
		// Enable grants for one plugin.
		pluginName := args[1]
		log.Info("amb: adding grants for plugin: %v", pluginName)

		err := rc.Post(fmt.Sprintf("/plugins/%v/grant", pluginName), nil, nil)
		if err != nil {
			log.Error("amb: cloud not enable plugin (%v) grants: %v", pluginName, err.Error())
		}
	}
}

// Completer returns a list of suggestions based on the user input.
func (c *CmdGrant) Completer(d prompt.Document, args []string) []prompt.Suggest {
	// Show a secondary list of plugin suggestions.
	if len(args) == 2 {
		return prompt.FilterHasPrefix(pluginNames(), args[1], true)
	}

	// Else return nothing.
	return prompt.FilterHasPrefix([]prompt.Suggest{}, d.TextBeforeCursor(), true)
}
