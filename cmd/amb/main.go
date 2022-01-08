package main

import (
	"fmt"
	stdlog "log"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/lib/requestclient"
	"github.com/josephspurrier/ambient/plugin/logger/zaplogger"
)

var (
	// App information.
	appName    = "amb"
	appVersion = "1.0"

	// Commands available.
	execEnable  = "enable"
	execGrants  = "grant"
	execEncrypt = "encryptstorage"
	execDecrypt = "decryptstorage"
	execExit    = "exit"

	// Prompts should match 1:1 with the commands above.
	promptSuggestions = []prompt.Suggest{
		{Text: execEnable, Description: "Enable plugin..."},
		{Text: execGrants, Description: "Add grants for plugin..."},
		{Text: execEncrypt, Description: "Encrypt storage"},
		{Text: execDecrypt, Description: "Decrypt storage"},
		{Text: execExit, Description: "Exit the CLI (or press Ctrl+C)"},
	}

	// Key bindings.
	quit = prompt.KeyBind{
		Key: prompt.ControlC,
		Fn: func(b *prompt.Buffer) {
			os.Exit(0)
		},
	}

	// Globals.
	log ambient.AppLogger
	rc  *requestclient.RequestClient
)

func main() {
	// Use an Ambient logger for consistency.
	var err error
	log, err = ambient.NewAppLogger(appName, appVersion, zaplogger.New())
	if err != nil {
		if log != nil {
			// Use the logger if it's available.
			log.Fatal(err.Error())
		} else {
			// Else use the standard logger.
			stdlog.Fatalln(err.Error())
		}
	}

	// TODO: Need to make this configurable.
	rc = requestclient.New("http://localhost:8081", "")

	// Start the read–eval–print loop (REPL).
	p := prompt.New(
		executer,
		completer,
		prompt.OptionTitle(fmt.Sprintf("%v: Ambient Interactive Client", appName)),
		prompt.OptionPrefix(">>> "),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionSetExitCheckerOnInput(exitChecker),
		prompt.OptionAddKeyBind(quit),
	)
	p.Run()
}

func exitChecker(in string, breakline bool) bool {
	return in == execExit && breakline
}

func executer(s string) {
	args := strings.Split(s, " ")

	switch args[0] {
	case execEnable:
		if len(args) < 2 {
			log.Info("amb: command not recognized")
			break
		}

		if args[1] == "all" {
			// Enable all plugins.
			log.Info("amb: enabling all trusted plugins")

			err := rc.Post("/plugins/enable", nil, nil)
			if err != nil {
				log.Error("amb: could not enable all plugins: %v", err.Error())
			}
		} else {
			// Enable one plugin.
			pluginName := args[1]
			log.Info("amb: enabling plugin: %v", pluginName)

			err := rc.Post(fmt.Sprintf("/plugins/%v/enable", pluginName), nil, nil)
			if err != nil {
				log.Error("amb: could not enable plugin, %v: %v", pluginName, err.Error())
			}
		}
	case execGrants:
		if len(args) < 2 {
			log.Info("amb: command not recognized")
			break
		}

		if args[1] == "all" {
			// Enable grants for all plugins.
			log.Info("amb: adding grants for all trusted plugins")

			err := rc.Post("/plugins/grant", nil, nil)
			if err != nil {
				log.Error("amb: cloud not enable all plugins grants: %v", err.Error())
			}
		} else {
			// Enable grants for one plugin.
			pluginName := args[1]
			log.Info("amb: adding grants for plugin: %v", pluginName)

			err := rc.Post(fmt.Sprintf("/plugins/%v/grant", pluginName), nil, nil)
			if err != nil {
				log.Error("amb: cloud not enable plugin (%v) grants: %v", pluginName, err.Error())
			}
		}
	case execEncrypt:
		err := rc.Post("/storage/encrypt", nil, nil)
		if err != nil {
			log.Error("amb: error encrypting storage: %v", err)
		} else {
			log.Info("amb: encrypted storage file: site.bin")
		}
	case execDecrypt:
		err := rc.Post("/storage/decrypt", nil, nil)
		if err != nil {
			log.Error("amb: error decrypted storage: %v", err)
		} else {
			log.Info("amb: decrypted storage file: site.bin")
		}
	case execExit:
		os.Exit(0)
	default:
		log.Info("amb: command not recognized")
	}
}

// completer handles the auto completion.
func completer(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{}

	// if d.TextBeforeCursor() == "" {
	// 	return suggestions
	// }

	args := strings.Split(d.TextBeforeCursor(), " ")

	if len(args) <= 1 {
		return prompt.FilterHasPrefix(promptSuggestions, args[0], true)
	}

	switch args[0] {
	case execEnable, execGrants:
		// For these commands, show a secondary list of plugin suggestions.
		if len(args) == 2 {
			return prompt.FilterHasPrefix(pluginSuggestions(), args[1], true)
		}
	}

	return prompt.FilterHasPrefix(suggestions, d.TextBeforeCursor(), true)
}

// pluginSuggestions returns a list of suggestions for plugins.
func pluginSuggestions() []prompt.Suggest {
	arr := make([]prompt.Suggest, 0)
	arr = append(arr, prompt.Suggest{Text: "all", Description: ""})

	// Get the plugin names.
	pluginNames := make([]string, 0)
	err := rc.Get("/plugins", &pluginNames)
	if err != nil {
		log.Error("amb: could not get plugin names: %v", err.Error())
		return nil
	}

	// Build a list of suggestions.
	for _, pluginName := range pluginNames {
		arr = append(arr, prompt.Suggest{Text: pluginName, Description: ""})
	}

	return arr
}
