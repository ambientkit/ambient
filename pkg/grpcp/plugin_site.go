package grpcp

import (
	"context"
	"net/http"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/internal/secureconfig"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
)

// GRPCSitePlugin is the plugin side implementation of secure site.
type GRPCSitePlugin struct {
	client protodef.SiteClient
	Log    ambient.Logger
}

// Error handler.
func (c *GRPCSitePlugin) Error(siteError error) error {
	return secureconfig.Error(siteError)
}

// Load handler.
func (c *GRPCSitePlugin) Load() error {
	_, err := c.client.Load(context.Background(), &protodef.Empty{})
	return ErrorHandler(err)
}

// Authorized handler.
func (c *GRPCSitePlugin) Authorized(grant ambient.Grant) bool {
	resp, err := c.client.Authorized(context.Background(), &protodef.SiteAuthorizedRequest{
		Grant: string(grant),
	})
	if err != nil {
		c.Log.Error("grpc-plugin: site.Authorized error: %v", err.Error())
	}
	return resp.Authorized
}

// NeighborPluginGrantList handler.
func (c *GRPCSitePlugin) NeighborPluginGrantList(pluginName string) ([]ambient.GrantRequest, error) {
	resp, err := c.client.NeighborPluginGrantList(context.Background(), &protodef.SiteNeighborPluginGrantListRequest{
		Pluginname: pluginName,
	})
	if err != nil {
		return []ambient.GrantRequest{}, ErrorHandler(err)
	}

	arr := make([]ambient.GrantRequest, 0)
	for _, v := range resp.Grants {
		arr = append(arr, ambient.GrantRequest{
			Grant:       ambient.Grant(v.Grant),
			Description: v.Description,
		})
	}

	return arr, nil
}

// NeighborPluginGrants handler.
func (c *GRPCSitePlugin) NeighborPluginGrants(pluginName string) (map[ambient.Grant]bool, error) {
	resp, err := c.client.NeighborPluginGrants(context.Background(), &protodef.SiteNeighborPluginGrantsRequest{
		Pluginname: pluginName,
	})
	if err != nil {
		return make(map[ambient.Grant]bool), ErrorHandler(err)
	}

	sm := make(map[ambient.Grant]bool)
	err = ProtobufStructToObject(resp.Grants, &sm)
	if err != nil {
		return make(map[ambient.Grant]bool), ErrorHandler(err)
	}

	return sm, nil
}

// NeighborPluginGranted handler.
func (c *GRPCSitePlugin) NeighborPluginGranted(pluginName string, grantName ambient.Grant) (bool, error) {
	resp, err := c.client.NeighborPluginGranted(context.Background(), &protodef.SiteNeighborPluginGrantedRequest{
		Pluginname: pluginName,
		Grant:      string(grantName),
	})
	if err != nil {
		return false, ErrorHandler(err)
	}

	return resp.Granted, nil
}

// NeighborPluginRequestedGrant handler.
func (c *GRPCSitePlugin) NeighborPluginRequestedGrant(pluginName string, grantName ambient.Grant) (bool, error) {
	resp, err := c.client.NeighborPluginRequestedGrant(context.Background(), &protodef.SiteNeighborPluginRequestedGrantRequest{
		Pluginname: pluginName,
		Grant:      string(grantName),
	})
	if err != nil {
		return false, ErrorHandler(err)
	}

	return resp.Granted, nil
}

