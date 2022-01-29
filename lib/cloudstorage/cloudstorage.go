// Package cloudstorage helps detect where an Ambient app is running and provides
// the correct storage plugin.
package cloudstorage

import (
	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/lib/envdetect"
	"github.com/ambientkit/ambient/plugin/storage/awsbucketstorage"
	"github.com/ambientkit/ambient/plugin/storage/azureblobstorage"
	"github.com/ambientkit/ambient/plugin/storage/gcpbucketstorage"
	"github.com/ambientkit/ambient/plugin/storage/localstorage"
)

// StorageBasedOnCloud returns storage engine based on the environment it's
// running in.
func StorageBasedOnCloud(sitePath string, sessionPath string) ambient.StoragePlugin {
	// Select the storage engine for site and session information.
	var storage ambient.StoragePlugin
	if envdetect.RunningLocalDev() {
		storage = localstorage.New(sitePath, sessionPath)
	} else if envdetect.RunningInGoogle() {
		storage = gcpbucketstorage.New(sitePath, sessionPath)
	} else if envdetect.RunningInAWS() {
		storage = awsbucketstorage.New(sitePath, sessionPath)
	} else if envdetect.RunningInAzureFunction() {
		storage = azureblobstorage.New(sitePath, sessionPath)
	} else {
		// Defaulting to local storage.
		storage = localstorage.New(sitePath, sessionPath)
	}

	return storage
}
