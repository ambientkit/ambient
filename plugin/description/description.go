// Package description provides description functionality
// for an Ambient application.
package description

import (
	"embed"
	"fmt"

	"github.com/josephspurrier/ambient/app/core"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginBase
	*core.Toolkit
}

// New returns a new description plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &core.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "description"
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
	// Description allows user to set the description.
	Description = "Description"
)

// Settings returns a list of user settable fields.
func (p *Plugin) Settings() []core.Setting {
	return []core.Setting{
		{
			Name: Description,
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]core.Asset, *embed.FS) {
	siteDescription, err := p.Site.PluginSettingString(Description)
	if err != nil || len(siteDescription) == 0 {
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
					Value: "description",
				},
				{
					Name:  "content",
					Value: fmt.Sprintf("{{if .pagedescription}}{{.pagedescription}}{{else}}%v{{end}}", siteDescription),
				},
			},
		},
	}, nil
}
