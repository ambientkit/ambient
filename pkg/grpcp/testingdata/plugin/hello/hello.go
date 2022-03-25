// Package hello provides a hello page for an Ambient app.
package hello

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/ambientkit/ambient"
)

//go:embed template/*.tmpl template/content/*.tmpl
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
}

// New returns a new hello plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "hello"
}

// PluginVersion returns the plugin version.
func (p *Plugin) PluginVersion() string {
	return "1.0.0"
}

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *ambient.Toolkit) error {
	err := p.PluginBase.Enable(toolkit)
	if err != nil {
		return err
	}

	//p.Log.Info("plugin: enabled called")

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
		{Grant: ambient.GrantRouterMiddlewareWrite, Description: "Access to add middleware."},
	}
}

// Routes sets routes for the plugin.
func (p *Plugin) Routes() {
	//p.Log.Error("plugin: routes called")
	p.Mux.Get("/", p.index)
	p.Mux.Get("/another", p.another)
	p.Mux.Get("/name/{name}", p.name)
	p.Mux.Get("/nameold/{name}", p.Mux.Wrap(p.nameOld))
	p.Mux.Get("/error", p.errorFunc)
	p.Mux.Get("/created", p.created)
	p.Mux.Get("/headers", p.headers)
	p.Mux.Get("/form", p.formGet)
	p.Mux.Post("/form", p.formPOST)
	p.Mux.Get("/login", p.login)
	p.Mux.Get("/loggedin", p.loggedin)
	p.Mux.Get("/errors", p.errorsFunc)
	p.Mux.Get("/neighborPluginGrantList", p.neighborPluginGrantList)
	p.Mux.Get("/neighborPluginGrantListBad", p.neighborPluginGrantListBad)
	p.Mux.Get("/neighborPluginGrants", p.neighborPluginGrants)
	p.Mux.Get("/neighborPluginGranted", p.neighborPluginGranted)
	p.Mux.Get("/neighborPluginGrantedBad", p.neighborPluginGrantedBad)
	p.Mux.Get("/setNeighborPluginGrantFalse", p.setNeighborPluginGrantFalse)
	p.Mux.Get("/setNeighborPluginGrantTrue", p.setNeighborPluginGrantTrue)
	p.Mux.Get("/neighborPluginRequestedGrant", p.neighborPluginRequestedGrant)
	p.Mux.Get("/neighborPluginRequestedGrantBad", p.neighborPluginRequestedGrantBad)
	p.Mux.Get("/plugins", p.plugins)
	p.Mux.Get("/pluginNames", p.pluginNames)
	p.Mux.Delete("/deletePlugin", p.deletePlugin)
	p.Mux.Delete("/deletePluginBad", p.deletePluginBad)
	p.Mux.Get("/enablePlugin", p.enablePlugin)
	p.Mux.Get("/enablePluginBad", p.enablePluginBad)
	p.Mux.Get("/disablePlugin", p.disablePlugin)
	p.Mux.Get("/disablePluginBad", p.disablePluginBad)
	p.Mux.Post("/savePost", p.savePost)
	p.Mux.Get("/publishedPosts", p.publishedPosts)
	p.Mux.Get("/publishedPages", p.publishedPages)
	p.Mux.Get("/postBySlug", p.postBySlug)
	p.Mux.Get("/postBySlugBad", p.postBySlugBad)
	p.Mux.Get("/postByID", p.postByID)
	p.Mux.Get("/postByIDBad", p.postByID)
	p.Mux.Delete("/deletePostByID", p.deletePostByID)
	p.Mux.Get("/pluginNeighborRoutesList", p.pluginNeighborRoutesList)
	p.Mux.Get("/pluginNeighborRoutesListBad", p.pluginNeighborRoutesListBad)
	p.Mux.Get("/userPersist", p.userPersist)
	p.Mux.Get("/userPersistFalse", p.userPersistFalse)
	p.Mux.Get("/grantRequests", p.grantRequests)
	p.Mux.Get("/userLogout", p.userLogout)
	p.Mux.Get("/logoutAllUsers", p.logoutAllUsers)
	p.Mux.Get("/csrf", p.setCSRF)
	p.Mux.Post("/csrf", p.cSRF)
	p.Mux.Get("/sessionValue", p.sessionValue)
	p.Mux.Get("/PluginNeighborSettingsList", p.pluginNeighborSettingsList)
	p.Mux.Get("/setPluginSetting", p.setPluginSetting)
	p.Mux.Get("/setNeighborPluginSetting", p.setNeighborPluginSetting)
	p.Mux.Get("/pluginTrusted", p.pluginTrusted)
	p.Mux.Get("/title", p.title)
	p.Mux.Get("/scheme", p.scheme)
	p.Mux.Get("/url", p.uRL)
	p.Mux.Get("/updated", p.updated)
	p.Mux.Get("/content", p.content)
	p.Mux.Get("/tags", p.tags)
	p.Mux.Get("/assets", p.assets)
	p.Mux.Get("/assetsHello", p.assetsHello)
	p.Mux.Get("/assetsError", p.assetsError)
	p.Mux.Get("/pageHello", p.pageHello)
	p.Mux.Get("/context", p.context)
	p.Mux.Get("/redirect", p.redirect)
}

