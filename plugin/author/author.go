// Package author provides author functionality
// for an Ambient application.
package author

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

// New returns a new author plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "author",
			Version:    "1.0.0",
			AppVersion: "1.0.0",
		},
	}
}

const (
	// Author allows user to set the author.
	Author = "Author"
)

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *core.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

// Grants returns a list of grants requested by the plugin.
func (p *Plugin) Grants() []core.Grant {
	return []core.Grant{
		core.GrantPluginFieldRead,
		core.GrantPluginFieldWrite,
	}
}

// Settings returns a list of user settable fields.
func (p *Plugin) Settings() []core.Setting {
	return []core.Setting{
		{
			Name: Author,
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]core.Asset, *embed.FS, func(r *http.Request) template.FuncMap) {
	name, err := p.Site.PluginSettingString(Author)
	if err != nil || len(name) == 0 {
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
					Value: "author",
				},
				{
					Name:  "content",
					Value: name,
				},
			},
		},
	}, nil, nil
}
