package ambient

import (
	"fmt"
)

// LoadStorage returns the storage.
func LoadStorage(log IAppLogger, plugin IPlugin) (*Storage, SessionStorer, error) {
	// Define the storage managers.
	var ds IDataStorer
	var ss SessionStorer

	// Get the storage manager from the plugins.
	pds, pss, err := plugin.Storage(log)
	if err != nil {
		log.Error("", err.Error())
	} else if pds != nil && pss != nil {
		log.Info("boot: using storage from first plugin: %v", plugin.PluginName())
		ds = pds
		ss = pss
	}
	if ds == nil || ss == nil {
		return nil, nil, fmt.Errorf("boot: no default storage found")
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
