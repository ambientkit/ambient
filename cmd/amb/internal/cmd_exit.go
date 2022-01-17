package internal

import (
	"os"

	"github.com/c-bata/go-prompt"
)

// CmdExit represents a command object.
type CmdExit struct {
	CmdBase
}

// Command returns the initial command.
func (c *CmdExit) Command() string {
	return "exit"
}

// Suggestion returns the suggestion for the initial command.
func (c *CmdExit) Suggestion() prompt.Suggest {
	return prompt.Suggest{Text: c.Command(), Description: "Exit the CLI (or press Ctrl+D)"}
}

// Executer executes the command.
func (c *CmdExit) Executer(args []string) {
	os.Exit(0)
}

// Checker returns true if exiting.
func (c *CmdExit) Checker(in string, breakline bool) bool {
	return in == c.Command() && breakline
}
