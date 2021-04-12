// Package charset provides charset functionality
// for an Ambient application.
package charset

import (
	"embed"

	"github.com/josephspurrier/ambient/app/core"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginBase
	*core.Toolkit
}

// New returns a new charset plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &core.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "charset"
}

// PluginVersion returns the plugin version.
func (p *Plugin) PluginVersion() string {
	return "1.0.0"
}

const (
	// Charset allows user to set the charset.
	Charset = "Charset"
)

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

// Settings returns a list of user settable fields.
func (p *Plugin) Settings() []core.Setting {
	return []core.Setting{
		{
			Name:    Charset,
			Default: "utf-8",
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]core.Asset, *embed.FS) {
	cs, err := p.Site.PluginSettingString(Charset)
	if err != nil || len(cs) == 0 {
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
					Name:  "charset",
					Value: cs,
				},
			},
		},
	}, nil
}
