package internal

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
	paramFolder     = "-folder"
	paramTemplate   = "-template"
	defaultTemplate = "https://github.com/ambientkit/ambient-template"
	defaultFolder   = "ambapp"
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
		{Suggest: prompt.Suggest{Text: paramFolder, Description: fmt.Sprintf("Folder to create the project (default: %v)", defaultFolder)}, Required: false},
		{Suggest: prompt.Suggest{Text: paramTemplate, Description: fmt.Sprintf("Template project to git clone (default: %v)", defaultTemplate)}, Required: false},
	}
}

// Executer executes the command.
func (c *CmdCreateApp) Executer(args []string) {
	// Get folder name.
	folderName, err := c.Param(args, paramFolder)
	if err != nil {
		folderName = defaultFolder
	}

	// Determine if folder already exists.
	if _, err := os.Stat(folderName); !os.IsNotExist(err) {
		log.Error("amb: folder already exists: %v", folderName)
		return
	}

	// Get template name.
	templateName, err := c.Param(args, paramTemplate)
	if err != nil {
		templateName = defaultTemplate
	}

	// Perform git clone on the template.
	log.Info("amb: creating new project from template: %v", templateName)
	gitArgs := []string{"clone", "--depth=1", "--branch=main", templateName, folderName}
	cmd := exec.Command("git", gitArgs...)
	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr
	err = cmd.Run()
	if err != nil {
		log.Error("amb: couldn't create project (git %v): %v %v", strings.Join(gitArgs, " "), err.Error(), stdErr.String())
		return
	}

	// Remove .git folder.
	gitFolder := filepath.Join(folderName, ".git")
	err = os.RemoveAll(gitFolder)
	if err != nil {
		log.Error("amb: couldn't remove .git folder: %v", err.Error())
	}

	// Make bin folder.
	binFolder := filepath.Join(folderName, "bin")
	err = os.Mkdir(binFolder, 0755)
	if err != nil {
		log.Error("amb: couldn't create bin folder: %v", err.Error())
	}

	log.Info("amb: removing folder: %v", gitFolder)
	log.Info("amb: created project successfully in folder: %v", folderName)
}

// Completer returns a list of suggestions based on the user input.
func (c *CmdCreateApp) Completer(d prompt.Document, args []string) []prompt.Suggest {
	// Don't show any suggestions if type types: --parameter SPACE
	prevCursor := d.GetWordBeforeCursorWithSpace()
	if strings.HasPrefix(prevCursor, "-") && strings.HasSuffix(prevCursor, " ") {
		return nil
	}

	// Remove duplicates from autocomplete if they've already been typed in.
	list := filterAlreadyUsed(c.ArgumentSuggestions().ToSuggest(), args)

	// Only show autocomplete when the word matches.
	list = prompt.FilterHasPrefix(list, d.GetWordBeforeCursor(), true)

	return list
}
