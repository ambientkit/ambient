// Package prism is an Ambient plugin that provides syntax highlighting using Prism (https://prismjs.com/).
package prism

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/ambientkit/ambient"
)

//go:embed css/*.css
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit
}

// New returns an Ambient plugin that provides syntax highlighting using Prism (https://prismjs.com/).
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "prism"
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
		{Grant: ambient.GrantSiteAssetWrite, Description: "Access to add stylesheets and javascript to each page."},
		{Grant: ambient.GrantRouterRouteWrite, Description: "Access to create routes for accessing stylesheets."},
		{Grant: ambient.GrantPluginSettingRead, Description: "Read own plugin settings."},
	}
}

const (
	// Version allows user to set the library version.
	Version = "Version"
	// Styles allows user to set the styles.
	Styles = "Styles"
)

// Settings returns a list of user settable fields.
func (p *Plugin) Settings() []ambient.Setting {
	return []ambient.Setting{
		{
			Name:    Version,
			Default: "1.23.0",
			Description: ambient.SettingDescription{
				Text: "View releases (ex: 1.23.0)",
				URL:  "https://github.com/PrismJS/prism/releases",
			},
		},
		{
			Name: Styles,
			Type: ambient.Textarea,
			Description: ambient.SettingDescription{
				Text: "You can paste a theme from https://github.com/PrismJS/prism-themes/tree/master/themes or an import like this using https://gitcdn.link/: @import 'https://gitcdn.link/cdn/PrismJS/prism-themes/d00360c3b3cfe495f45cc06865969c7731a94763/themes/prism-vsc-dark-plus.css'",
				URL:  "https://github.com/PrismJS/prism-themes/tree/master/themes",
			},
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]ambient.Asset, *embed.FS) {
	version, err := p.Site.PluginSettingString(Version)
	if err != nil {
		return nil, nil
	}

	arr := []ambient.Asset{
		{
			Path:     "css/clean.css",
			Filetype: ambient.AssetStylesheet,
			Location: ambient.LocationHead,
		},
		{
			Path:     fmt.Sprintf("https://unpkg.com/prismjs@%v/components/prism-core.min.js", version),
			Filetype: ambient.AssetJavaScript,
			Location: ambient.LocationBody,
			External: true,
		},
		{
			Path:     fmt.Sprintf("https://unpkg.com/prismjs@%v/plugins/autoloader/prism-autoloader.min.js", version),
			Filetype: ambient.AssetJavaScript,
			Location: ambient.LocationBody,
			External: true,
		},
	}

	s, err := p.Site.PluginSettingString(Styles)
	if err == nil && len(s) > 0 {
		arr = append(arr, ambient.Asset{
			Path:           "css/style.css",
			Filetype:       ambient.AssetStylesheet,
			Location:       ambient.LocationHead,
			SkipExistCheck: true,
		})
	}

	return arr, &assets
}

// Routes gets routes for the plugin.
func (p *Plugin) Routes() {
	p.Mux.Get(fmt.Sprintf("/plugins/%v/css/style.css", p.PluginName()), p.index)
}

// index returns CSS file.
func (p *Plugin) index(w http.ResponseWriter, r *http.Request) (status int, err error) {
	// Get the styles.
	s, err := p.Site.PluginSetting(Styles)
	if err != nil {
		return p.Site.Error(err)
	}

	w.Header().Set("Content-Type", "text/css")

	fmt.Fprint(w, s)
	return
}
