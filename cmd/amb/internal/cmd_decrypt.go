package internal

import (
	"github.com/c-bata/go-prompt"
)

// CmdDecrypt represents a command object.
type CmdDecrypt struct {
	CmdBase
}

// Command returns the initial command.
func (c *CmdDecrypt) Command() string {
	return "decryptstorage"
}

// Suggestion returns the suggestion for the initial command.
func (c *CmdDecrypt) Suggestion() prompt.Suggest {
	return prompt.Suggest{Text: c.Command(), Description: "Decrypt storage"}
}

// Executer executes the command.
func (c *CmdDecrypt) Executer(args []string) {
	err := rc.Post("/storage/decrypt", nil, nil)
	if err != nil {
		log.Error("amb: error decrypted storage: %v", err)
	} else {
		log.Info("amb: decrypted storage file: site.bin")
	}
}

// Completer returns a list of suggestions based on the user input.
func (c *CmdDecrypt) Completer(d prompt.Document, args []string) []prompt.Suggest {
	// Return nothing.
	return prompt.FilterHasPrefix([]prompt.Suggest{}, d.TextBeforeCursor(), true)
}
