// Package disqus provides Disqus commenting
// for an Ambient application.
package disqus

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/josephspurrier/ambient/app/core"
)

//go:embed css/*.css js/*.js
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginBase
	*core.Toolkit
}

// New returns a new disqus plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &core.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "disqus"
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

const (
	// DisqusID allows user to set the Disqus ID.
	DisqusID = "Disqus ID"
)

// Settings returns a list of user settable fields.
func (p *Plugin) Settings() []core.Setting {
	return []core.Setting{
		{
			Name: DisqusID,
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]core.Asset, *embed.FS) {
	// Get the Disqus ID.
	disqusID, err := p.Site.PluginSettingString(DisqusID)
	if err != nil || len(disqusID) == 0 {
		// Otherwise don't set the assets.
		return nil, nil
	}

	siteURL, err := p.Site.FullURL()
	if err != nil || len(siteURL) == 0 {
		// Otherwise don't set the assets.
		return nil, nil
	}

	return []core.Asset{
		{
			Path:     "css/disqus.css",
			Filetype: core.AssetStylesheet,
			Location: core.LocationHead,
			LayoutOnly: []core.LayoutType{
				core.LayoutPost,
			},
		},
		{
			Path:     "js/disqus.js",
			Filetype: core.AssetJavaScript,
			Location: core.LocationBody,
			LayoutOnly: []core.LayoutType{
				core.LayoutPost,
			},
			Inline: true,
			Replace: []core.Replace{
				{
					Find:    "{{.DisqusID}}",
					Replace: disqusID,
				},
				{
					Find:    "{{.SiteURL}}",
					Replace: siteURL,
				},
			},
		},
		{
			Filetype: core.AssetGeneric,
			Location: core.LocationMain,
			LayoutOnly: []core.LayoutType{
				core.LayoutPost,
			},
			TagName:    "div",
			ClosingTag: true,
			Attributes: []core.Attribute{
				{
					Name:  "id",
					Value: "disqus_thread",
				},
			},
		},
	}, &assets
}

// FuncMap returns a callable function when passed in a request.
func (p *Plugin) FuncMap() func(r *http.Request) template.FuncMap {
	return p.funcMap
}

func (p *Plugin) funcMap(r *http.Request) template.FuncMap {
	fm := make(template.FuncMap)
	fm["PageURL"] = func() string {
		return r.URL.Path
	}

	return fm
}
