// Package gcpbucketstorage provides GCP storage and local when AMB_LOCAL is set
// for an Ambient application.
package gcpbucketstorage

import (
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

const (
	// Bucket allows user to set the GCP bucket.
	Bucket = "Bucket"
)

// Settings returns a list of user settable fields.
func (p *Plugin) Settings() []ambient.Setting {
	return []ambient.Setting{
		{
			Name: Bucket,
		},
	}
}

var (
	storageSitePath    = "storage/site.json"
	storageSessionPath = "storage/session.bin"
)

// Storage returns data and session storage.
func (p *Plugin) Storage(logger ambient.Logger) (ambient.DataStorer, ambient.SessionStorer, error) {
	var ds ambient.DataStorer
	var ss ambient.SessionStorer

	if runningLocalDev() {
		// Use local filesytem when developing.
		ds = store.NewLocalStorage(storageSitePath)
		ss = store.NewLocalStorage(storageSessionPath)
	} else {
		bucket, err := p.Site.PluginSettingString(Bucket)
		if err != nil {
			return nil, nil, err
		}

		// Use Google when running in GCP.
		ds = store.NewGCPStorage(bucket, storageSitePath)
		ss = store.NewGCPStorage(bucket, storageSessionPath)
	}

	return ds, ss, nil
}

// runningLocalDev returns true if the AMB_LOCAL environment variable is set.
func runningLocalDev() bool {
	s := os.Getenv("AMB_LOCAL")
	return len(s) > 0
}
