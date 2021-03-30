// Package viewport provides viewport functionality
// for an Ambient application.
package viewport

import (
	"embed"

	"github.com/josephspurrier/ambient/app/core"
)

//go:embed *
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginMeta
	*core.Toolkit
}

// New returns a new viewport plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "viewport",
			Version:    "1.0.0",
			AppVersion: "1.0.0",
		},
	}
}

const (
	// Viewport allows user to set the viewport.
	Viewport = "Viewport"
)

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *core.Toolkit) error {
	p.Toolkit = toolkit

	// Set the default on enable.
	cs, _ := p.Site.PluginField(Viewport)
	if len(cs) == 0 {
		err := p.Site.SetPluginField(Viewport, "width=device-width, initial-scale=1.0, maximum-scale=1.0,user-scalable=0")
		if err != nil {
			return err
		}
	}

	return nil
}

// Fields returns a list of user settable fields.
func (p *Plugin) Fields() []core.Field {
	return []core.Field{
		{
			Name: Viewport,
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]core.Asset, *embed.FS) {
	vp, err := p.Site.PluginField(Viewport)
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
	}, &assets
}
