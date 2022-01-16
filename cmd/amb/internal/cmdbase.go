package internal

import (
	"fmt"

	"github.com/c-bata/go-prompt"
)

// CmdBase is a base object for structs. This reduces creating methods that are
// optional and provides methods all can be used.
type CmdBase struct{}

// ArgumentSuggestions returns a smart suggestion group that includes validation.
func (c *CmdBase) ArgumentSuggestions() SmartSuggestGroup {
	return SmartSuggestGroup{}
}

// Completer returns a list of suggestions based on the user input.
func (c *CmdBase) Completer(d prompt.Document, args []string) []prompt.Suggest {
	// Return nothing.
	return prompt.FilterHasPrefix([]prompt.Suggest{}, d.TextBeforeCursor(), true)
}

// Param returns the named parameter value or an error if it doesn't exist.
func (c *CmdBase) Param(args []string, name string) (string, error) {
	for k, v := range args {
		if v == name {
			if len(args) >= k+2 {
				return args[k+1], nil
			}
		}
	}

	return "", fmt.Errorf("param not found: %v", name)
}
