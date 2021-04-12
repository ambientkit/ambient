package main

import (
	"fmt"
	"os"

	"github.com/c-bata/go-prompt"
	"github.com/josephspurrier/ambient/app"
	"github.com/josephspurrier/ambient/app/core"
	"github.com/josephspurrier/ambient/app/lib/logger"
	"github.com/sirupsen/logrus"
)

var (
	appName = "amb"
	quit    = prompt.KeyBind{
		Key: prompt.ControlC,
		Fn: func(b *prompt.Buffer) {
			os.Exit(0)
		},
	}
	log           *logger.Logger
	plugins       *core.PluginSystem
	securestorage *core.SecureSite
)

func init() {
	// Set the time zone.
	tz := os.Getenv("AMB_TIMEZONE")
	if len(tz) == 0 {
		// Set the default to eastern time.
		tz = "America/New_York"
	}
	os.Setenv("TZ", tz)
}

func main() {
	// Create the logger.
	log = logger.NewLogger(appName, "1.0")
	//log.SetLevel(uint32(logrus.DebugLevel))
	log.SetLevel(uint32(logrus.InfoLevel))

	// Ensure there is at least the storage plugin.
	if len(app.Plugins) == 0 {
		log.Fatal("", "boot: no plugins found")
	}

	// Get the plugins and initialize storage.
	storage, _, err := app.Storage(log, app.Plugins[0])
	if err != nil {
		log.Fatal("", err.Error())
	}

	// Initialize the plugin system.
	plugins, err = core.NewPluginSystem(log, app.Plugins, storage)
	if err != nil {
		log.Fatal("", err.Error())
	}

	// Set up the secure storage.
	securestorage = core.NewSecureSite(appName, log, storage, nil, nil, plugins)

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
	err := securestorage.EnablePlugin(name)
	if err != nil {
		log.Error("", err.Error())
	}
}

func enableGrants(name string) {
	log.Info("add plugin grants: %v", name)

	p, err := plugins.Plugin(name)
	if err != nil {
		log.Error("", err.Error())
	}

	for _, v := range p.Grants() {
		log.Info("%v - add grant: %v", name, v)
		err := securestorage.SetNeighborPluginGrant(name, v, true)
		if err != nil {
			log.Error("", err.Error())
		}
	}
}

func addGrantAll(name string) error {
	// Set the grants for the CLI tool.
	err := plugins.SetGrant(name, core.GrantAll)
	if err != nil {
		return err
	}
	return plugins.Save()
}

func removeGrantAll(name string) error {
	// Remove the grants for the CLI tool.
	err := plugins.RemoveGrant(name, core.GrantAll)
	if err != nil {
		return err
	}

	return plugins.Save()
}

func enableCLIGrant() bool {
	// Initialize the plugin in storage.
	err := plugins.InitializePlugin(appName)
	if err != nil {
		log.Error("could not initialize plugin %v: %v", appName, err.Error())
		return true
	}

	// Add admin grant for the CLI app.
	err = addGrantAll(appName)
	if err != nil {
		log.Error("could not enable GrantAll on plugin %v: %v", appName, err.Error())
		return true
	}
	log.Info("temporarily enabling GrantAll for plugin: %v", appName)

	return false
}

func disableCLIGrant() {
	// Remove admin grant for the CLI app.
	log.Info("remove GrantAll for plugin: %v", appName)
	err := removeGrantAll(appName)
	if err != nil {
		log.Error("could not remove GrantAll grant from plugin %v: %v", appName, err.Error())
	}

	// Remove the plugin from storage.
	err = plugins.RemovePlugin(appName)
	if err != nil {
		log.Error("could not remove plugin %v: %v", appName, err.Error())
	}

}

// List of core plugins.
var pluginList = []string{
	"awayrouter",
	"scssession",
	"htmltemplate",
	"plugins",
	"bearcss",
	"bearblog",
}

func executer(s string) {
	switch s {
	case execEnable:
		log.Info("", "enabling core plugins")

		// Enable grants temporarily.
		fail := enableCLIGrant()
		if fail {
			return
		}

		// Enable plugins.
		for _, p := range pluginList {
			enablePlugin(p)
		}

		// Remove temporary grants.
		disableCLIGrant()
	case execGrants:
		log.Info("", "adding core plugin grants")

		// Enable grants temporarily.
		fail := enableCLIGrant()
		if fail {
			return
		}

		// Enable plugin grants.
		for _, p := range pluginList {
			enableGrants(p)
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
	s := []prompt.Suggest{
		{Text: execEnable, Description: "Enable the core plugins"},
		{Text: execGrants, Description: "Add grants for the core plugins"},
		{Text: execExit, Description: "Exit the CLI (or press Ctrl+C)"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func exitChecker(in string, breakline bool) bool {
	return in == execExit && breakline
}
