// Package bearblog provides basic blog functionality
// for an Ambient application.
package bearblog

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/josephspurrier/ambient/app/core"
)

//go:embed template/partial/*.tmpl template/content/*.tmpl
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginMeta
	*core.Toolkit
}

// New returns a new bearblog plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "bearblog",
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
)

// Fields returns a list of user settable fields.
func (p *Plugin) Fields() []core.Field {
	return []core.Field{
		{
			Name:    LoginURL,
			Default: "admin",
		},
		{
			Name: Author,
		},
		{
			Name: Subtitle,
		},
		{
			Name: Description,
			Type: core.Textarea,
		},
		{
			Name: Footer,
			Type: core.Textarea,
		},
		{
			Name: AllowHTMLinMarkdown,
			Type: core.Checkbox,
		},
	}
}

// Routes gets routes for the plugin.
func (p *Plugin) Routes() {
	p.Mux.Get("/blog", p.postIndex)
	p.Mux.Get("/:slug", p.postShow)

	p.Mux.Get("/login/:slug", p.login)
	p.Mux.Post("/login/:slug", p.loginPost)
	p.Mux.Get("/dashboard/logout", p.logout)

	p.Mux.Get("/", p.index)
	p.Mux.Get("/dashboard", p.edit)
	p.Mux.Post("/dashboard", p.update)
	p.Mux.Get("/dashboard/reload", p.reload)

	p.Mux.Get("/dashboard/posts", p.postAdminIndex)
	p.Mux.Get("/dashboard/posts/new", p.postAdminCreate)
	p.Mux.Post("/dashboard/posts/new", p.postAdminStore)
	p.Mux.Get("/dashboard/posts/:id", p.postAdminEdit)
	p.Mux.Post("/dashboard/posts/:id", p.postAdminUpdate)
	p.Mux.Get("/dashboard/posts/:id/delete", p.postAdminDestroy)
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]core.Asset, *embed.FS, func(r *http.Request) template.FuncMap) {
	arr := make([]core.Asset, 0)

	siteTitle, err := p.Site.Title()
	if err == nil && len(siteTitle) > 0 {
		arr = append(arr, core.Asset{
			Filetype: core.AssetGeneric,
			Location: core.LocationHead,
			TagName:  "title",
			Inline:   true,
			Content:  fmt.Sprintf(`{{if .pagetitle}}{{.pagetitle}} | %v{{else}}%v{{end}}`, siteTitle, siteTitle),
		})
	}

	siteDescription, err := p.Site.PluginFieldString(Description)
	if err == nil && len(siteDescription) > 0 {
		arr = append(arr, core.Asset{
			Filetype:   core.AssetGeneric,
			Location:   core.LocationHead,
			TagName:    "meta",
			ClosingTag: false,
			Attributes: []core.Attribute{
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

	arr = append(arr, core.Asset{
		Filetype:   core.AssetGeneric,
		Location:   core.LocationHead,
		TagName:    "link",
		ClosingTag: false,
		Attributes: []core.Attribute{
			{
				Name:  "rel",
				Value: "canonical",
			},
			{
				Name:  "href",
				Value: `{{if .canonical}}{{.canonical}}{{else}}{{PageURL}}{{end}}`,
			},
		},
	})

	siteAuthor, err := p.Site.PluginFieldString(Author)
	if err == nil && len(siteAuthor) > 0 {
		arr = append(arr, core.Asset{
			Filetype:   core.AssetGeneric,
			Location:   core.LocationHead,
			TagName:    "meta",
			ClosingTag: false,
			Attributes: []core.Attribute{
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

	arr = append(arr, core.Asset{
		Path:     "template/partial/nav.tmpl",
		Filetype: core.AssetGeneric,
		Location: core.LocationHeader,
		Inline:   true,
	})

	arr = append(arr, core.Asset{
		Path:     "template/partial/footer.tmpl",
		Filetype: core.AssetGeneric,
		Location: core.LocationFooter,
		Inline:   true,
	})

	return arr, &assets, p.FuncMap
}
