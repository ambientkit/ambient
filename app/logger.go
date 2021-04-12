package app

import (
	"fmt"

	"github.com/josephspurrier/ambient/app/core"
)

// Logger returns the logger.
func Logger(appName string, appVersion string, plugin core.IPlugin) (core.IAppLogger, error) {
	// Get the logger from the plugins.
	log, err := plugin.Logger(appName, appVersion)
	if err != nil {
		log.Error("", err.Error())
	} else if log != nil {
		log.Info("boot: using logger from first plugin: %v", plugin.PluginName())
	}
	if log == nil {
		return nil, fmt.Errorf("boot: no default logger found")
	}

	return log, nil
}
