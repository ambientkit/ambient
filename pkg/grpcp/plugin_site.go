package grpcp

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/internal/secureconfig"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"github.com/ambientkit/ambient/pkg/requestuuid"
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
func (c *GRPCSitePlugin) Load(ctx context.Context) error {
	_, err := c.client.Load(ctx, &protodef.Empty{})
	return ErrorHandler(err)
}

// LoadSinglePluginPages handler.
func (c *GRPCSitePlugin) LoadSinglePluginPages(ctx context.Context, name string) {
	_, err := c.client.LoadSinglePluginPages(ctx, &protodef.SiteLoadSinglePluginPagesRequest{
		Pluginname: name,
	})
	if err != nil {
		c.Log.Error("site.LoadSinglePluginPages error: %v", err.Error())
		return
	}
}

// Authorized handler.
func (c *GRPCSitePlugin) Authorized(ctx context.Context, grant ambient.Grant) bool {
	resp, err := c.client.Authorized(ctx, &protodef.SiteAuthorizedRequest{
		Grant: string(grant),
	})
	if err != nil {
		c.Log.Error("site.Authorized error: %v", err.Error())
		return false
	}

	return resp.Authorized
}

// NeighborPluginGrantList handler.
func (c *GRPCSitePlugin) NeighborPluginGrantList(ctx context.Context, pluginName string) ([]ambient.GrantRequest, error) {
	resp, err := c.client.NeighborPluginGrantList(ctx, &protodef.SiteNeighborPluginGrantListRequest{
		Pluginname: pluginName,
	})
	if err != nil {
		// FIXME: Need to determine if it should return nil and error? Or what kind of error?
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
func (c *GRPCSitePlugin) NeighborPluginGrants(ctx context.Context, pluginName string) (map[ambient.Grant]bool, error) {
	resp, err := c.client.NeighborPluginGrants(ctx, &protodef.SiteNeighborPluginGrantsRequest{
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
func (c *GRPCSitePlugin) NeighborPluginGranted(ctx context.Context, pluginName string, grantName ambient.Grant) (bool, error) {
	resp, err := c.client.NeighborPluginGranted(ctx, &protodef.SiteNeighborPluginGrantedRequest{
		Pluginname: pluginName,
		Grant:      string(grantName),
	})
	if err != nil {
		return false, ErrorHandler(err)
	}

	return resp.Granted, nil
}

// NeighborPluginRequestedGrant handler.
func (c *GRPCSitePlugin) NeighborPluginRequestedGrant(ctx context.Context, pluginName string, grantName ambient.Grant) (bool, error) {
	resp, err := c.client.NeighborPluginRequestedGrant(ctx, &protodef.SiteNeighborPluginRequestedGrantRequest{
		Pluginname: pluginName,
		Grant:      string(grantName),
	})
	if err != nil {
		return false, ErrorHandler(err)
	}

	return resp.Granted, nil
}

// SetNeighborPluginGrant handler.
func (c *GRPCSitePlugin) SetNeighborPluginGrant(ctx context.Context, pluginName string, grantName ambient.Grant, granted bool) error {
	_, err := c.client.SetNeighborPluginGrant(ctx, &protodef.SiteSetNeighborPluginGrantRequest{
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
func (c *GRPCSitePlugin) Plugins(ctx context.Context) (map[string]ambient.PluginData, error) {
	resp, err := c.client.Plugins(ctx, &protodef.Empty{})
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
func (c *GRPCSitePlugin) PluginNames(ctx context.Context) ([]string, error) {
	resp, err := c.client.PluginNames(ctx, &protodef.Empty{})
	if err != nil {
		return make([]string, 0), ErrorHandler(err)
	}

	return resp.Names, nil
}

// DeletePlugin handler.
func (c *GRPCSitePlugin) DeletePlugin(ctx context.Context, pluginName string) error {
	_, err := c.client.DeletePlugin(ctx, &protodef.SiteDeletePluginRequest{
		Name: pluginName,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// EnablePlugin handler.
func (c *GRPCSitePlugin) EnablePlugin(ctx context.Context, pluginName string, loadPlugin bool) error {
	_, err := c.client.EnablePlugin(ctx, &protodef.SiteEnablePluginRequest{
		Name: pluginName,
		Load: loadPlugin,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// DisablePlugin handler.
func (c *GRPCSitePlugin) DisablePlugin(ctx context.Context, pluginName string, unloadPlugin bool) error {
	_, err := c.client.DisablePlugin(ctx, &protodef.SiteDisablePluginRequest{
		Name:   pluginName,
		Unload: unloadPlugin,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// SavePost handler.
func (c *GRPCSitePlugin) SavePost(ctx context.Context, ID string, post ambient.Post) error {
	ps, err := ObjectToProtobufStruct(post)
	if err != nil {
		return ErrorHandler(err)
	}

	_, err = c.client.SavePost(ctx, &protodef.SiteSavePostRequest{
		Id:   ID,
		Post: ps,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// PostsAndPages handler.
func (c *GRPCSitePlugin) PostsAndPages(ctx context.Context, onlyPublished bool) (ambient.PostWithIDList, error) {
	resp, err := c.client.PostsAndPages(ctx, &protodef.SitePostsAndPagesRequest{
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
func (c *GRPCSitePlugin) PublishedPosts(ctx context.Context) ([]ambient.Post, error) {
	resp, err := c.client.PublishedPosts(ctx, &protodef.Empty{})
	if err != nil {
		return make([]ambient.Post, 0), ErrorHandler(err)
	}

	post := make([]ambient.Post, 0)
	err = ProtobufStructToArray(resp.Posts, &post)
	return post, err
}

// PublishedPages handler.
func (c *GRPCSitePlugin) PublishedPages(ctx context.Context) ([]ambient.Post, error) {
	resp, err := c.client.PublishedPages(ctx, &protodef.Empty{})
	if err != nil {
		return make([]ambient.Post, 0), ErrorHandler(err)
	}

	post := make([]ambient.Post, 0)
	err = ProtobufStructToArray(resp.Posts, &post)
	return post, err
}

// PostBySlug handler.
func (c *GRPCSitePlugin) PostBySlug(ctx context.Context, slug string) (ambient.PostWithID, error) {
	resp, err := c.client.PostBySlug(ctx, &protodef.SitePostBySlugRequest{
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
func (c *GRPCSitePlugin) PostByID(ctx context.Context, ID string) (ambient.Post, error) {
	resp, err := c.client.PostByID(ctx, &protodef.SitePostByIDRequest{
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
func (c *GRPCSitePlugin) DeletePostByID(ctx context.Context, ID string) error {
	_, err := c.client.DeletePostByID(ctx, &protodef.SiteDeletePostByIDRequest{
		Id: ID,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// PluginNeighborRoutesList handler.
func (c *GRPCSitePlugin) PluginNeighborRoutesList(ctx context.Context, pluginName string) ([]ambient.Route, error) {
	resp, err := c.client.PluginNeighborRoutesList(ctx, &protodef.SitePluginNeighborRoutesListRequest{
		Pluginname: pluginName,
	})
	if err != nil {
		return make([]ambient.Route, 0), ErrorHandler(err)
	}

	post := make([]ambient.Route, 0)
	err = ProtobufStructToArray(resp.Routes, &post)
	return post, err
}

// UserPersist handler.
func (c *GRPCSitePlugin) UserPersist(r *http.Request, persist bool) error {
	_, err := c.client.UserPersist(r.Context(), &protodef.SiteUserPersistRequest{
		Requestid: requestuuid.Get(r),
		Persist:   persist,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// UserLogin handler.
func (c *GRPCSitePlugin) UserLogin(r *http.Request, username string) error {
	_, err := c.client.UserLogin(r.Context(), &protodef.SiteUserLoginRequest{
		Username:  username,
		Requestid: requestuuid.Get(r),
	})
	return err
}

// AuthenticatedUser handler.
func (c *GRPCSitePlugin) AuthenticatedUser(r *http.Request) (string, error) {
	out, err := c.client.AuthenticatedUser(r.Context(), &protodef.SiteAuthenticatedUserRequest{
		Requestid: requestuuid.Get(r),
	})
	if err != nil {
		return "", ErrorHandler(err)
	}

	if len(out.Username) == 0 {
		return "", errors.New("user not found")
	}

	return out.Username, nil
}

// UserLogout handler.
func (c *GRPCSitePlugin) UserLogout(r *http.Request) error {
	_, err := c.client.UserLogout(r.Context(), &protodef.SiteUserLogoutRequest{
		Requestid: requestuuid.Get(r),
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// LogoutAllUsers handler.
func (c *GRPCSitePlugin) LogoutAllUsers(r *http.Request) error {
	_, err := c.client.LogoutAllUsers(r.Context(), &protodef.SiteLogoutAllUsersRequest{
		Requestid: requestuuid.Get(r),
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// SetCSRF handler.
func (c *GRPCSitePlugin) SetCSRF(r *http.Request) string {
	resp, err := c.client.SetCSRF(r.Context(), &protodef.SiteSetCSRFRequest{
		Requestid: requestuuid.Get(r),
	})
	if err != nil {
		c.Log.Error("site.SetCSRF error: %v", err.Error())
		return ""
	}

	return resp.Token
}

// CSRF handler.
func (c *GRPCSitePlugin) CSRF(r *http.Request, token string) bool {
	resp, err := c.client.CSRF(r.Context(), &protodef.SiteCSRFRequest{
		Requestid: requestuuid.Get(r),
		Token:     token,
	})
	if err != nil {
		c.Log.Error("site.CSRF error: %v", err.Error())
		return false
	}

	return resp.Valid
}

// SessionValue handler.
func (c *GRPCSitePlugin) SessionValue(r *http.Request, name string) string {
	resp, err := c.client.SessionValue(r.Context(), &protodef.SiteSessionValueRequest{
		Requestid: requestuuid.Get(r),
		Name:      name,
	})
	if err != nil {
		c.Log.Error("site.SessionValue error: %v", err.Error())
		return ""
	}

	return resp.Value
}

// SetSessionValue handler.
func (c *GRPCSitePlugin) SetSessionValue(r *http.Request, name string, value string) error {
	_, err := c.client.SetSessionValue(r.Context(), &protodef.SiteSetSessionValueRequest{
		Requestid: requestuuid.Get(r),
		Name:      name,
		Value:     value,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// DeleteSessionValue handler.
func (c *GRPCSitePlugin) DeleteSessionValue(r *http.Request, name string) {
	_, err := c.client.DeleteSessionValue(r.Context(), &protodef.SiteDeleteSessionValueRequest{
		Requestid: requestuuid.Get(r),
		Name:      name,
	})
	if err != nil {
		ErrorHandler(err)
	}
}

// PluginNeighborSettingsList handler.
func (c *GRPCSitePlugin) PluginNeighborSettingsList(ctx context.Context, pluginName string) ([]ambient.Setting, error) {
	settings := make([]ambient.Setting, 0)

	resp, err := c.client.PluginNeighborSettingsList(ctx, &protodef.SitePluginNeighborSettingsListRequest{
		Pluginname: pluginName,
	})
	if err != nil {
		return settings, ErrorHandler(err)
	}

	err = ProtobufStructToArray(resp.Settings, &settings)
	return settings, err
}

// SetPluginSetting handler.
func (c *GRPCSitePlugin) SetPluginSetting(ctx context.Context, settingName string, value string) error {
	_, err := c.client.SetPluginSetting(ctx, &protodef.SiteSetPluginSettingRequest{
		Settingname: settingName,
		Value:       value,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// PluginSettingBool handler.
func (c *GRPCSitePlugin) PluginSettingBool(ctx context.Context, fieldName string) (bool, error) {
	resp, err := c.client.PluginSettingBool(ctx, &protodef.SitePluginSettingBoolRequest{
		Fieldname: fieldName,
	})
	if err != nil {
		return false, ErrorHandler(err)
	}

	return resp.Value, nil
}

// PluginSettingString handler.
func (c *GRPCSitePlugin) PluginSettingString(ctx context.Context, fieldName string) (string, error) {
	resp, err := c.client.PluginSettingString(ctx, &protodef.SitePluginSettingStringRequest{
		Fieldname: fieldName,
	})
	if err != nil {
		return "", ErrorHandler(err)
	}

	return resp.Value, nil
}

// PluginSetting handler.
func (c *GRPCSitePlugin) PluginSetting(ctx context.Context, fieldName string) (interface{}, error) {
	resp, err := c.client.PluginSetting(ctx, &protodef.SitePluginSettingRequest{
		Fieldname: fieldName,
	})
	if err != nil {
		return "", ErrorHandler(err)
	}

	var i interface{}
	err = ProtobufAnyToInterface(resp.Value, &i)
	return i, err
}

// SetNeighborPluginSetting handler.
func (c *GRPCSitePlugin) SetNeighborPluginSetting(ctx context.Context, pluginName string, settingName string, settingValue string) error {
	_, err := c.client.SetNeighborPluginSetting(ctx, &protodef.SiteSetNeighborPluginSettingRequest{
		Pluginname:   pluginName,
		Settingname:  settingName,
		Settingvalue: settingValue,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// NeighborPluginSettingString handler.
func (c *GRPCSitePlugin) NeighborPluginSettingString(ctx context.Context, pluginName string, fieldName string) (string, error) {
	resp, err := c.client.NeighborPluginSettingString(ctx, &protodef.SiteNeighborPluginSettingStringRequest{
		Pluginname: pluginName,
		Fieldname:  fieldName,
	})
	if err != nil {
		return "", ErrorHandler(err)
	}

	return resp.Value, nil
}

// NeighborPluginSetting handler.
func (c *GRPCSitePlugin) NeighborPluginSetting(ctx context.Context, pluginName string, fieldName string) (interface{}, error) {
	resp, err := c.client.NeighborPluginSetting(ctx, &protodef.SiteNeighborPluginSettingRequest{
		Pluginname: pluginName,
		Fieldname:  fieldName,
	})
	if err != nil {
		return "", ErrorHandler(err)
	}

	var i interface{}
	err = ProtobufAnyToInterface(resp.Value, &i)
	return i, err
}

// PluginTrusted handler.
func (c *GRPCSitePlugin) PluginTrusted(ctx context.Context, pluginName string) (bool, error) {
	resp, err := c.client.PluginTrusted(ctx, &protodef.SitePluginTrustedRequest{
		Pluginname: pluginName,
	})
	if err != nil {
		return false, ErrorHandler(err)
	}

	return resp.Trusted, nil
}

// SetTitle handler.
func (c *GRPCSitePlugin) SetTitle(ctx context.Context, title string) error {
	_, err := c.client.SetTitle(ctx, &protodef.SiteSetTitleRequest{
		Title: title,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// Title handler.
func (c *GRPCSitePlugin) Title(ctx context.Context) (string, error) {
	resp, err := c.client.Title(ctx, &protodef.Empty{})
	if err != nil {
		return "", ErrorHandler(err)
	}

	return resp.Title, nil
}

// SetScheme handler.
func (c *GRPCSitePlugin) SetScheme(ctx context.Context, scheme string) error {
	_, err := c.client.SetScheme(ctx, &protodef.SiteSetSchemeRequest{
		Scheme: scheme,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// Scheme handler.
func (c *GRPCSitePlugin) Scheme(ctx context.Context) (string, error) {
	resp, err := c.client.Scheme(ctx, &protodef.Empty{})
	if err != nil {
		return "", ErrorHandler(err)
	}

	return resp.Scheme, nil
}

// SetURL handler.
func (c *GRPCSitePlugin) SetURL(ctx context.Context, URL string) error {
	_, err := c.client.SetURL(ctx, &protodef.SiteSetURLRequest{
		Url: URL,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// URL handler.
func (c *GRPCSitePlugin) URL(ctx context.Context) (string, error) {
	resp, err := c.client.URL(ctx, &protodef.Empty{})
	if err != nil {
		return "", ErrorHandler(err)
	}

	return resp.Url, nil
}

// FullURL handler.
func (c *GRPCSitePlugin) FullURL(ctx context.Context) (string, error) {
	resp, err := c.client.FullURL(ctx, &protodef.Empty{})
	if err != nil {
		return "", ErrorHandler(err)
	}

	return resp.Fullurl, nil
}

// Updated handler.
func (c *GRPCSitePlugin) Updated(ctx context.Context) (time.Time, error) {
	resp, err := c.client.Updated(ctx, &protodef.Empty{})
	if err != nil {
		return time.Time{}, ErrorHandler(err)
	}

	return resp.Timestamp.AsTime(), nil
}

// SetContent handler.
func (c *GRPCSitePlugin) SetContent(ctx context.Context, content string) error {
	_, err := c.client.SetContent(ctx, &protodef.SiteSetContentRequest{
		Content: content,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// Content handler.
func (c *GRPCSitePlugin) Content(ctx context.Context) (string, error) {
	resp, err := c.client.Content(ctx, &protodef.Empty{})
	if err != nil {
		return "", ErrorHandler(err)
	}

	return resp.Content, nil
}

// Tags handler.
func (c *GRPCSitePlugin) Tags(ctx context.Context, onlyPublished bool) (ambient.TagList, error) {
	tags := make(ambient.TagList, 0)
	resp, err := c.client.Tags(ctx, &protodef.SiteTagsRequest{
		Onlypublished: onlyPublished,
	})
	if err != nil {
		return tags, ErrorHandler(err)
	}

	for _, v := range resp.Tags {
		tags = append(tags, ambient.Tag{
			Name:      v.Name,
			Timestamp: v.Timestamp.AsTime(),
		})
	}

	return tags, nil
}
