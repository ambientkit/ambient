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

const (
	// Description allows user to set the description.
	Description = "Description"
)

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *core.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

// Fields returns a list of user settable fields.
func (p *Plugin) Fields() []core.Field {
	return []core.Field{
		{
			Name: Description,
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]core.Asset, *embed.FS, func(r *http.Request) template.FuncMap) {
	siteDescription, err := p.Site.PluginField(Description)
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
