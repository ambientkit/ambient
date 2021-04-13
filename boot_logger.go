package ambient

import (
	"fmt"
)

// LoadLogger returns the logger.
func LoadLogger(appName string, appVersion string, plugin IPlugin) (IAppLogger, error) {
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
