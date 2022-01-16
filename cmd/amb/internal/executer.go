package internal

import (
	"strings"
)

// Executer executes the typed in command.
func (cl *CommandList) Executer(s string) {
	// Split and remove empty items.
	args := filterString(strings.Split(s, " "), "")

	firstCommand := args[0]

	// Loop through each command to find a match to execute.
	handled := false
	for _, v := range cl.cmd {
		if firstCommand == v.Command() {
			// Ensure all required arguments exist.
			if valid, missing := v.ArgumentSuggestions().Valid(args); !valid {
				log.Error("amb: missing argument for '%v': %v", firstCommand, missing)
				return
			}
			// Execute the command.
			v.Executer(args)
			handled = true
			break
		}
	}

	if !handled {
		log.Info("amb: command not recognized")
	}
}
