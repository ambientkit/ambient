// Package googleanalytics is an Ambient plugin that provides Google Analytics tracking.
package googleanalytics

import (
	"embed"
	"fmt"

	"github.com/ambientkit/ambient"
)

//go:embed js/*.js
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit
}

// New returns an Ambient plugin that provides Google Analytics tracking.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "googleanalytics"
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

// GrantRequests returns a list of grants requested by the plugin.
func (p *Plugin) GrantRequests() []ambient.GrantRequest {
	return []ambient.GrantRequest{
		{Grant: ambient.GrantPluginSettingRead, Description: "Access to the tracking ID."},
	}
}

const (
	// TrackingID allows the user to set the Google Analytics property ID.
	TrackingID = "Tracking ID"
)

// Settings returns a list of user settable fields.
func (p *Plugin) Settings() []ambient.Setting {
	return []ambient.Setting{
		{
			Name: TrackingID,
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]ambient.Asset, *embed.FS) {
	// Get the tracking ID.
	trackingID, err := p.Site.PluginSettingString(TrackingID)
	if err != nil || len(trackingID) == 0 {
		// Otherwise don't set the assets.
		return nil, nil
	}

	return []ambient.Asset{
		{
			Path:     fmt.Sprintf("https://www.googletagmanager.com/gtag/js?id=%v", trackingID),
			Filetype: ambient.AssetJavaScript,
			Location: ambient.LocationBody,
			External: true,
			Auth:     ambient.AuthAnonymousOnly,
			Attributes: []ambient.Attribute{
				{
					Name:  "async",
					Value: nil,
				},
			},
		},
		{
			Path:     "js/ga.js",
			Filetype: ambient.AssetJavaScript,
			Location: ambient.LocationBody,
			Auth:     ambient.AuthAnonymousOnly,
			Replace: []ambient.Replace{
				{
					Find:    "{{TrackingID}}",
					Replace: trackingID,
				},
			},
		},
	}, &assets
}
