package ambient

// GrantRequest represents a plugin grant request.
type GrantRequest struct {
	Grant       Grant
	Description string
}

// PluginGrants represents an unordered map of grants.
type PluginGrants map[Grant]bool

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

	// GrantSitePluginRead allows read access to the site plugins.
	GrantSitePluginRead Grant = "site.plugin:read"
	// GrantSitePluginEnable allows enable access to the site plugins.
	GrantSitePluginEnable Grant = "site.plugin:enable"
	// GrantSitePluginDisable allows disable access to the site plugins.
	GrantSitePluginDisable Grant = "site.plugin:disable"
	// GrantSitePluginDelete allows delete access to the site plugins.
	GrantSitePluginDelete Grant = "site.plugin:delete"

	// GrantRouterRouteWrite allows write access to routes.
	GrantRouterRouteWrite Grant = "router.route:write"
	// GrantRouterRouteClear allows clear access to a route.
	GrantRouterRouteClear Grant = "router.route:clear"
	// GrantRouterNeighborRouteClear allows clear access to a route in another plugin.
	GrantRouterNeighborRouteClear Grant = "router.neighborroute:clear"

	// GrantPluginSettingRead allows read access to the plugin setting.
	GrantPluginSettingRead Grant = "plugin.setting:read"
	// GrantPluginSettingWrite allows write access to the plugin setting.
	GrantPluginSettingWrite Grant = "plugin.setting:write"
	// GrantPluginNeighborSettingRead allows read access to a setting in another plugin.
	GrantPluginNeighborSettingRead Grant = "plugin.neighborsetting:read"
	// GrantPluginNeighborSettingWrite allows write access to a setting in another plugin.
	GrantPluginNeighborSettingWrite Grant = "plugin.neighborsetting:write"
	// GrantPluginNeighborGrantRead allows read access to a grant in another plugin.
	GrantPluginNeighborGrantRead Grant = "plugin.neighborgrant:read"
	// GrantPluginNeighborGrantWrite allows write access to a grant in another plugin.
	GrantPluginNeighborGrantWrite Grant = "plugin.neighborgrant:write"

	// GrantUserAuthenticatedRead allows read access whether the current user is logged in or not.
	GrantUserAuthenticatedRead Grant = "user.authenticated:read"
	// GrantUserAuthenticatedWrite allows write access to login or logout a user.
	GrantUserAuthenticatedWrite Grant = "user.authenticated:write"
	// GrantUserPersistWrite allows write access to login or logout a user.
	GrantUserPersistWrite Grant = "user.persist:write"

	// GrantSiteAssetWrite allows write access to site assets.
	GrantSiteAssetWrite Grant = "site.asset:write"
	// GrantSiteFuncMapWrite allows write access to site FuncMap for templates.
	GrantSiteFuncMapWrite Grant = "site.funcmap:write"
)
