// Package viewport provides viewport functionality
// for an Ambient application.
package viewport

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/josephspurrier/ambient/app/core"
)

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
	return nil
}

// Fields returns a list of user settable fields.
func (p *Plugin) Fields() []core.Setting {
	return []core.Setting{
		{
			Name:    Viewport,
			Default: "width=device-width, initial-scale=1.0, maximum-scale=1.0,user-scalable=0",
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]core.Asset, *embed.FS, func(r *http.Request) template.FuncMap) {
	vp, err := p.Site.PluginSettingString(Viewport)
	if err != nil || len(vp) == 0 {
		// Otherwise don't set the assets.
		return nil, nil, nil
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
	}, nil, nil
}
