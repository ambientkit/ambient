package main

import (
	"fmt"
	pkglog "log"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/app"
	"github.com/josephspurrier/ambient/plugin/gcpbucketstorage"
	"github.com/josephspurrier/ambient/plugin/zaplogger"
)

var (
	appName    = "amb"
	appVersion = "1.0"
	quit       = prompt.KeyBind{
		Key: prompt.ControlC,
		Fn: func(b *prompt.Buffer) {
			os.Exit(0)
		},
	}
	log           ambient.AppLogger
	pluginsystem  *ambient.PluginSystem
	securestorage *ambient.SecureSite
	plugins       *ambient.PluginLoader
)

func main() {
	plugins = app.Plugins()

	// Create the ambient app.
	ambientApp, err := ambient.NewApp(appName, appVersion,
		zaplogger.New(),
		gcpbucketstorage.New(app.StorageSitePath, app.StorageSessionPath),
		plugins)
	if err != nil {
		pkglog.Fatalln(err.Error())
	}

	// Get the
	log = ambientApp.Logger()
	pluginsystem = ambientApp.PluginSystem()

	// Create secure site for the core application and use "ambient" so it gets
	// full permissions.
	securestorage = ambient.NewSecureSite("ambient", log, pluginsystem, nil, nil, nil)

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

var (
	execEnable = "enable"
	execGrants = "grant"
	execExit   = "exit"
)

func enablePlugin(name string) {
	log.Info("enabling plugin: %v", name)
	err := securestorage.EnablePlugin(name, false)
	if err != nil {
		log.Error("", err.Error())
	}
}

func enableGrants(name string) {
	log.Info("add plugin grants: %v", name)

	p, err := pluginsystem.Plugin(name)
	if err != nil {
		log.Error("error with plugin (%v): %v", name, err.Error())
		return
	}

	for _, request := range p.GrantRequests() {
		log.Info("%v - add grant: %v", name, request.Grant)
		err := securestorage.SetNeighborPluginGrant(name, request.Grant, true)
		if err != nil {
			log.Error("", err.Error())
		}
	}
}

func addGrantAll(name string) error {
	// Set the grants for the CLI tool.
	err := pluginsystem.SetGrant(name, ambient.GrantAll)
	if err != nil {
		return err
	}
	return pluginsystem.Save()
}

func removeGrantAll(name string) error {
	// Remove the grants for the CLI tool.
	err := pluginsystem.RemoveGrant(name, ambient.GrantAll)
	if err != nil {
		return err
	}

	return pluginsystem.Save()
}

func pluginSuggestions() []prompt.Suggest {
	arr := make([]prompt.Suggest, 0)
	arr = append(arr, prompt.Suggest{Text: "all", Description: ""})

	for pluginName, trusted := range plugins.TrustedPlugins {
		if trusted {
			arr = append(arr, prompt.Suggest{Text: pluginName, Description: ""})
		}
	}

	return arr
}

func executer(s string) {
	args := strings.Split(s, " ")

	switch args[0] {
	case execEnable:
		if len(args) < 2 {
			log.Info("", "command not recognized")
			break
		}

		log.Info("", "enabling plugin")

		if args[1] == "all" {
			// Enable plugins.
			for pluginName, trusted := range plugins.TrustedPlugins {
				if trusted {
					enablePlugin(pluginName)
				}
			}
		} else {
			enablePlugin(args[1])
		}
	case execGrants:
		if len(args) < 2 {
			log.Info("", "command not recognized")
			break
		}

		log.Info("", "adding plugin grants")

		if args[1] == "all" {
			// Enable plugin grants.
			for pluginName, trusted := range plugins.TrustedPlugins {
				if trusted {
					enableGrants(pluginName)
				}
			}
		} else {
			enableGrants(args[1])
		}
	case execExit:
		os.Exit(0)
	default:
		log.Info("", "command not recognized")
	}
}

func completer(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{}

	// if d.TextBeforeCursor() == "" {
	// 	return suggestions
	// }

	args := strings.Split(d.TextBeforeCursor(), " ")

	if len(args) <= 1 {
		return prompt.FilterHasPrefix([]prompt.Suggest{
			{Text: execEnable, Description: "Enable the core plugins"},
			{Text: execGrants, Description: "Add grants for the core plugins"},
			{Text: execExit, Description: "Exit the CLI (or press Ctrl+C)"},
		}, args[0], true)
	}

	switch args[0] {
	case "enable":
		if len(args) == 2 {
			return prompt.FilterHasPrefix(pluginSuggestions(), args[1], true)
		}
	case "grant":
		if len(args) == 2 {
			return prompt.FilterHasPrefix(pluginSuggestions(), args[1], true)
		}
	}

	return prompt.FilterHasPrefix(suggestions, d.TextBeforeCursor(), true)
}

func exitChecker(in string, breakline bool) bool {
	return in == execExit && breakline
}
