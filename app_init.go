package ambient

import "fmt"

// LoadLogger returns the logger.
func loadLogger(appName string, appVersion string, plugin Plugin) (AppLogger, error) {
	// Get the logger from the plugins.
	log, err := plugin.Logger(appName, appVersion)
	if err != nil {
		return nil, err
	} else if log == nil {
		return nil, fmt.Errorf("ambient: no logger found")
	} else {
		log.Info("ambient: using logger from plugin: %v", plugin.PluginName())
	}

	return log, nil
}

// LoadStorage returns the storage.
func loadStorage(log AppLogger, plugin Plugin) (*Storage, SessionStorer, error) {
	// Define the storage managers.
	var ds DataStorer
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
		return nil, nil, fmt.Errorf("ambient: no storage manager found")
	}

	// Set up the data storage provider.
	storage, err := NewStorage(ds)
	if err != nil {
		return nil, nil, err
	}

	return storage, ss, err
}
