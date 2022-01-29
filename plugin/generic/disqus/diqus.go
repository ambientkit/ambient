// Package disqus is an Ambient plugin that provides Disqus commenting.
package disqus

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/ambientkit/ambient"
)

//go:embed css/*.css js/*.js
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit
}

// New returns an Ambient plugin that provides Disqus commenting.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
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
func (p *Plugin) Enable(toolkit *ambient.Toolkit) error {
	p.Toolkit = toolkit
	return nil
}

// GrantRequests returns a list of grants requested by the plugin.
func (p *Plugin) GrantRequests() []ambient.GrantRequest {
	return []ambient.GrantRequest{
		{Grant: ambient.GrantPluginSettingRead, Description: "Access to the Disqus ID."},
		{Grant: ambient.GrantSiteURLRead, Description: "Access to read the site URL."},
		{Grant: ambient.GrantSiteSchemeRead, Description: "Access to read the site scheme."},
		{Grant: ambient.GrantSiteFuncMapWrite, Description: "Access to create global FuncMaps for templates."},
		{Grant: ambient.GrantSiteAssetWrite, Description: "Access to write meta tags to the header and add a nav and footer."},
		{Grant: ambient.GrantRouterRouteWrite, Description: "Access to create routes for serving javascript and stylesheets."},
	}
}

const (
	// DisqusID allows user to set the Disqus ID.
	DisqusID = "Disqus ID"
)

// Settings returns a list of user settable fields.
func (p *Plugin) Settings() []ambient.Setting {
	return []ambient.Setting{
		{
			Name: DisqusID,
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]ambient.Asset, *embed.FS) {
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

	return []ambient.Asset{
		{
			Path:     "css/disqus.css",
			Filetype: ambient.AssetStylesheet,
			Location: ambient.LocationHead,
			LayoutOnly: []ambient.LayoutType{
				ambient.LayoutPost,
			},
		},
		{
			Path:     "js/disqus.js",
			Filetype: ambient.AssetJavaScript,
			Location: ambient.LocationBody,
			LayoutOnly: []ambient.LayoutType{
				ambient.LayoutPost,
			},
			Inline: true,
			Replace: []ambient.Replace{
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
			Filetype: ambient.AssetGeneric,
			Location: ambient.LocationMain,
			LayoutOnly: []ambient.LayoutType{
				ambient.LayoutPost,
			},
			TagName:    "div",
			ClosingTag: true,
			Attributes: []ambient.Attribute{
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
	return func(r *http.Request) template.FuncMap {
		fm := make(template.FuncMap)
		fm["disqus_PageURL"] = func() string {
			return r.URL.Path
		}

		return fm
	}
}
