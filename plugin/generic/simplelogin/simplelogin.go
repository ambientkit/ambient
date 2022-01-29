// Package simplelogin is an Ambient plugin that provides a basic website template with a login page.
package simplelogin

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/ambientkit/ambient"
)

//go:embed template/partial/*.tmpl template/content/*.tmpl
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit

	passwordHash string
}

// New returns an Ambient plugin that provides a basic website template with a login page.
func New(passwordHash string) *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},

		passwordHash: passwordHash,
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "simplelogin"
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
		{Grant: ambient.GrantUserAuthenticatedRead, Description: "Show different menus to authenticated vs unauthenticated users."},
		{Grant: ambient.GrantUserAuthenticatedWrite, Description: "Access to login and logout the user."},
		{Grant: ambient.GrantUserPersistWrite, Description: "Access to set session as persistent."},
		{Grant: ambient.GrantPluginSettingRead, Description: "Read own plugin settings."},
		{Grant: ambient.GrantPluginSettingWrite, Description: "Write own plugin settings."},
		{Grant: ambient.GrantSitePostRead, Description: "Read all site posts."},
		{Grant: ambient.GrantSitePostWrite, Description: "Create and edit site posts."},
		{Grant: ambient.GrantSiteSchemeRead, Description: "Read site scheme."},
		{Grant: ambient.GrantSiteSchemeWrite, Description: "Update the site scheme."},
		{Grant: ambient.GrantSiteURLRead, Description: "Read the site URL."},
		{Grant: ambient.GrantSiteURLWrite, Description: "Update the site URL."},
		{Grant: ambient.GrantSiteTitleRead, Description: "Read the site title."},
		{Grant: ambient.GrantSiteTitleWrite, Description: "Update the site title."},
		{Grant: ambient.GrantSiteContentRead, Description: "Read home page content."},
		{Grant: ambient.GrantSiteContentWrite, Description: "Update home page content."},
		{Grant: ambient.GrantSiteAssetWrite, Description: "Access to write blog meta tags to the header and add a nav and footer."},
		{Grant: ambient.GrantSiteFuncMapWrite, Description: "Access to create global FuncMaps for templates."},
		{Grant: ambient.GrantRouterRouteWrite, Description: "Access to create routes for editing the blog posts."},
	}
}

const (
	// LoginURL allows user to set the login URL.
	LoginURL = "Login URL"
	// Author allows user to set the author.
	Author = "Author"
	// Subtitle allows user to set the Subtitle.
	Subtitle = "Subtitle"
	// Description allows user to set the description.
	Description = "Description"
	// Footer allows user to set the footer.
	Footer = "Footer"
	// AllowHTMLinMarkdown allows user to set if they allow HTML in markdown.
	AllowHTMLinMarkdown = "Allow HTML in Markdown"

	// Username allows user to set the login username.
	Username = "Username"
	// Password allows user to set the login password.
	Password = "Password"
	// MFAKey allows user to set the MFA key.
	MFAKey = "MFA Key"
)

// Settings returns a list of user settable fields.
func (p *Plugin) Settings() []ambient.Setting {
	return []ambient.Setting{
		{
			Name:    Username,
			Default: "admin",
		},
		{
			Name:    Password,
			Default: p.passwordHash,
			Type:    ambient.InputPassword,
			Hide:    true,
		},
		{
			Name: MFAKey,
			Type: ambient.InputPassword,
			Description: ambient.SettingDescription{
				Text: "Generate an MFA key. Plugin must be enabled first.",
				URL:  "/dashboard/mfa",
			},
		},
		{
			Name: Author,
		},
		{
			Name: Subtitle,
			Hide: true,
		},
		{
			Name: Description,
			Type: ambient.Textarea,
		},
		{
			Name: Footer,
			Type: ambient.Textarea,
			Hide: true,
		},
		{
			Name: AllowHTMLinMarkdown,
			Type: ambient.Checkbox,
		},
	}
}

// Routes gets routes for the plugin.
func (p *Plugin) Routes() {
	p.Mux.Get("/", p.index)
	p.Mux.Get("/dashboard", p.edit)
	p.Mux.Post("/dashboard", p.update)
	p.Mux.Get("/dashboard/reload", p.reload)
	p.Mux.Get("/login", p.login)
	p.Mux.Post("/login", p.loginPost)
	p.Mux.Get("/dashboard/mfa", p.mfa)
	p.Mux.Post("/dashboard/mfa", p.mfaPost)
	p.Mux.Get("/dashboard/logout", p.logout)
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]ambient.Asset, *embed.FS) {
	arr := make([]ambient.Asset, 0)

	siteTitle, err := p.Site.Title()
	if err == nil && len(siteTitle) > 0 {
		arr = append(arr, ambient.Asset{
			Filetype: ambient.AssetGeneric,
			Location: ambient.LocationHead,
			TagName:  "title",
			Inline:   true,
			Content:  fmt.Sprintf(`{{if .pagetitle}}{{.pagetitle}} | %v{{else}}%v{{end}}`, siteTitle, siteTitle),
		})
	}

	siteDescription, err := p.Site.PluginSettingString(Description)
	if err == nil && len(siteDescription) > 0 {
		arr = append(arr, ambient.Asset{
			Filetype:   ambient.AssetGeneric,
			Location:   ambient.LocationHead,
			TagName:    "meta",
			ClosingTag: false,
			Attributes: []ambient.Attribute{
				{
					Name:  "name",
					Value: "description",
				},
				{
					Name:  "content",
					Value: fmt.Sprintf("{{if .pagedescription}}{{.pagedescription}}{{else}}%v{{end}}", siteDescription),
				},
			},
		})
	}

	siteAuthor, err := p.Site.PluginSettingString(Author)
	if err == nil && len(siteAuthor) > 0 {
		arr = append(arr, ambient.Asset{
			Filetype:   ambient.AssetGeneric,
			Location:   ambient.LocationHead,
			TagName:    "meta",
			ClosingTag: false,
			Attributes: []ambient.Attribute{
				{
					Name:  "name",
					Value: "author",
				},
				{
					Name:  "content",
					Value: siteAuthor,
				},
			},
		})
	}

	arr = append(arr, ambient.Asset{
		Path:     "template/partial/nav.tmpl",
		Filetype: ambient.AssetGeneric,
		Location: ambient.LocationHeader,
		Inline:   true,
	})

	arr = append(arr, ambient.Asset{
		Path:     "template/partial/footer.tmpl",
		Filetype: ambient.AssetGeneric,
		Location: ambient.LocationFooter,
		Inline:   true,
	})

	return arr, &assets
}

// FuncMap returns a callable function when passed in a request.
func (p *Plugin) FuncMap() func(r *http.Request) template.FuncMap {
	return p.funcMap
}
