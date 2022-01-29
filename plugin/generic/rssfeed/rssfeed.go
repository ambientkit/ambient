// Package rssfeed is an Ambient plugin that provides an RSS feed.
package rssfeed

import (
	"embed"

	"github.com/ambientkit/ambient"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit
}

// New returns an Ambient plugin that provides an RSS feed.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "rssfeed"
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
		{Grant: ambient.GrantSiteTitleRead, Description: "Access to read the site title."},
		{Grant: ambient.GrantSiteSchemeRead, Description: "Access to read the site scheme."},
		{Grant: ambient.GrantSiteURLRead, Description: "Access to read the site URL."},
		{Grant: ambient.GrantSitePostRead, Description: "Access to read all the site posts."},
		{Grant: ambient.GrantPluginSettingRead, Description: "Access to read the plugin settings."},
		{Grant: ambient.GrantPluginSettingWrite, Description: "Access to write to the plugin settings."},
	}
}

const (
	// FeedURL allows user to set the feed URL>
	FeedURL = "Feed URL"
	// Description allows user to set the description.
	Description = "Description"
)

// Settings returns a list of user settable fields.
func (p *Plugin) Settings() []ambient.Setting {
	return []ambient.Setting{
		{
			Name:    FeedURL,
			Default: "/rss.xml",
			Description: ambient.SettingDescription{
				Text: "Must start with a slash like this: /rss.xml",
			},
		},
		{
			Name: Description,
			Type: ambient.Textarea,
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]ambient.Asset, *embed.FS) {
	siteTitle, err := p.Site.Title()
	if err != nil {
		return nil, nil
	}

	feedURL, err := p.Site.PluginSettingString(FeedURL)
	if err != nil {
		return nil, nil
	}

	return []ambient.Asset{
		{
			Filetype:   ambient.AssetGeneric,
			Location:   ambient.LocationHead,
			TagName:    "link",
			ClosingTag: false,
			Attributes: []ambient.Attribute{
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
	}, nil
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
