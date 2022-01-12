package helper

import (
	"fmt"
	"os"
	"strings"
)

// ExitChecker -
func ExitChecker(in string, breakline bool) bool {
	return in == execExit && breakline
}

// Executer runs the logic.
func Executer(s string) {
	// Split and remove empty items.
	args := filterString(strings.Split(s, " "), "")

	switch args[0] {
	case execCreateApp:
		if valid, missing := createAppSuggest.Valid(args); !valid {
			log.Error("amb: missing argument: %v", missing)
		}
		//log.Info("amb: args: %#v %v", args, len(args))
		log.Info("amb: creating new project")
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
