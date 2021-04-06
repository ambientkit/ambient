package core

// Grant is a type of permission.
type Grant string

const (
	// GrantAll allows all access.
	GrantAll Grant = "*"
	// GrantSiteTitleRead allows read access to the site title.
	GrantSiteTitleRead Grant = "site.title:read"
	// GrantSiteTitleWrite allows write access to the site title.
	GrantSiteTitleWrite Grant = "site.title:write"
	// GrantSiteContentRead allows read access to the site content.
	GrantSiteContentRead Grant = "site.content:read"
	// GrantSiteContentWrite allows write access to the site content.
	GrantSiteContentWrite Grant = "site.content:write"
	// GrantSiteSchemeRead allows read access to the site scheme.
	GrantSiteSchemeRead Grant = "site.scheme:read"
	// GrantSiteSchemeWrite allows write access to the site scheme.
	GrantSiteSchemeWrite Grant = "site.scheme:write"
	// GrantSiteURLRead allows read access to the site URL.
	GrantSiteURLRead Grant = "site.url:read"
	// GrantSiteURLWrite allows write access to the site URL.
	GrantSiteURLWrite Grant = "site.url:write"
	// GrantSiteUpdatedRead allows read access to the site updated time.
	// TODO: This doesn't have a write associated with it.
	GrantSiteUpdatedRead Grant = "site.updated:read"

	// GrantSiteLoadTrigger allows trigger access to the site load from data storage.
	GrantSiteLoadTrigger Grant = "site.load:trigger"

	// GrantSitePostRead allows read access to the site posts.
	// Allows access to calls like: postsandpages, publishedpages, postbyslug, tags.
	GrantSitePostRead Grant = "site.post:read"
	// GrantSitePostWrite allows write access to the site posts.
	GrantSitePostWrite Grant = "site.post:write"
	// GrantSitePostDelete allows delete access to the site posts.
	GrantSitePostDelete Grant = "site.post:delete"

	// GrantRouterRouteClear allows clear access to a route.
	GrantRouterRouteClear Grant = "router.route:clear"
	// GrantRouterNeighborRouteClear allows clear access to a route in another plugin.
	GrantRouterNeighborRouteClear Grant = "router.neighborroute:clear"

	// GrantSitePluginRead allows read access to the site plugins.
	GrantSitePluginRead Grant = "site.plugin:read"
	// GrantSitePluginEnable allows enable access to the site plugins.
	GrantSitePluginEnable Grant = "site.plugin:enable"
	// GrantSitePluginDisable allows disable access to the site plugins.
	GrantSitePluginDisable Grant = "site.plugin:disable"
	// GrantSitePluginDelete allows delete access to the site plugins.
	GrantSitePluginDelete Grant = "site.plugin:delete"

	// GrantPluginFieldRead allows read access to the plugin field.
	GrantPluginFieldRead Grant = "plugin.field:read"
	// GrantPluginFieldWrite allows write access to the plugin field.
	GrantPluginFieldWrite Grant = "plugin.field:write"
	// GrantPluginNeighborfieldRead allows read access to a field in another plugin.
	GrantPluginNeighborfieldRead Grant = "plugin.neighborfield:read"
	// GrantPluginNeighborfieldWrite allows write access to a field in another plugin.
	GrantPluginNeighborfieldWrite Grant = "plugin.neighborfield:write"

	// GrantUserAuthenticatedRead allows read access whether the current user is logged in or not.
	GrantUserAuthenticatedRead Grant = "user.authenticated:read"
)
