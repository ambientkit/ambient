// Package googleanalytics provides Google Analytics tracking
// for an Ambient application.
package googleanalytics

import (
	"embed"
	"fmt"

	"github.com/josephspurrier/ambient/app/core"
)

//go:embed js/*.js
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginBase
	*core.Toolkit
}

// New returns a new googleanalytics plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &core.PluginBase{},
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
func (p *Plugin) Enable(toolkit *core.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

// Grants returns a list of grants requested by the plugin.
func (p *Plugin) Grants() []core.Grant {
	return []core.Grant{
		core.GrantPluginSettingRead,
	}
}

const (
	// TrackingID allows the user to set the Google Analytics property ID.
	TrackingID = "Tracking ID"
)

// Settings returns a list of user settable fields.
func (p *Plugin) Settings() []core.Setting {
	return []core.Setting{
		{
			Name: TrackingID,
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]core.Asset, *embed.FS) {
	// Get the tracking ID.
	trackingID, err := p.Site.PluginSettingString(TrackingID)
	if err != nil || len(trackingID) == 0 {
		// Otherwise don't set the assets.
		return nil, nil
	}

	return []core.Asset{
		{
			Path:     fmt.Sprintf("https://www.googletagmanager.com/gtag/js?id=%v", trackingID),
			Filetype: core.AssetJavaScript,
			Location: core.LocationBody,
			External: true,
			Auth:     core.AuthAnonymousOnly,
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
			Auth:     core.AuthAnonymousOnly,
			Replace: []core.Replace{
				{
					Find:    "{{TrackingID}}",
					Replace: trackingID,
				},
			},
		},
	}, &assets
}
