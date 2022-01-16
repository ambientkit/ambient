package internal

import (
	"github.com/c-bata/go-prompt"
)

// CmdEncrypt represents a command object.
type CmdEncrypt struct {
	CmdBase
}

// Command returns the initial command.
func (c *CmdEncrypt) Command() string {
	return "encryptstorage"
}

// Suggestion returns the suggestion for the initial command.
func (c *CmdEncrypt) Suggestion() prompt.Suggest {
	return prompt.Suggest{Text: c.Command(), Description: "Encrypt storage"}
}

// Executer executes the command.
func (c *CmdEncrypt) Executer(args []string) {
	err := rc.Post("/storage/encrypt", nil, nil)
	if err != nil {
		log.Error("amb: error encrypting storage: %v", err)
	} else {
		log.Info("amb: encrypted storage file: site.bin")
	}
}

// Completer returns a list of suggestions based on the user input.
func (c *CmdEncrypt) Completer(d prompt.Document, args []string) []prompt.Suggest {
	// Return nothing.
	return prompt.FilterHasPrefix([]prompt.Suggest{}, d.TextBeforeCursor(), true)
}
