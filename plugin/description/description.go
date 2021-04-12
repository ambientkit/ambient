// Package description provides description functionality
// for an Ambient application.
package description

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/josephspurrier/ambient/app/core"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginMeta
	*core.Toolkit
}

// New returns a new description plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "description",
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
func (p *Plugin) Assets() ([]core.Asset, *embed.FS, func(r *http.Request) template.FuncMap) {
	siteDescription, err := p.Site.PluginSettingString(Description)
	if err != nil || len(siteDescription) == 0 {
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
					Value: "description",
				},
				{
					Name:  "content",
					Value: fmt.Sprintf("{{if .pagedescription}}{{.pagedescription}}{{else}}%v{{end}}", siteDescription),
				},
			},
		},
	}, nil, nil
}
