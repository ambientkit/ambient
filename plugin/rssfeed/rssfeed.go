// Package rssfeed provides rss feed functionality
// for an Ambient application.
package rssfeed

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

// New returns a new rssfeed plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "rssfeed",
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
	rssURL = "/rss.xml"
)

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]core.Asset, *embed.FS) {
	siteTitle, err := p.Site.Title()
	if err != nil {
		return nil, nil
	}

	return []core.Asset{
		{
			Filetype:   core.AssetGeneric,
			Location:   core.LocationHead,
			TagName:    "link",
			ClosingTag: false,
			Attributes: []core.Attribute{
				{
					Name:  "rel",
					Value: "alternative",
				},
				{
					Name:  "href",
					Value: rssURL,
				},
				{
					Name:  "application",
					Value: "rss+xml",
				},
				{
					Name:  "title",
					Value: siteTitle,
				},
			},
		},
	}, &assets
}

// Routes gets routes for the plugin.
func (p *Plugin) Routes() {
	p.Router.Get(rssURL, p.index)
}
