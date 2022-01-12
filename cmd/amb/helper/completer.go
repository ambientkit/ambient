package helper

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

// Completer handles the auto completion.
func (cl *CommandList) Completer(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{}

	// if d.TextBeforeCursor() == "" {
	// 	return suggestions
	// }

	// Split arguments by spaces.
	args := strings.Split(d.TextBeforeCursor(), " ")

	if len(args) <= 1 {
		return prompt.FilterHasPrefix(cl.InitialCommandSuggestions(), args[0], true)
	}

	firstCommand := args[0]
	for _, v := range cl.cmd {
		if firstCommand == v.Command() {
			return v.Completer(d, args)
		}
	}

	// switch args[0] {
	// case execEnable, execGrants:
	// 	// For these commands, show a secondary list of plugin suggestions.
	// 	if len(args) == 2 {
	// 		return prompt.FilterHasPrefix(pluginSuggestions(), args[1], true)
	// 	}
	// case execCreateApp:
	// 	return createAppCompeter(d, args)
	// }

	return prompt.FilterHasPrefix(suggestions, d.TextBeforeCursor(), true)
}
