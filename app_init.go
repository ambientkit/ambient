package ambient

import "fmt"

// LoadLogger returns the logger.
func loadLogger(appName string, appVersion string, plugin LoggingPlugin) (AppLogger, error) {
	// Don't allow certain plugin names.
	if allowed, ok := disallowedPluginNames[plugin.PluginName()]; ok && !allowed {
		return nil, fmt.Errorf("ambient: plugin name not allowed: %v", plugin.PluginName())
	}

	// Get the logger from the plugins.
	log, err := plugin.Logger(appName, appVersion, nil)
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
func loadStorage(log AppLogger, pluginGroup StoragePluginGroup) (*Storage, SessionStorer, error) {
	// Detect if storage plugin is missing.
	if pluginGroup.Storage == nil {
		return nil, nil, fmt.Errorf("ambient: storage plugin is missing")
	}

	plugin := pluginGroup.Storage

	// Don't allow certain plugin names.
	if allowed, ok := disallowedPluginNames[plugin.PluginName()]; ok && !allowed {
		return nil, nil, fmt.Errorf("ambient: plugin name not allowed: %v", plugin.PluginName())
	}

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
	storage, err := NewStorage(log, ds, pluginGroup.Encryption)
	if err != nil {
		return nil, nil, err
	}

	return storage, ss, err
}
