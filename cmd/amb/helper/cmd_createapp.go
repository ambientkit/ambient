package helper

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

// CmdCreateApp represents a command object.
type CmdCreateApp struct{}

// Command returns the initial command.
func (c *CmdCreateApp) Command() string {
	return "createapp"
}

// Suggestion returns the suggestion for the initial command.
func (c *CmdCreateApp) Suggestion() prompt.Suggest {
	return prompt.Suggest{Text: c.Command(), Description: "Create new Ambient app..."}
}

// ArgumentSuggestions returns a smart suggestion group that includes validation.
func (c *CmdCreateApp) ArgumentSuggestions() SmartSuggestGroup {
	return SmartSuggestGroup{
		{Suggest: prompt.Suggest{Text: "--folder", Description: "Folder to create the project. Ex: . or ./newdir"}, Required: true},
		{Suggest: prompt.Suggest{Text: "--template", Description: "Template project to use. Ex: default or github.com/josephspurrier/template"}, Required: true},
	}
}

// Executer executes the command.
func (c *CmdCreateApp) Executer(args []string) {
	if valid, missing := c.ArgumentSuggestions().Valid(args); !valid {
		log.Error("amb: missing argument: %v", missing)
	}
	//log.Info("amb: args: %#v %v", args, len(args))
	log.Info("amb: creating new project")
}

// Completer returns a list of suggestions based on the user input.
func (c *CmdCreateApp) Completer(d prompt.Document, args []string) []prompt.Suggest {
	// Don't show any suggestions if type types: --parameter SPACE
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
