// Package gcpbucketstorage is an Ambient plugin that provides storage in GCP Cloud Storage.
package gcpbucketstorage

import (
	"fmt"
	"os"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/plugin/storage/gcpbucketstorage/store"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	sitePath    string
	sessionPath string
}

// New returns an Ambient plugin that provides storage in GCP Cloud Storage.
func New(sitePath string, sessionPath string) *Plugin {
	return &Plugin{
		sitePath:    sitePath,
		sessionPath: sessionPath,
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

const (
	// BucketEnv is the Google Storage bucket environment variable.
	BucketEnv = "AMB_GCP_BUCKET"
)

// Storage returns data and session storage.
func (p *Plugin) Storage(logger ambient.Logger) (ambient.DataStorer, ambient.SessionStorer, error) {
	// Load the bucket from the environment variable.
	bucket := os.Getenv(BucketEnv)
	if len(bucket) == 0 {
		return nil, nil, fmt.Errorf("%v: environment variable, %v, is missing", p.PluginName(), BucketEnv)
	}

	ds := store.NewGCPStorage(bucket, p.sitePath)
	ss := store.NewGCPStorage(bucket, p.sessionPath)

	return ds, ss, nil
}
