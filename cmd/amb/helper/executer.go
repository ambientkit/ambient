package helper

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
			v.Executer(args)
			handled = true
			break
		}
	}

	if !handled {
		log.Info("amb: command not recognized")
	}
}
