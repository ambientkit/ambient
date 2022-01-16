package internal

import (
	"github.com/c-bata/go-prompt"
)

// Command represents all the information required to run a command, including
// suggestions.
type Command interface {
	Command() string
	Suggestion() prompt.Suggest
	Executer(args []string)
	Completer(d prompt.Document, args []string) []prompt.Suggest
	ArgumentSuggestions() SmartSuggestGroup
}

// CommandList is a collection of commands.
type CommandList struct {
	cmd []Command
}

// NewCommandList returns a new collection of commands.
func NewCommandList() *CommandList {
	return &CommandList{
		cmd: make([]Command, 0),
	}
}

// Add a command to the list.
func (cl *CommandList) Add(c Command) {
	cl.cmd = append(cl.cmd, c)
}

// InitialCommandSuggestions returns a list of the initial or top-level
// commands.
func (cl *CommandList) InitialCommandSuggestions() []prompt.Suggest {
	arr := make([]prompt.Suggest, 0)

	for _, v := range cl.cmd {
		arr = append(arr, v.Suggestion())
	}

	return arr
}
