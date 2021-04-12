// Package viewport provides viewport functionality
// for an Ambient application.
package viewport

import (
	"embed"

	"github.com/josephspurrier/ambient/app/core"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginBase
	*core.Toolkit
}

// New returns a new viewport plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &core.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "viewport"
}

// PluginVersion returns the plugin version.
func (p *Plugin) PluginVersion() string {
	return "1.0.0"
}

const (
	// Viewport allows user to set the viewport.
	Viewport = "Viewport"
)

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *core.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

// Settings returns a list of user settable fields.
func (p *Plugin) Settings() []core.Setting {
	return []core.Setting{
		{
			Name:    Viewport,
			Default: "width=device-width, initial-scale=1.0, maximum-scale=1.0,user-scalable=0",
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]core.Asset, *embed.FS) {
	vp, err := p.Site.PluginSettingString(Viewport)
	if err != nil || len(vp) == 0 {
		// Otherwise don't set the assets.
		return nil, nil
	}

	return []core.Asset{
		{
			Filetype:   core.AssetGeneric,
			Location:   core.LocationHead,
			TagName:    "meta",
			ClosingTag: false,
			Attributes: []core.Attribute{
				{
					Name:  "name",
					Value: "viewport",
				},
				{
					Name:  "content",
					Value: vp,
				},
			},
		},
	}, nil
}
