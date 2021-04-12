// Package styles provides a styles page
// for an Ambient application.
package styles

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

// New returns a new styles plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "styles",
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
	// Favicon allows user to set the favicon.
	Favicon = "Favicon"
	// Styles allows user to set the styles.
	Styles = "Styles"
)

// Fields returns a list of user settable fields.
func (p *Plugin) Fields() []core.Setting {
	return []core.Setting{
		{
			Name: Favicon,
			Description: core.SettingDescription{
				Text: "Emoji cheatsheet",
				URL:  "https://emojicheatsheet.com/",
			},
		},
		{
			Name: Styles,
			Type: core.Textarea,
			Description: core.SettingDescription{
				Text: "No-class css themes. You can also paste a link like this: @import 'https://cdn.jsdelivr.net/gh/kognise/water.css@latest/dist/dark.css'",
				URL:  "https://www.cssbed.com/",
			},
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]core.Asset, *embed.FS, func(r *http.Request) template.FuncMap) {
	arr := make([]core.Asset, 0)

	favicon, err := p.Site.PluginSettingString(Favicon)
	if err == nil && len(favicon) > 0 {
		arr = append(arr, core.Asset{
			Filetype:   core.AssetGeneric,
			Location:   core.LocationHead,
			TagName:    "link",
			ClosingTag: false,
			Attributes: []core.Attribute{
				{
					Name:  "rel",
					Value: "icon",
				},
				{
					Name:  "href",
					Value: "data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 100 100%22><text y=%22.9em%22 font-size=%2290%22>" + favicon + "</text></svg>",
				},
			},
		})
	}

	s, err := p.Site.PluginSettingString(Styles)
	if err == nil && len(s) > 0 {
		arr = append(arr, core.Asset{
			Path:     "/plugins/styles/css/style.css",
			Filetype: core.AssetStylesheet,
			Location: core.LocationHead,
			External: true,
		})
	}

	return arr, nil, nil
}

// Routes gets routes for the plugin.
func (p *Plugin) Routes() {
	p.Mux.Get("/plugins/styles/css/style.css", p.index)
}
