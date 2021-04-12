// Package charset provides charset functionality
// for an Ambient application.
package charset

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

// New returns a new charset plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "charset",
			Version:    "1.0.0",
			AppVersion: "1.0.0",
		},
	}
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

// Fields returns a list of user settable fields.
func (p *Plugin) Fields() []core.Setting {
	return []core.Setting{
		{
			Name:    Charset,
			Default: "utf-8",
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]core.Asset, *embed.FS, func(r *http.Request) template.FuncMap) {
	cs, err := p.Site.PluginSettingString(Charset)
	if err != nil || len(cs) == 0 {
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
					Name:  "charset",
					Value: cs,
				},
			},
		},
	}, nil, nil
}
