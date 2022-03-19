// Package hello provides a hello page for an Ambient app.
package hello

import (
	"embed"

	"github.com/ambientkit/ambient"
)

//go:embed template/*.tmpl
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

	p.Log.Info("plugin: enabled called")

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

// Routes sets routes for the plugin.
func (p *Plugin) Routes() {
	p.Log.Error("plugin: routes called")
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
	p.Mux.Get("/loadAllPluginPages", p.loadAllPluginPages)
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
}
