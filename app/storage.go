package app

import (
	"fmt"

	"github.com/josephspurrier/ambient/app/core"
	"github.com/josephspurrier/ambient/app/lib/logger"
)

// Storage returns the storage.
func Storage(l *logger.Logger, arrPlugins core.IPluginList) (*core.Storage, core.SessionStorer, error) {
	// Define the storage managers.
	var ds core.DataStorer
	var ss core.SessionStorer

	// Get the storage manager from the plugins.
	// This must be the first plugin or else it fails.
	if len(arrPlugins) > 0 {
		firstPlugin := arrPlugins[0]
		// Get the storage system.
		pds, pss, err := firstPlugin.Storage(l)
		if err != nil {
			l.Error("", err.Error())
		} else if pds != nil && pss != nil {
			l.Info("boot: using storage from first plugin: %v", firstPlugin.PluginName())
			ds = pds
			ss = pss
		}
	}
	if ds == nil || ss == nil {
		return nil, nil, fmt.Errorf("boot: no default storage found")
	}

	// Create new store object with the defaults.
	site := &core.Site{}

	// Set up the data storage provider.
	storage, err := core.NewDatastore(ds, site)
	if err != nil {
		return nil, nil, err
	}

	return storage, ss, err
}
