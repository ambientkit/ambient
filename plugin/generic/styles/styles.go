// Package styles is an Ambient plugin that provides a page to edit styles.
package styles

import (
	"embed"

	"github.com/ambientkit/ambient"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit
}

// New returns an Ambient plugin that provides a page to edit styles.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "styles"
}

// PluginVersion returns the plugin version.
func (p *Plugin) PluginVersion() string {
	return "1.0.0"
}

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *ambient.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

// GrantRequests returns a list of grants requested by the plugin.
func (p *Plugin) GrantRequests() []ambient.GrantRequest {
	return []ambient.GrantRequest{
		{Grant: ambient.GrantPluginSettingRead, Description: "Access to read the plugin settings."},
		{Grant: ambient.GrantSiteAssetWrite, Description: "Access to add favicon."},
		{Grant: ambient.GrantRouterRouteWrite, Description: "Access to add a route for styles."},
	}
}

const (
	// Favicon allows user to set the favicon.
	Favicon = "Favicon"
	// Styles allows user to set the styles.
	Styles = "Styles"
)

// Settings returns a list of user settable fields.
func (p *Plugin) Settings() []ambient.Setting {
	return []ambient.Setting{
		{
			Name: Favicon,
			Description: ambient.SettingDescription{
				Text: "Emoji cheatsheet",
				URL:  "https://github.com/ikatyang/emoji-cheat-sheet/blob/master/README.md",
			},
		},
		{
			Name: Styles,
			Type: ambient.Textarea,
			Description: ambient.SettingDescription{
				Text: "No-class css themes. You can also paste a link like this: @import 'https://cdn.jsdelivr.net/gh/kognise/water.css@latest/dist/dark.css'",
				URL:  "https://www.cssbed.com/",
			},
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]ambient.Asset, *embed.FS) {
	arr := make([]ambient.Asset, 0)

	favicon, err := p.Site.PluginSettingString(Favicon)
	if err == nil && len(favicon) > 0 {
		arr = append(arr, ambient.Asset{
			Filetype:   ambient.AssetGeneric,
			Location:   ambient.LocationHead,
			TagName:    "link",
			ClosingTag: false,
			Attributes: []ambient.Attribute{
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
		arr = append(arr, ambient.Asset{
			Path:     "css/style.css",
			Filetype: ambient.AssetStylesheet,
			Location: ambient.LocationHead,
		})
	}

	return arr, nil
}

// Routes gets routes for the plugin.
func (p *Plugin) Routes() {
	p.Mux.Get("/plugins/styles/css/style.css", p.index)
}
