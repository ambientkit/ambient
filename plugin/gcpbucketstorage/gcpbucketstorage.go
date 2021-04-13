// Package gcpbucketstorage provides GCP storage and local when AMB_LOCAL is set
// for an Ambient application.
package gcpbucketstorage

import (
	"fmt"
	"os"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/plugin/gcpbucketstorage/store"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit
}

// New returns a new gcpbucketstorage plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "gcpbucketstorage"
}

// PluginVersion returns the plugin version.
func (p *Plugin) PluginVersion() string {
	return "1.0.0"
}

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *ambient.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

var (
	storageSitePath    = "storage/site.json"
	storageSessionPath = "storage/session.bin"
)

// Storage returns data and session storage.
func (p *Plugin) Storage(logger ambient.ILogger) (ambient.IDataStorer, ambient.SessionStorer, error) {
	bucket := os.Getenv("AMB_GCP_BUCKET_NAME")
	if len(bucket) == 0 {
		return nil, nil, fmt.Errorf("environment variable missing: %v", "AMB_GCP_BUCKET_NAME")
	}

	// Set the storage and session environment variables.
	// sitePath := os.Getenv("AMB_SITE_PATH")
	// if len(sitePath) > 0 {
	// 	storageSitePath = sitePath
	// }

	var ds ambient.IDataStorer
	var ss ambient.SessionStorer

	if RunningLocalDev() {
		// Use local filesytem when developing.
		ds = store.NewLocalStorage(storageSitePath)
		ss = store.NewLocalStorage(storageSessionPath)

	} else {
		// Use Google when running in GCP.
		ds = store.NewGCPStorage(bucket, storageSitePath)
		ss = store.NewGCPStorage(bucket, storageSessionPath)
	}

	return ds, ss, nil
}

// RunningLocalDev returns true if the AMB_LOCAL environment variable is set.
func RunningLocalDev() bool {
	s := os.Getenv("AMB_LOCAL")
	return len(s) > 0
}
