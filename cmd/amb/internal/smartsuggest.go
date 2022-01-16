package internal

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

// SmartSuggest is an advanced struct that has additional properties.
type SmartSuggest struct {
	prompt.Suggest

	Required bool
}

// SmartSuggestGroup is a group.
type SmartSuggestGroup []SmartSuggest

// Valid returns true if all of the required fields are found or false with the
// field that is missing.
func (g SmartSuggestGroup) Valid(list []string) (bool, string) {
	for _, suggestion := range g {
		if suggestion.Required {
			found := false
			missingArg := suggestion.Text
			for k, typed := range list {
				// Determine if the argument is found.
				if typed == suggestion.Text {
					// Determine if there is a value that is not a parameter
					// after it.
					if len(list) >= k+2 && !strings.HasPrefix(list[k+1], "-") {
						found = true
					}
					break
				}
			}

			if !found {
				return false, missingArg
			}
		}
	}

	return true, ""
}

// ToSuggest returns an array of suggestions.
func (g SmartSuggestGroup) ToSuggest() []prompt.Suggest {
	arr := make([]prompt.Suggest, 0)

	for _, v := range g {
		arr = append(arr, v.Suggest)
	}
	return arr
}
