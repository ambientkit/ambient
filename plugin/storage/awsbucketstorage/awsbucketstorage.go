// Package awsbucketstorage is an Ambient plugin that provides storage in AWS S3.
package awsbucketstorage

import (
	"os"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/plugin/storage/awsbucketstorage/store"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit

	sitePath    string
	sessionPath string
}

// New returns an Ambient plugin that provides storage in AWS S3.
func New(sitePath string, sessionPath string) *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},

		sitePath:    sitePath,
		sessionPath: sessionPath,
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "awsbucketstorage"
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
	// Bucket allows user to set the AWS bucket.
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

// Storage returns data and session storage.
func (p *Plugin) Storage(logger ambient.Logger) (ambient.DataStorer, ambient.SessionStorer, error) {
	bucket := os.Getenv("AMB_AWS_BUCKET")
	if len(bucket) == 0 {
		var err error
		bucket, err = p.Site.PluginSettingString(Bucket)
		if err != nil {
			return nil, nil, err
		}
	}

	ds := store.NewAWSStorage(bucket, p.sitePath)
	ss := store.NewAWSStorage(bucket, p.sessionPath)

	return ds, ss, nil
}
