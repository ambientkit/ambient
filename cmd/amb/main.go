package main

import (
	"fmt"
	syslog "log"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/app"
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
	log           ambient.IAppLogger
	pluginsystem  *ambient.PluginSystem
	securestorage *ambient.SecureSite
	plugins       ambient.IPluginList
)

func main() {
	// Ensure there is at least the logger and storage plugins.
	plugins = app.Plugins()
	if len(plugins) < 2 {
		syslog.Fatalln("boot: no log and storage plugins found")
	}

	// Set up the logger.
	var err error
	log, err = ambient.LoadLogger(appName, appVersion, plugins[0])
	if err != nil {
		syslog.Fatalln(err.Error())
	}

	// Get the plugins and initialize storage.
	storage, _, err := ambient.LoadStorage(log, plugins[1])
	if err != nil {
		log.Fatal("", err.Error())
	}

	// Initialize the plugin system.
	pluginsystem, err = ambient.NewPluginSystem(log, plugins, storage)
	if err != nil {
		log.Fatal("", err.Error())
	}

	// Set up the secure storage.
	securestorage = ambient.NewSecureSite(appName, log, storage, pluginsystem, nil, nil, nil)

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
		log.Error("", err.Error())
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

func enableCLIGrant() bool {
	// Initialize the plugin in storage.
	err := pluginsystem.InitializePlugin(appName)
	if err != nil {
		log.Error("could not initialize plugin %v: %v", appName, err.Error())
		return true
	}

	// Add admin grant for the CLI ambient.
	err = addGrantAll(appName)
	if err != nil {
		log.Error("could not enable GrantAll on plugin %v: %v", appName, err.Error())
		return true
	}
	log.Info("temporarily enabling GrantAll for plugin: %v", appName)

	return false
}

func disableCLIGrant() {
	// Remove admin grant for the CLI ambient.
	log.Info("remove GrantAll for plugin: %v", appName)
	err := removeGrantAll(appName)
	if err != nil {
		log.Error("could not remove GrantAll grant from plugin %v: %v", appName, err.Error())
	}

	// Remove the plugin from storage.
	err = pluginsystem.RemovePlugin(appName)
	if err != nil {
		log.Error("could not remove plugin %v: %v", appName, err.Error())
	}

}

func pluginSuggestions() []prompt.Suggest {
	arr := make([]prompt.Suggest, 0)
	arr = append(arr, prompt.Suggest{Text: "all", Description: ""})

	for _, pluginName := range app.MinimalPlugins {
		arr = append(arr, prompt.Suggest{Text: pluginName, Description: ""})
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

		// Enable grants temporarily.
		fail := enableCLIGrant()
		if fail {
			return
		}

		if args[1] == "all" {
			// Enable plugins.
			for _, p := range app.MinimalPlugins {
				enablePlugin(p)
			}
		} else {
			enablePlugin(args[1])
		}

		// Remove temporary grants.
		disableCLIGrant()
	case execGrants:
		if len(args) < 2 {
			log.Info("", "command not recognized")
			break
		}

		log.Info("", "adding plugin grants")

		// Enable grants temporarily.
		fail := enableCLIGrant()
		if fail {
			return
		}

		if args[1] == "all" {
			// Enable plugin grants.
			for _, p := range app.MinimalPlugins {
				enableGrants(p)
			}
		} else {
			enableGrants(args[1])
		}

		// Remove temporary grants.
		disableCLIGrant()
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