// SetNeighborPluginGrant handler.
func (c *GRPCSitePlugin) SetNeighborPluginGrant(pluginName string, grantName ambient.Grant, granted bool) error {
	_, err := c.client.SetNeighborPluginGrant(context.Background(), &protodef.SiteSetNeighborPluginGrantRequest{
		Pluginname: pluginName,
		Grant:      string(grantName),
		Granted:    granted,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// Plugins handler.
func (c *GRPCSitePlugin) Plugins() (map[string]ambient.PluginData, error) {
	resp, err := c.client.Plugins(context.Background(), &protodef.Empty{})
	if err != nil {
		return make(map[string]ambient.PluginData), ErrorHandler(err)
	}

	sm := make(map[string]ambient.PluginData)
	err = ProtobufStructToObject(resp.Plugindata, &sm)
	if err != nil {
		return make(map[string]ambient.PluginData), ErrorHandler(err)
	}

	return sm, nil
}

// PluginNames handler.
func (c *GRPCSitePlugin) PluginNames() ([]string, error) {
	resp, err := c.client.PluginNames(context.Background(), &protodef.Empty{})
	if err != nil {
		return make([]string, 0), ErrorHandler(err)
	}

	return resp.Names, nil
}

// DeletePlugin handler.
func (c *GRPCSitePlugin) DeletePlugin(pluginName string) error {
	_, err := c.client.DeletePlugin(context.Background(), &protodef.SiteDeletePluginRequest{
		Name: pluginName,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// EnablePlugin handler.
func (c *GRPCSitePlugin) EnablePlugin(pluginName string, loadPlugin bool) error {
	_, err := c.client.EnablePlugin(context.Background(), &protodef.SiteEnablePluginRequest{
		Name: pluginName,
		Load: loadPlugin,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// LoadAllPluginPages handler.
func (c *GRPCSitePlugin) LoadAllPluginPages() error {
	_, err := c.client.LoadAllPluginPages(context.Background(), &protodef.Empty{})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// DisablePlugin handler.
func (c *GRPCSitePlugin) DisablePlugin(pluginName string, unloadPlugin bool) error {
	_, err := c.client.DisablePlugin(context.Background(), &protodef.SiteDisablePluginRequest{
		Name:   pluginName,
		Unload: unloadPlugin,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// SavePost handler.
func (c *GRPCSitePlugin) SavePost(ID string, post ambient.Post) error {
	ps, err := ObjectToProtobufStruct(post)
	if err != nil {
		return ErrorHandler(err)
	}

	_, err = c.client.SavePost(context.Background(), &protodef.SiteSavePostRequest{
		Id:   ID,
		Post: ps,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// PostsAndPages handler.
func (c *GRPCSitePlugin) PostsAndPages(onlyPublished bool) (ambient.PostWithIDList, error) {
	resp, err := c.client.PostsAndPages(context.Background(), &protodef.SitePostsAndPagesRequest{
		Onlypublished: onlyPublished,
	})
	if err != nil {
		return ambient.PostWithIDList{}, ErrorHandler(err)
	}

	post := make(ambient.PostWithIDList, 0)
	err = ProtobufStructToArray(resp.Postwithidlist, &post)
	return post, err
}

// PublishedPosts handler.
func (c *GRPCSitePlugin) PublishedPosts() ([]ambient.Post, error) {
	resp, err := c.client.PublishedPosts(context.Background(), &protodef.Empty{})
	if err != nil {
		return make([]ambient.Post, 0), ErrorHandler(err)
	}

	post := make([]ambient.Post, 0)
	err = ProtobufStructToArray(resp.Posts, &post)
	return post, err
}

// PublishedPages handler.
func (c *GRPCSitePlugin) PublishedPages() ([]ambient.Post, error) {
	resp, err := c.client.PublishedPages(context.Background(), &protodef.Empty{})
	if err != nil {
		return make([]ambient.Post, 0), ErrorHandler(err)
	}

	post := make([]ambient.Post, 0)
	err = ProtobufStructToArray(resp.Posts, &post)
	return post, err
}

// PostBySlug handler.
func (c *GRPCSitePlugin) PostBySlug(slug string) (ambient.PostWithID, error) {
	resp, err := c.client.PostBySlug(context.Background(), &protodef.SitePostBySlugRequest{
		Slug: slug,
	})
	if err != nil {
		return ambient.PostWithID{}, ErrorHandler(err)
	}

	post := ambient.PostWithID{}
	err = ProtobufStructToObject(resp.Post, &post)
	return post, err
}

// PostByID handler.
func (c *GRPCSitePlugin) PostByID(ID string) (ambient.Post, error) {
	resp, err := c.client.PostByID(context.Background(), &protodef.SitePostByIDRequest{
		Id: ID,
	})
	if err != nil {
		return ambient.Post{}, ErrorHandler(err)
	}

	post := ambient.Post{}
	err = ProtobufStructToObject(resp.Post, &post)
	return post, err
}

// DeletePostByID handler.
func (c *GRPCSitePlugin) DeletePostByID(ID string) error {
	_, err := c.client.DeletePostByID(context.Background(), &protodef.SiteDeletePostByIDRequest{
		Id: ID,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

/////////////////////////////////////////////////////

// UserLogin handler.
func (c *GRPCSitePlugin) UserLogin(r *http.Request, username string) error {
	_, err := c.client.UserLogin(context.Background(), &protodef.SiteUserLoginRequest{
		Username:  username,
		Requestid: requestID(r),
	})
	return err
}

// AuthenticatedUser handler.
func (c *GRPCSitePlugin) AuthenticatedUser(r *http.Request) (string, error) {
	out, err := c.client.AuthenticatedUser(context.Background(), &protodef.SiteAuthenticatedUserRequest{
		Requestid: requestID(r),
	})
	if err != nil {
		return "", ErrorHandler(err)
	}

	return out.Username, err
}
