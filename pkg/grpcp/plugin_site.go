package grpcp

import (
	"context"
	"net/http"
	"time"

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
		return false
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

// PluginNeighborRoutesList handler.
func (c *GRPCSitePlugin) PluginNeighborRoutesList(pluginName string) ([]ambient.Route, error) {
	resp, err := c.client.PluginNeighborRoutesList(context.Background(), &protodef.SitePluginNeighborRoutesListRequest{
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
	_, err := c.client.UserPersist(context.Background(), &protodef.SiteUserPersistRequest{
		Requestid: requestID(r),
		Persist:   persist,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

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

	return out.Username, nil
}

// UserLogout handler.
func (c *GRPCSitePlugin) UserLogout(r *http.Request) error {
	_, err := c.client.UserLogout(context.Background(), &protodef.SiteUserLogoutRequest{
		Requestid: requestID(r),
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// LogoutAllUsers handler.
func (c *GRPCSitePlugin) LogoutAllUsers(r *http.Request) error {
	_, err := c.client.LogoutAllUsers(context.Background(), &protodef.SiteLogoutAllUsersRequest{
		Requestid: requestID(r),
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// SetCSRF handler.
func (c *GRPCSitePlugin) SetCSRF(r *http.Request) string {
	resp, err := c.client.SetCSRF(context.Background(), &protodef.SiteSetCSRFRequest{
		Requestid: requestID(r),
	})
	if err != nil {
		c.Log.Error("grpc-plugin: site.SetCSRF error: %v", err.Error())
		return ""
	}

	return resp.Token
}

// CSRF handler.
func (c *GRPCSitePlugin) CSRF(r *http.Request, token string) bool {
	resp, err := c.client.CSRF(context.Background(), &protodef.SiteCSRFRequest{
		Requestid: requestID(r),
		Token:     token,
	})
	if err != nil {
		c.Log.Error("grpc-plugin: site.CSRF error: %v", err.Error())
		return false
	}

	return resp.Valid
}

// SessionValue handler.
func (c *GRPCSitePlugin) SessionValue(r *http.Request, name string) string {
	resp, err := c.client.SessionValue(context.Background(), &protodef.SiteSessionValueRequest{
		Requestid: requestID(r),
		Name:      name,
	})
	if err != nil {
		c.Log.Error("grpc-plugin: site.SessionValue error: %v", err.Error())
		return ""
	}

	return resp.Value
}

// SetSessionValue handler.
func (c *GRPCSitePlugin) SetSessionValue(r *http.Request, name string, value string) error {
	_, err := c.client.SetSessionValue(context.Background(), &protodef.SiteSetSessionValueRequest{
		Requestid: requestID(r),
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
	_, err := c.client.DeleteSessionValue(context.Background(), &protodef.SiteDeleteSessionValueRequest{
		Requestid: requestID(r),
		Name:      name,
	})
	if err != nil {
		ErrorHandler(err)
	}
}

// PluginNeighborSettingsList handler.
func (c *GRPCSitePlugin) PluginNeighborSettingsList(pluginName string) ([]ambient.Setting, error) {
	settings := make([]ambient.Setting, 0)

	resp, err := c.client.PluginNeighborSettingsList(context.Background(), &protodef.SitePluginNeighborSettingsListRequest{
		Pluginname: pluginName,
	})
	if err != nil {
		return settings, ErrorHandler(err)
	}

	err = ProtobufStructToArray(resp.Settings, &settings)
	return settings, err
}

// SetPluginSetting handler.
func (c *GRPCSitePlugin) SetPluginSetting(settingName string, value string) error {
	_, err := c.client.SetPluginSetting(context.Background(), &protodef.SiteSetPluginSettingRequest{
		Settingname: settingName,
		Value:       value,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// PluginSettingBool handler.
func (c *GRPCSitePlugin) PluginSettingBool(fieldName string) (bool, error) {
	resp, err := c.client.PluginSettingBool(context.Background(), &protodef.SitePluginSettingBoolRequest{
		Fieldname: fieldName,
	})
	if err != nil {
		return false, ErrorHandler(err)
	}

	return resp.Value, nil
}

// PluginSettingString handler.
func (c *GRPCSitePlugin) PluginSettingString(fieldName string) (string, error) {
	resp, err := c.client.PluginSettingString(context.Background(), &protodef.SitePluginSettingStringRequest{
		Fieldname: fieldName,
	})
	if err != nil {
		return "", ErrorHandler(err)
	}

	return resp.Value, nil
}

// PluginSetting handler.
func (c *GRPCSitePlugin) PluginSetting(fieldName string) (interface{}, error) {
	resp, err := c.client.PluginSetting(context.Background(), &protodef.SitePluginSettingRequest{
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
func (c *GRPCSitePlugin) SetNeighborPluginSetting(pluginName string, settingName string, settingValue string) error {
	_, err := c.client.SetNeighborPluginSetting(context.Background(), &protodef.SiteSetNeighborPluginSettingRequest{
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
func (c *GRPCSitePlugin) NeighborPluginSettingString(pluginName string, fieldName string) (string, error) {
	resp, err := c.client.NeighborPluginSettingString(context.Background(), &protodef.SiteNeighborPluginSettingStringRequest{
		Pluginname: pluginName,
		Fieldname:  fieldName,
	})
	if err != nil {
		return "", ErrorHandler(err)
	}

	return resp.Value, nil
}

// NeighborPluginSetting handler.
func (c *GRPCSitePlugin) NeighborPluginSetting(pluginName string, fieldName string) (interface{}, error) {
	resp, err := c.client.NeighborPluginSetting(context.Background(), &protodef.SiteNeighborPluginSettingRequest{
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
func (c *GRPCSitePlugin) PluginTrusted(pluginName string) (bool, error) {
	resp, err := c.client.PluginTrusted(context.Background(), &protodef.SitePluginTrustedRequest{
		Pluginname: pluginName,
	})
	if err != nil {
		return false, ErrorHandler(err)
	}

	return resp.Trusted, nil
}

// SetTitle handler.
func (c *GRPCSitePlugin) SetTitle(title string) error {
	_, err := c.client.SetTitle(context.Background(), &protodef.SiteSetTitleRequest{
		Title: title,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// Title handler.
func (c *GRPCSitePlugin) Title() (string, error) {
	resp, err := c.client.Title(context.Background(), &protodef.Empty{})
	if err != nil {
		return "", ErrorHandler(err)
	}

	return resp.Title, nil
}

// SetScheme handler.
func (c *GRPCSitePlugin) SetScheme(scheme string) error {
	_, err := c.client.SetScheme(context.Background(), &protodef.SiteSetSchemeRequest{
		Scheme: scheme,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// Scheme handler.
func (c *GRPCSitePlugin) Scheme() (string, error) {
	resp, err := c.client.Scheme(context.Background(), &protodef.Empty{})
	if err != nil {
		return "", ErrorHandler(err)
	}

	return resp.Scheme, nil
}

// SetURL handler.
func (c *GRPCSitePlugin) SetURL(URL string) error {
	_, err := c.client.SetURL(context.Background(), &protodef.SiteSetURLRequest{
		Url: URL,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// URL handler.
func (c *GRPCSitePlugin) URL() (string, error) {
	resp, err := c.client.URL(context.Background(), &protodef.Empty{})
	if err != nil {
		return "", ErrorHandler(err)
	}

	return resp.Url, nil
}

// FullURL handler.
func (c *GRPCSitePlugin) FullURL() (string, error) {
	resp, err := c.client.FullURL(context.Background(), &protodef.Empty{})
	if err != nil {
		return "", ErrorHandler(err)
	}

	return resp.Fullurl, nil
}

// Updated handler.
func (c *GRPCSitePlugin) Updated() (time.Time, error) {
	resp, err := c.client.Updated(context.Background(), &protodef.Empty{})
	if err != nil {
		return time.Time{}, ErrorHandler(err)
	}

	return resp.Timestamp.AsTime(), nil
}
