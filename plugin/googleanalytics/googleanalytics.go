// Package googleanalytics provides Google Analytics tracking
// for an Ambient application.
package googleanalytics

import (
	"embed"
	"fmt"

	"github.com/josephspurrier/ambient/app/core"
)

//go:embed *
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginMeta
	*core.Toolkit
}

// New returns a new googleanalytics plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "googleanalytics",
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

const (
	trackingID = "Tracking ID"
)

// Fields returns a list of user settable fields.
func (p *Plugin) Fields() []string {
	return []string{
		trackingID,
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]core.Asset, *embed.FS) {
	// Get the tracking ID.
	trackingid, err := p.Site.PluginField(trackingID)
	if err != nil || len(trackingid) == 0 {
		// Otherwise don't set the assets.
		return nil, nil
	}

	return []core.Asset{
		{
			Path:     fmt.Sprintf("https://www.googletagmanager.com/gtag/js?id=%v", trackingid),
			Filetype: core.AssetJavaScript,
			Location: core.LocationBody,
			Embedded: false,
			Auth:     core.AnonymousOnly,
			Attributes: []core.Attribute{
				{
					Name:  "async",
					Value: nil,
				},
			},
		},
		{
			Path:     "js/ga.js",
			Filetype: core.AssetJavaScript,
			Location: core.LocationBody,
			Embedded: true,
			Auth:     core.AnonymousOnly,
			Replace: []core.Replace{
				{
					Find:    "{{TrackingID}}",
					Replace: trackingid,
				},
			},
		},
	}, &assets
}
