package helper

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

// Completer handles the auto completion.
func Completer(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{}

	// if d.TextBeforeCursor() == "" {
	// 	return suggestions
	// }

	// Split arguments by spaces.
	args := strings.Split(d.TextBeforeCursor(), " ")

	if len(args) <= 1 {
		return prompt.FilterHasPrefix(promptSuggestions, args[0], true)
	}

	switch args[0] {
	case execEnable, execGrants:
		// For these commands, show a secondary list of plugin suggestions.
		if len(args) == 2 {
			return prompt.FilterHasPrefix(pluginSuggestions(), args[1], true)
		}
	case execCreateApp:
		// Don't show any suggestions if type types: --parameter SPACE
		// TODO: This should probably go to the top.
		prevCursor := d.GetWordBeforeCursorWithSpace()
		if strings.HasPrefix(prevCursor, "--") && strings.HasSuffix(prevCursor, " ") {
			return nil
		}

		// Remove duplicates from autocomplete if they've already been typed in.
		list := filterAlreadyUsed(createAppSuggest.Suggest(), args)

		// Only show autocomplete when the word matches.
		list = prompt.FilterHasPrefix(list, d.GetWordBeforeCursor(), true)

		return list
	}

	return prompt.FilterHasPrefix(suggestions, d.TextBeforeCursor(), true)
}
