// Package googleanalytics provides Google Analytics tracking
// for an Ambient application.
package googleanalytics

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/josephspurrier/ambient/app/core"
)

//go:embed js/*.js
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
	// TrackingID allows the user to set the Google Analytics property ID.
	TrackingID = "Tracking ID"
)

// Fields returns a list of user settable fields.
func (p *Plugin) Fields() []core.Setting {
	return []core.Setting{
		{
			Name: TrackingID,
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]core.Asset, *embed.FS, func(r *http.Request) template.FuncMap) {
	// Get the tracking ID.
	trackingID, err := p.Site.PluginSettingString(TrackingID)
	if err != nil || len(trackingID) == 0 {
		// Otherwise don't set the assets.
		return nil, nil, nil
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
	}, &assets, nil
}