const (
	// Username allows user to set the login username.
	Username = "Username"
	// SafeMode is a boolean value.
	SafeMode = "Safe Mode"
)

// Settings returns a list of user settable fields.
func (p *Plugin) Settings() []ambient.Setting {
	return []ambient.Setting{
		{
			Name:    Username,
			Default: "admin",
		},
		{
			Name:    SafeMode,
			Default: true,
		},
	}
}

// Assets returns a list of assets and an embedded filesystem.
func (p *Plugin) Assets() ([]ambient.Asset, ambient.FileSystemReader) {
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

	arr = append(arr, ambient.Asset{
		Filetype:   ambient.AssetGeneric,
		Location:   ambient.LocationHead,
		TagName:    "link",
		ClosingTag: false,
		Attributes: []ambient.Attribute{
			{
				Name:  "rel",
				Value: "canonical",
			},
			{
				Name:  "href",
				Value: `{{if .canonical}}{{.canonical}}{{else}}{{hello_Cool}}{{end}}`,
			},
		},
	})

	// arr = append(arr, ambient.Asset{
	// 	Path:     "template/partial/nav.tmpl",
	// 	Filetype: ambient.AssetGeneric,
	// 	Location: ambient.LocationHeader,
	// 	Inline:   true,
	// })

	// arr = append(arr, ambient.Asset{
	// 	Path:     "template/partial/footer.tmpl",
	// 	Filetype: ambient.AssetGeneric,
	// 	Location: ambient.LocationFooter,
	// 	Inline:   true,
	// })

	//return arr, &assets
	return arr, &assets
}

// FuncMap returns a callable function that accepts a request.
func (p *Plugin) FuncMap() func(r *http.Request) template.FuncMap {
	return func(r *http.Request) template.FuncMap {
		fm := make(template.FuncMap)
		fm["hello_Cool"] = func() string {
			return "cool"
		}
		fm["hello_Foo"] = func(name string) string {
			return "hello: " + name
		}
		fm["hello_Error"] = func(name string) error {
			return errors.New("this is an error")
		}

		return fm
	}
}

// Middleware returns router middleware.
func (p *Plugin) Middleware() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{
		func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Set a context key.
				h.ServeHTTP(w, Set(r, "foo"))
			})
		},
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/api/healthcheck" {
					m := make(map[string]interface{})
					m["message"] = "ok"
					b, _ := json.Marshal(m)
					w.WriteHeader(http.StatusCreated)
					w.Header().Add("Content-Type", "application/json")
					fmt.Fprint(w, string(b))
					return
				}
				next.ServeHTTP(w, r)
			})
		},
	}
}
