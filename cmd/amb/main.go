package main

import (
	"fmt"
	stdlog "log"
	"os"

	"github.com/c-bata/go-prompt"
	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/cmd/amb/helper"
	"github.com/josephspurrier/ambient/lib/envdetect"
	"github.com/josephspurrier/ambient/lib/requestclient"
	"github.com/josephspurrier/ambient/plugin/logger/zaplogger"
)

var (
	// App information.
	appName    = "amb"
	appVersion = "1.0"

	// Key bindings.
	quit = prompt.KeyBind{
		Key: prompt.ControlC,
		Fn: func(b *prompt.Buffer) {
			os.Exit(0)
		},
	}
)

func main() {
	// Use an Ambient logger for consistency.
	log, err := ambient.NewAppLogger(appName, appVersion, zaplogger.New())
	if err != nil {
		if log != nil {
			// Use the logger if it's available.
			log.Fatal(err.Error())
		} else {
			// Else use the standard logger.
			stdlog.Fatalln(err.Error())
		}
	}

	// Set the URL for the Dev Console.
	rc := requestclient.New(
		fmt.Sprintf("%v:%v", envdetect.DevConsoleURL(), envdetect.DevConsolePort()),
		"")

	// TODO: Should make this a struct instead.
	helper.SetGlobals(log, rc)

	cmds := helper.NewCommandList()
	cmds.Add(&helper.CmdCreateApp{})

	// Start the read–eval–print loop (REPL).
	p := prompt.New(
		cmds.Executer,
		cmds.Completer,
		prompt.OptionTitle(fmt.Sprintf("%v: Ambient Interactive Client", appName)),
		prompt.OptionPrefix(">>> "),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionSetExitCheckerOnInput(helper.ExitChecker),
		prompt.OptionAddKeyBind(quit),
	)
	p.Run()
}
