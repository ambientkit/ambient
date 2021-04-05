// Package gcpbucketstorage provides GCP storage and local when AMB_LOCAL is set
// for an Ambient application.
package gcpbucketstorage

import (
	"embed"
	"fmt"
	"os"

	"github.com/josephspurrier/ambient/app/core"
	"github.com/josephspurrier/ambient/plugin/gcpbucketstorage/store"
)

//go:embed *
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginMeta
	*core.Toolkit
}

// New returns a new gcpbucketstorage plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "gcpbucketstorage",
			Version:    "1.0.0",
			AppVersion: "1.0.0",
		},
	}
}

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *core.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

var (
	storageSitePath    = "storage/site.json"
	storageSessionPath = "storage/session.bin"
)

// Storage returns data and session storage.
func (p *Plugin) Storage(logger core.ILogger) (core.DataStorer, core.SessionStorer, error) {
	bucket := os.Getenv("AMB_GCP_BUCKET_NAME")
	if len(bucket) == 0 {
		return nil, nil, fmt.Errorf("environment variable missing: %v", "AMB_GCP_BUCKET_NAME")
	}

	// Set the storage and session environment variables.
	// sitePath := os.Getenv("AMB_SITE_PATH")
	// if len(sitePath) > 0 {
	// 	storageSitePath = sitePath
	// }

	var ds core.DataStorer
	var ss core.SessionStorer

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
