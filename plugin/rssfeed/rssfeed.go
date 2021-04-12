// Package rssfeed provides rss feed functionality
// for an Ambient application.
package rssfeed

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
	// FeedURL allows user to set the feed URL>
	FeedURL = "Feed URL"
	// Description allows user to set the description.
	Description = "Description"
)

// Settings returns a list of user settable fields.
func (p *Plugin) Settings() []core.Setting {
	return []core.Setting{
		{
			Name:    FeedURL,
			Default: "/rss.xml",
			Description: core.SettingDescription{
				Text: "Must start with a slash like this: /rss.xml",
			},
		},
		{
			Name: Description,
			Type: core.Textarea,
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]core.Asset, *embed.FS, func(r *http.Request) template.FuncMap) {
	siteTitle, err := p.Site.Title()
	if err != nil {
		return nil, nil, nil
	}

	feedURL, err := p.Site.PluginSettingString(FeedURL)
	if err != nil {
		return nil, nil, nil
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
					Value: feedURL,
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
	}, nil, nil
}

// Routes gets routes for the plugin.
func (p *Plugin) Routes() {
	// FIXME: This can't be changed dynamically.
	feedURL, err := p.Site.PluginSettingString(FeedURL)
	if err != nil {
		return
	}

	p.Mux.Get(feedURL, p.index)
}
