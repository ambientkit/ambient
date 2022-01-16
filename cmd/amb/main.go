package main

import (
	"flag"
	"fmt"
	stdlog "log"
	"os"

	"github.com/c-bata/go-prompt"
	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/cmd/amb/internal"
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
	// Detect debug flag.
	var debugEnable bool
	flag.BoolVar(&debugEnable, "debug", false, "Enable debug output")

	flag.Parse()

	// Determine log level.
	logLevel := ambient.LogLevelInfo
	if debugEnable {
		logLevel = ambient.LogLevelDebug
	}

	// Use an Ambient logger for consistency.
	log, err := ambient.NewAppLogger(appName, appVersion, zaplogger.New(), logLevel)
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
	internal.SetGlobals(log, rc)

	// Get the exit command.
	exit := &internal.CmdExit{}

	cmds := internal.NewCommandList()
	cmds.Add(&internal.CmdCreateApp{})
	cmds.Add(&internal.CmdEnable{})
	cmds.Add(&internal.CmdGrant{})
	cmds.Add(&internal.CmdEncrypt{})
	cmds.Add(&internal.CmdDecrypt{})
	cmds.Add(exit)

	// Start the read–eval–print loop (REPL).
	p := prompt.New(
		cmds.Executer,
		cmds.Completer,
		prompt.OptionTitle(fmt.Sprintf("%v: Ambient Interactive Client", appName)),
		prompt.OptionPrefix(">>> "),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionSetExitCheckerOnInput(exit.Checker),
		prompt.OptionAddKeyBind(quit),
	)
	p.Run()
}
