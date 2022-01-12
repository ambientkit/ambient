package helper

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

// CmdCreateApp -
type CmdCreateApp struct{}

// Command -
func (c *CmdCreateApp) Command() string {
	return "createapp"
}

// Suggestion -
func (c *CmdCreateApp) Suggestion() prompt.Suggest {
	return prompt.Suggest{Text: c.Command(), Description: "Create new Ambient app"}
}

// ArgumentSuggestions -
func (c *CmdCreateApp) ArgumentSuggestions() SmartSuggestGroup {
	return SmartSuggestGroup{
		{Suggest: prompt.Suggest{Text: "--folder", Description: "Folder to create the project. Ex: . or ./newdir"}, Required: true},
		{Suggest: prompt.Suggest{Text: "--template", Description: "Template project to use. Ex: default or github.com/josephspurrier/template"}, Required: true},
	}
}

// Executer -
func (c *CmdCreateApp) Executer(args []string) {
	if valid, missing := c.ArgumentSuggestions().Valid(args); !valid {
		log.Error("amb: missing argument: %v", missing)
	}
	//log.Info("amb: args: %#v %v", args, len(args))
	log.Info("amb: creating new project")
}

// Completer -
func (c *CmdCreateApp) Completer(d prompt.Document, args []string) []prompt.Suggest {
	// Don't show any suggestions if type types: --parameter SPACE
	// TODO: This should probably go to the top.
	prevCursor := d.GetWordBeforeCursorWithSpace()
	if strings.HasPrefix(prevCursor, "--") && strings.HasSuffix(prevCursor, " ") {
		return nil
	}

	// Remove duplicates from autocomplete if they've already been typed in.
	list := filterAlreadyUsed(c.ArgumentSuggestions().ToSuggest(), args)

	// Only show autocomplete when the word matches.
	list = prompt.FilterHasPrefix(list, d.GetWordBeforeCursor(), true)

	return list
}
