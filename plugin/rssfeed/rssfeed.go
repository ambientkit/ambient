// Package rssfeed provides rss feed functionality
// for an Ambient application.
package rssfeed

import (
	"embed"

	"github.com/josephspurrier/ambient/app/core"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginBase
	*core.Toolkit
}

// New returns a new rssfeed plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &core.PluginBase{},
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
func (p *Plugin) Enable(toolkit *core.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

// Grants returns a list of grants requested by the plugin.
func (p *Plugin) Grants() []core.Grant {
	return []core.Grant{
		core.GrantSiteTitleRead,
		core.GrantSiteSchemeRead,
		core.GrantSiteURLRead,
		core.GrantSitePostRead,
		core.GrantPluginSettingRead,
		core.GrantPluginSettingWrite,
	}
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
func (p *Plugin) Assets() ([]core.Asset, *embed.FS) {
	siteTitle, err := p.Site.Title()
	if err != nil {
		return nil, nil
	}

	feedURL, err := p.Site.PluginSettingString(FeedURL)
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
