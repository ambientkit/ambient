// Package disqus provides Disqus commenting
// for an Ambient application.
package disqus

import (
	"embed"

	"github.com/josephspurrier/ambient/app/core"
)

//go:embed *
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginMeta
	*core.Toolkit
}

// New returns a new disqus plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "disqus",
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
	// DisqusID allows user to set the Disqus ID.
	DisqusID = "Disqus ID"
)

// Fields returns a list of user settable fields.
func (p *Plugin) Fields() []string {
	return []string{
		DisqusID,
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]core.Asset, *embed.FS) {
	// Get the Disqus ID.
	disqusID, err := p.Site.PluginField(DisqusID)
	if err != nil || len(disqusID) == 0 {
		// Otherwise don't set the assets.
		return nil, nil
	}

	URL, err := p.Site.URL()
	if err != nil || len(disqusID) == 0 {
		// Otherwise don't set the assets.
		return nil, nil
	}

	return []core.Asset{
		{
			Path:     "css/disqus.css",
			Filetype: core.AssetStylesheet,
			Location: core.LocationHead,
			Embedded: true,
		},
		{
			Path:     "js/disqus.js",
			Filetype: core.AssetJavaScript,
			Location: core.LocationBody,
			Embedded: true,
			Replace: []core.Replace{
				{
					Find:    "{{DisqusID}}",
					Replace: disqusID,
				},
				{
					Find:    "{{SiteURL}}",
					Replace: URL,
				},
				{
					Find:    "{{.posturl}}",
					Replace: "789",
				},
				{
					Find:    "{{.id}}",
					Replace: "101",
				},
			},
		},
		{
			Filetype:   core.AssetGeneric,
			Location:   core.LocationMain,
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
