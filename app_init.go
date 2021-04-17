package ambient

import "fmt"

// LoadLogger returns the logger.
func loadLogger(appName string, appVersion string, plugin IPlugin) (IAppLogger, error) {
	// Get the logger from the plugins.
	log, err := plugin.Logger(appName, appVersion)
	if err != nil {
		return nil, err
	} else if log == nil {
		return nil, fmt.Errorf("ambient: no default logger found")
	} else {
		log.Info("ambient: using logger from plugin: %v", plugin.PluginName())
	}

	return log, nil
}

// LoadStorage returns the storage.
func loadStorage(log IAppLogger, plugin IPlugin) (*Storage, SessionStorer, error) {
	// Define the storage managers.
	var ds IDataStorer
	var ss SessionStorer

	// Get the storage manager from the plugins.
	pds, pss, err := plugin.Storage(log)
	if err != nil {
		log.Error("", err.Error())
	} else if pds != nil && pss != nil {
		log.Info("ambient: using storage from first plugin: %v", plugin.PluginName())
		ds = pds
		ss = pss
	}
	if ds == nil || ss == nil {
		return nil, nil, fmt.Errorf("ambient: no default storage found")
	}

	// Create new store object with the defaults.
	site := &Site{}

	// Set up the data storage provider.
	storage, err := NewStorage(ds, site)
	if err != nil {
		return nil, nil, err
	}

	return storage, ss, err
}
