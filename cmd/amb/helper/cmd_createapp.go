package helper

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/c-bata/go-prompt"
)

const (
	paramFolder   = "--folder"
	paramTemplate = "--template"
)

// CmdCreateApp represents a command object.
type CmdCreateApp struct {
	CmdBase
}

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
		{Suggest: prompt.Suggest{Text: paramFolder, Description: "Folder to create the project (default: myambapp)"}, Required: true},
		{Suggest: prompt.Suggest{Text: paramTemplate, Description: "Template project to git clone (default: https://github.com/josephspurrier/ball)"}, Required: false},
	}
}

// param returns the named parameter value.
func (c *CmdCreateApp) param(args []string, name string) (string, error) {
	for k, v := range args {
		if v == name {
			if len(args) >= k+2 {
				return args[k+1], nil
			}
		}
	}

	return "", fmt.Errorf("param not found: %v", name)
}

// git clone --depth=1 --branch=master git@github.com:josephspurrier/ball.git .

// Executer executes the command.
func (c *CmdCreateApp) Executer(args []string) {
	// Ensure all required arguments exist.
	// TODO: Maybe move this up a level so it's called automatically?
	if valid, missing := c.ArgumentSuggestions().Valid(args); !valid {
		log.Error("amb: missing argument: %v", missing)
		return
	}

	// Get folder name.
	folderName, err := c.param(args, paramFolder)
	if err != nil {
		folderName = "myambapp"
	}

	// Determine if folder already exists.
	if _, err := os.Stat(folderName); !os.IsNotExist(err) {
		log.Error("amb: folder already exists: %v", folderName)
		return
	}

	// Get template name.
	templateName, err := c.param(args, paramTemplate)
	if err != nil {
		templateName = "https://github.com/josephspurrier/ball"
	}

	// Perform git clone on the template.
	log.Info("amb: creating new project from template: %v", templateName)
	gitArgs := []string{"clone", "--depth=1", "--branch=master", templateName, folderName}
	cmd := exec.Command("git", gitArgs...)
	var out bytes.Buffer
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		log.Error("amb: couldn't create project (git %v): %v | %v", strings.Join(gitArgs, " "), err.Error(), out.String())
		return
	}

	// Remove .git folder.
	gitFolder := filepath.Join(folderName, ".git")
	err = os.RemoveAll(gitFolder)
	if err != nil {
		log.Error("amb: couldn't remove .git folder: %v", err.Error())
	}

	log.Info("amb: removing folder: %v", gitFolder)
	log.Info("amb: created project successfully in folder: %v", folderName)
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
