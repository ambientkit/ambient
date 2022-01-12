package helper

import (
	"github.com/c-bata/go-prompt"
	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/lib/requestclient"
)

var (
	// Globals.
	log ambient.AppLogger
	rc  *requestclient.RequestClient
)

// Command -
type Command interface {
	Command() string
	Suggestion() prompt.Suggest
	ArgumentSuggestions() SmartSuggestGroup
	Executer(args []string)
	Completer(d prompt.Document, args []string) []prompt.Suggest
}

// CommandList -
type CommandList struct {
	cmd []Command
}

// NewCommandList -
func NewCommandList() *CommandList {
	return &CommandList{
		cmd: make([]Command, 0),
	}
}

// Add a command to the list.
func (cl *CommandList) Add(c Command) {
	cl.cmd = append(cl.cmd, c)
}

// SetGlobals will set the variables used by the package.
func SetGlobals(l ambient.AppLogger, r *requestclient.RequestClient) {
	log = l
	rc = r
}
