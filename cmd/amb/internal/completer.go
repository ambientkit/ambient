package internal

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

// Completer handles the auto completion.
func (cl *CommandList) Completer(d prompt.Document) []prompt.Suggest {
	// Split arguments by spaces.
	args := strings.Split(d.TextBeforeCursor(), " ")

	// If there is no argument, then show a list of suggestions on initial tab.
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(cl.InitialCommandSuggestions(), args[0], true)
	}

	// Loop through each command to find a match to suggest.
	firstCommand := args[0]
	for _, v := range cl.cmd {
		if firstCommand == v.Command() {
			return v.Completer(d, args)
		}
	}

	// Return no suggestions.
	return prompt.FilterHasPrefix([]prompt.Suggest{}, d.TextBeforeCursor(), true)
}
